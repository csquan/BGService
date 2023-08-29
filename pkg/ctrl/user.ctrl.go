package ctrl

import (
	"context"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"user/pkg/conf"
	"user/pkg/db"
	"user/pkg/log"
	"user/pkg/model"
	"user/pkg/redis"
	"user/pkg/util"
	"user/pkg/util/rsa0"
	"user/pkg/web"
)

type UserCtrl struct {
	*Ctrl
	SafeCli *web.MyInnerCli
}

func NewUserCtrl() UserCtrl {
	rdb := redis.NewStore("go:user:controller:")
	logger := log.C.Logger().With().Str("ctrl", "user").Logger()
	safeCli, err := web.NewCli("http://safe-service/", conf.Conf.KeySvrPubKey)
	if err != nil {
		logger.Error().Msg(err.LStr())
		panic(err)
	}
	return UserCtrl{
		Ctrl:    &Ctrl{DB: db.DB, RDB: rdb, Log: &logger},
		SafeCli: safeCli,
	}
}

// GenGa godoc
//
//	@Summary		在校验手机或者校验邮箱成功后，创建谷歌验证
//	@Description	生成绑定谷歌验证需要的二维码
//	@Tags			user
//	@Accept			json
//	@Product		json
//	@Param			Authorization	header	string	true	"Authentication header"
//	@Security		ApiKeyAuth
//	@Success		200	{object}	web.HttpData[model.GenGaResponse]
//	@Failure		502	{object}	web.HttpMsg
//	@Router			/user/genGa [post]
func (uc *UserCtrl) GenGa(c *gin.Context) {
	tu, ok := mustGetTu(uc.Ctrl, c)
	if !ok {
		return
	}

	if tu.BID == "" {
		web.BadRes(c, util.ErrUnAuthed)
		uc.Log.Debug().Interface("tu", tu).Send()
		return
	}

	if _, ok := dbUserByID(uc.Ctrl, c, tu.BID); !ok {
		return
	}
	if gr, ok := genGa(uc.Ctrl, c, tu); ok {
		web.GoodResp(c, gr)
	}
}

// BindGa godoc
//
//	@Summary		绑定谷歌验证
//	@Description	成功返回一个 token 保存当前用户的所有信息
//	@Tags			user
//	@Accept			json
//	@Product		json
//	@Param			code			body	model.BindGAInput	true	"ga code & (email-code | sms-code)"
//	@Param			Authorization	header	string				true	"Authentication header"
//	@Security		ApiKeyAuth
//	@Success		200	{object}	web.HttpData[string]
//	@Failure		502	{object}	web.HttpMsg
//	@Router			/user/bindGa [post]
func (uc *UserCtrl) BindGa(c *gin.Context) {
	tu, ok := mustGetTu(uc.Ctrl, c)
	if !ok {
		return
	}

	var bindGAInput model.BindGAInput
	if err := c.ShouldBindJSON(&bindGAInput); err != nil {
		web.BadRes(c, util.ErrInvalidArgument)
		uc.Log.Err(err).Send()
		return
	}

	var user model.User
	user, ok = dbUserByID(uc.Ctrl, c, tu.BID)
	if !ok {
		return
	}

	redisMKey := ""
	if user.LastMVTime > 0 {
		redisMKey, ok = checkUserMobileCode(uc.Ctrl, c, &user, bindGAInput.MCode, true, true)
		if !ok {
			return
		}
	}
	redisEKey := ""
	if user.LastEVTime > 0 {
		redisEKey, ok = checkUserEmailCode(uc.Ctrl, c, &user, bindGAInput.ECode, true, true)
		if !ok {
			return
		}
	}
	redisGKey, ok := bindGa(uc.Ctrl, c, &user, bindGAInput.GCode)
	if !ok {
		return
	}
	if err := uc.DB.Model(&user).Updates(map[string]interface{}{
		"ga":           user.Ga,
		"last_gv_time": user.LastGVTime,
		"last_ev_time": user.LastEVTime,
		"last_mv_time": user.LastMVTime,
	}).Error; err != nil {
		web.BadRes(c, util.ErrDB)
		uc.Log.Err(err).Send()
		return
	}

	c0 := context.Background()
	if redisEKey != "" {
		_, _ = uc.RDB.Delete(c0, redisEKey)
	}
	if redisMKey != "" {
		_, _ = uc.RDB.Delete(c0, redisMKey)
	}
	_, _ = uc.RDB.Delete(c0, redisGKey)
	tu2 := user.ToTu(tu)
	AToken(uc.Ctrl, c, tu2)
}

// ResEmail godoc
//
//	@Summary		已登录用户再次验证电子邮箱，发起认证是否本人操作
//	@Description	带token操作，无需主动输入参数。后台生成 6 位随机数字邮件发送并保存在 Redis，5分钟过期
//	@Tags			user
//	@Accept			json
//	@Product		json
//	@Param			Authorization	header	string	true	"Authentication header"
//	@Security		ApiKeyAuth
//	@Success		200	{object}	web.HttpData[string]
//	@Failure		502	{object}	web.HttpMsg
//	@Router			/user/resEmail [post]
func (uc *UserCtrl) ResEmail(c *gin.Context) {
	tu, ok := mustGetTu(uc.Ctrl, c)
	if !ok {
		return
	}

	user, ok := dbUserByID(uc.Ctrl, c, tu.BID)
	if !ok {
		return
	}

	if user.Email == "" {
		err := util.ErrEmailNo
		web.BadRes(c, err)
		uc.Log.Error().Msg(err.LStr())
		return
	}

	sendEmailNStore(uc.Ctrl, c, user.Email)
}

// ResMobile godoc
//
//	@Summary		提交验证手机号，验证是否本人操作
//	@Description	已经登录情况下操作，无需其他参数。后台生成 6 位随机数字短信发送给用户并保存在 Redis，5分钟过期
//	@Tags			user
//	@Accept			json
//	@Product		json
//	@Param			Authorization	header	string	true	"Authentication header"
//	@Security		ApiKeyAuth
//	@Success		200	{object}	web.HttpData[string]
//	@Failure		502	{object}	web.HttpMsg
//	@Router			/user/resMobile [post]
func (uc *UserCtrl) ResMobile(c *gin.Context) {
	tu, ok := mustGetTu(uc.Ctrl, c)
	if !ok {
		return
	}

	user, ok := dbUserByID(uc.Ctrl, c, tu.BID)
	if !ok {
		return
	}

	if user.Mobile == "" {
		err := util.ErrMobileNo
		web.BadRes(c, err)
		uc.Log.Error().Msg(err.LStr())
		return
	}

	sendMobileNStore(uc.Ctrl, c, user.Mobile)
}

// RevEmail godoc
//
//	@Summary		已登录状态下校验邮件验证码，证明是本人操作
//	@Description	成功返回一个新 token
//	@Tags			user
//	@Accept			json
//	@Product		json
//	@Param			Authorization	header	string			true	"Authentication header"
//	@Param			input			body	model.CodeInput	true	"输入你收到的验证码"
//	@Security		ApiKeyAuth
//	@Success		200	{object}	web.HttpData[string]
//	@Failure		502	{object}	web.HttpMsg
//	@Router			/user/revEmail [post]
func (uc *UserCtrl) RevEmail(c *gin.Context) {
	var reVerifyInput model.CodeInput
	if err := c.ShouldBindJSON(&reVerifyInput); err != nil {
		web.BadRes(c, util.ErrInvalidArgument)
		uc.Log.Err(err).Send()
		return
	}
	code := reVerifyInput.Code

	tu, ok := mustGetTu(uc.Ctrl, c)
	if !ok {
		return
	}

	user, ok := dbUserByID(uc.Ctrl, c, tu.BID)
	if !ok {
		return
	}

	if user.Email == "" {
		err := util.ErrEmailNo
		web.BadRes(c, err)
		uc.Log.Error().Msg(err.LStr())
		return
	}

	redisKey, ok := checkUserEmailCode(uc.Ctrl, c, &user, code, false, true)
	if !ok {
		return
	}

	if err := uc.DB.Model(&user).Update("last_ev_time", user.LastEVTime).Error; err != nil {
		web.BadRes(c, util.ErrDB)
		uc.Log.Err(err).Send()
		return
	}

	_, _ = uc.RDB.Delete(context.Background(), redisKey)

	tu2 := user.ToTu(tu)
	AToken(uc.Ctrl, c, tu2)
}

// RevMobile godoc
//
//	@Summary		校验手机验证码，成功证明是本人手机
//	@Description	成功返回一个 token 可以用于保存当前用户的所有信息
//	@Tags			user
//	@Accept			json
//	@Product		json
//	@Param			Authorization	header	string			true	"Authentication header"
//	@Param			input			body	model.CodeInput	true	"输入你收到的验证码"
//	@Security		ApiKeyAuth
//	@Success		200	{object}	web.HttpData[string]
//	@Failure		502	{object}	web.HttpMsg
//	@Router			/user/revMobile [post]
func (uc *UserCtrl) RevMobile(c *gin.Context) {
	var reVerifyInput model.CodeInput
	if err := c.ShouldBindJSON(&reVerifyInput); err != nil {
		web.BadRes(c, util.ErrInvalidArgument)
		uc.Log.Err(err).Send()
		return
	}
	code := reVerifyInput.Code

	tu, ok := mustGetTu(uc.Ctrl, c)
	if !ok {
		return
	}

	user, ok := dbUserByID(uc.Ctrl, c, tu.BID)
	if !ok {
		return
	}

	if user.Mobile == "" {
		err := util.ErrMobileNo
		web.BadRes(c, err)
		uc.Log.Error().Msg(err.LStr())
		return
	}

	redisKey, ok := checkUserMobileCode(uc.Ctrl, c, &user, code, false, true)
	if !ok {
		return
	}

	if err := uc.DB.Model(&user).Update("last_mv_time", user.LastMVTime).Error; err != nil {
		web.BadRes(c, util.ErrDB)
		uc.Log.Err(err).Send()
		return
	}
	_, _ = uc.RDB.Delete(context.Background(), redisKey)

	tu2 := user.ToTu(tu)

	AToken(uc.Ctrl, c, tu2)
}

// MyEmail godoc
//
//	@Summary		已生成用户并登录手机的情况下提交电子邮箱，要求发起认证真实性
//	@Description	后台生成 6 位随机数字邮件发送并保存在 Redis，5分钟过期
//	@Tags			user
//	@Accept			json
//	@Product		json
//	@Param			email			body	model.EmailInput	true	"在 email 字段填入你的邮箱"
//	@Param			Authorization	header	string				true	"Authentication header"
//	@Security		ApiKeyAuth
//	@Success		200	{object}	web.HttpData[string]
//	@Failure		502	{object}	web.HttpMsg
//	@Router			/user/myEmail [post]
func (uc *UserCtrl) MyEmail(c *gin.Context) {
	var emailInput model.EmailInput
	if err := c.ShouldBindJSON(&emailInput); err != nil {
		web.BadRes(c, util.ErrInvalidArgument)
		uc.Log.Err(err).Send()
		return
	}
	if err := emailInput.Validate(); err != nil {
		web.BadRes(c, err)
		uc.Log.Error().Msg(err.LStr())
		return
	}
	email := emailInput.Email

	if !ensureEmailNotReg(uc.Ctrl, c, email) {
		return
	}

	tu, ok := mustGetTu(uc.Ctrl, c)
	if !ok {
		return
	}
	user, ok := dbUserByID(uc.Ctrl, c, tu.BID)
	if !ok {
		return
	}

	if tu.Mobile == "" {
		web.BadRes(c, util.ErrUnAuthed)
		uc.Log.Debug().Interface("tu=", tu).Send()
		return
	}
	if user.LastEVTime > 0 && user.Email != "" {
		err := util.ErrEmailAlready(user.Email)
		web.BadRes(c, err)
		uc.Log.Error().Msg(err.LStr())
		return
	}

	user.Email = email
	if err := uc.DB.Model(&user).Update("email", user.Email).Error; err != nil {
		web.BadRes(c, util.ErrDB)
		uc.Log.Err(err).Send()
		return
	}
	sendEmailNStore(uc.Ctrl, c, email)
}

// MyMobile godoc
//
//	@Summary		已注册并登录邮箱的情况下提交手机号，要求发起认证真实性
//	@Description	后台生成 6 位随机数字短信发送给用户并保存在 Redis，5分钟过期
//	@Tags			user
//	@Accept			json
//	@Product		json
//	@Param			mobile			body	model.MobileInput	true	"在 mobile 字段填入你的手机号，带国家字冠"
//	@Param			Authorization	header	string				true	"Authentication header"
//	@Security		ApiKeyAuth
//	@Success		200	{object}	web.HttpData[string]
//	@Failure		502	{object}	web.HttpMsg
//	@Router			/user/myMobile [post]
func (uc *UserCtrl) MyMobile(c *gin.Context) {
	var mobileInput model.MobileInput
	if err := c.ShouldBindJSON(&mobileInput); err != nil {
		web.BadRes(c, util.ErrInvalidArgument)
		uc.Log.Err(err).Send()
		return
	}
	if err := mobileInput.Validate(); err != nil {
		web.BadRes(c, err)
		uc.Log.Error().Msg(err.LStr())
		return
	}
	mobile := mobileInput.Mobile

	if !ensureMobileNotReg(uc.Ctrl, c, mobile) {
		return
	}

	tu, ok := mustGetTu(uc.Ctrl, c)
	if !ok {
		return
	}
	user, ok := dbUserByID(uc.Ctrl, c, tu.BID)
	if !ok {
		return
	}
	if user.LastMVTime > 0 && user.Mobile != "" {
		err := util.ErrMobileAlready(user.Mobile)
		web.BadRes(c, err)
		uc.Log.Error().Msg(err.LStr())
		return
	}

	user.Mobile = mobile
	if err := uc.DB.Model(&user).Update("mobile", user.Mobile).Error; err != nil {
		web.BadRes(c, util.ErrDB)
		uc.Log.Err(err).Send()
		return
	}

	sendMobileNStore(uc.Ctrl, c, mobile)
}

// VerifyGa godoc
//
//	@Summary		单独验证谷歌验证
//	@Description	成功返回一个 token 保存当前用户的所有信息
//	@Tags			user
//	@Accept			json
//	@Product		json
//	@Param			code			body	model.CodeInput	true	"ga code"
//	@Param			Authorization	header	string			true	"Authentication header"
//	@Security		ApiKeyAuth
//	@Success		200	{object}	web.HttpData[string]
//	@Failure		502	{object}	web.HttpMsg
//	@Router			/user/verifyGa [post]
func (uc *UserCtrl) VerifyGa(c *gin.Context) {
	var gaInput model.CodeInput
	if err := c.ShouldBindJSON(&gaInput); err != nil {
		web.BadRes(c, util.ErrInvalidArgument)
		uc.Log.Err(err).Send()
		return
	}

	tu, ok := mustGetTu(uc.Ctrl, c)
	if !ok {
		return
	}
	user, ok := dbUserByID(uc.Ctrl, c, tu.BID)
	if !ok {
		return
	}

	v, er := util.ValidateTOTP(gaInput.Code, user.Ga)
	if er != nil || !v {
		web.BadRes(c, util.ErrGaInvalid)
		uc.Log.Error().Msg(er.LStr())
		return
	}
	user.LastGVTime = time.Now().Unix()
	if err := uc.DB.Model(&user).Update("last_gv_time", user.LastGVTime).Error; err != nil {
		web.BadRes(c, util.ErrDB)
		uc.Log.Err(err).Send()
		return
	}
	tu2 := user.ToTu(tu)

	AToken(uc.Ctrl, c, tu2)
}

// MyNewEmail godoc
//
//	@Summary		已经认证过电子邮箱的情况下，要求换电子邮箱。如果从未验证过邮箱，需要调用 myEmail 接口
//	@Description	提交新电子邮箱，要求发起认证真实性。后台生成 6 位随机数字邮件发送并保存在 Redis，5分钟过期
//	@Tags			user
//	@Accept			json
//	@Product		json
//	@Param			email			body	model.EmailInput	true	"在 email 字段填入你的邮箱"
//	@Param			Authorization	header	string				true	"Authentication header"
//	@Security		ApiKeyAuth
//	@Success		200	{object}	web.HttpData[string]
//	@Failure		502	{object}	web.HttpMsg
//	@Router			/user/myNewEmail [post]
func (uc *UserCtrl) MyNewEmail(c *gin.Context) {
	var emailInput model.EmailInput
	if err := c.ShouldBindJSON(&emailInput); err != nil {
		web.BadRes(c, util.ErrInvalidArgument)
		uc.Log.Err(err).Send()
		return
	}
	if err := emailInput.Validate(); err != nil {
		web.BadRes(c, err)
		uc.Log.Error().Msg(err.LStr())
		return
	}
	email := emailInput.Email

	if !ensureEmailNotReg(uc.Ctrl, c, email) {
		return
	}

	tu, ok := mustGetTu(uc.Ctrl, c)
	if !ok {
		return
	}
	user, ok := dbUserByID(uc.Ctrl, c, tu.BID)
	if !ok {
		return
	}
	if user.LastEVTime == 0 {
		err := util.ErrEmailNo
		web.BadRes(c, err)
		uc.Log.Error().Msg(err.LStr())
		return
	}

	sendEmailNStore(uc.Ctrl, c, email)
}

// MyNewMobile godoc
//
//	@Summary		已经认证过手机的情况下要求更换手机号，要求对新手机号发起真实性认证。
//	@Description	后台生成 6 位随机数字短信发送给用户并保存在 Redis，5分钟过期
//	@Tags			user
//	@Accept			json
//	@Product		json
//	@Param			mobile			body	model.MobileInput	true	"在 mobile 字段填入你的手机号，带国家字冠"
//	@Param			Authorization	header	string				true	"Authentication header"
//	@Security		ApiKeyAuth
//	@Success		200	{object}	web.HttpData[string]
//	@Failure		502	{object}	web.HttpMsg
//	@Router			/user/myNewMobile [post]
func (uc *UserCtrl) MyNewMobile(c *gin.Context) {
	var mobileInput model.MobileInput
	if err := c.ShouldBindJSON(&mobileInput); err != nil {
		web.BadRes(c, util.ErrInvalidArgument)
		uc.Log.Err(err).Send()
		return
	}
	if err := mobileInput.Validate(); err != nil {
		web.BadRes(c, err)
		uc.Log.Error().Msg(err.LStr())
		return
	}
	mobile := mobileInput.Number()

	if !ensureMobileNotReg(uc.Ctrl, c, mobile) {
		return
	}

	tu, ok := mustGetTu(uc.Ctrl, c)
	if !ok {
		return
	}
	user, ok := dbUserByID(uc.Ctrl, c, tu.BID)
	if !ok {
		return
	}

	if user.LastMVTime == 0 {
		err := util.ErrMobileNo
		web.BadRes(c, err)
		uc.Log.Error().Msg(err.LStr())
		return
	}

	sendMobileNStore(uc.Ctrl, c, mobile)
}

// ChangeEmail godoc
//
//	@Summary		只认证了电子邮箱的情况下，要求换电子邮箱。提交新邮箱和验证码。
//	@Description	认证成功更换邮箱。要求刚验证过旧邮箱
//	@Tags			user
//	@Accept			json
//	@Product		json
//	@Param			email			body	model.VerifyEmailInput	true	"新邮箱和验证码"
//	@Param			Authorization	header	string					true	"Authentication header"
//	@Security		ApiKeyAuth
//	@Success		200	{object}	web.HttpData[string]
//	@Failure		502	{object}	web.HttpMsg
//	@Router			/user/changeEmail [post]
func (uc *UserCtrl) ChangeEmail(c *gin.Context) {
	var ceInput model.VerifyEmailInput
	if err := c.ShouldBindJSON(&ceInput); err != nil {
		web.BadRes(c, util.ErrInvalidArgument)
		uc.Log.Err(err).Send()
		return
	}
	if err := ceInput.Validate(); err != nil {
		web.BadRes(c, err)
		uc.Log.Error().Msg(err.LStr())
		return
	}
	email := ceInput.Email

	if !ensureEmailNotReg(uc.Ctrl, c, email) {
		return
	}

	tu, ok := mustGetTu(uc.Ctrl, c)
	if !ok {
		return
	}
	if !recentToken(uc.Ctrl, c, tu) {
		return
	}
	user, ok := dbUserByID(uc.Ctrl, c, tu.BID)
	if !ok {
		return
	}

	if user.LastEVTime == 0 {
		err := util.ErrEmailNo
		web.BadRes(c, err)
		uc.Log.Error().Msg(err.LStr())
		return
	}
	if user.LastMVTime > 0 {
		err := util.ErrEmailByMobile(user.Mobile)
		web.BadRes(c, err)
		uc.Log.Error().Msg(err.LStr())
		return
	}
	if user.LastGVTime > 0 {
		err := util.ErrEmailByGa
		web.BadRes(c, err)
		uc.Log.Error().Msg(err.LStr())
		return
	}

	if !emailCodeExists(uc.Ctrl, c, email, ceInput.Code) {
		return
	}

	user.Email = email
	user.LastEVTime = time.Now().Unix()

	if err := uc.DB.Model(&user).Updates(map[string]interface{}{
		"email":        user.Email,
		"last_ev_time": user.LastEVTime,
	}).Error; err != nil {
		web.BadRes(c, util.ErrDB)
		uc.Log.Err(err).Send()
		return
	}

	tu2 := user.ToTu(tu)
	AToken(uc.Ctrl, c, tu2)
}

// ChangeMobile godoc
//
//	@Summary		只认证了手机的情况下，要求换手机号码。提交新手机和验证码。
//	@Description	认证成功更换手机号。要求刚认证过旧手机
//	@Tags			user
//	@Accept			json
//	@Product		json
//	@Param			mobile			body	model.VerifyMobileInput	true	"新手机和验证码"
//	@Param			Authorization	header	string					true	"Authentication header"
//	@Security		ApiKeyAuth
//	@Success		200	{object}	web.HttpData[string]
//	@Failure		502	{object}	web.HttpMsg
//	@Router			/user/changeMobile [post]
func (uc *UserCtrl) ChangeMobile(c *gin.Context) {
	var cmInput model.VerifyMobileInput
	if err := c.ShouldBindJSON(&cmInput); err != nil {
		web.BadRes(c, util.ErrInvalidArgument)
		uc.Log.Err(err).Send()
		return
	}
	if err := cmInput.Validate(); err != nil {
		web.BadRes(c, err)
		uc.Log.Error().Msg(err.LStr())
		return
	}
	mobile := cmInput.Number()

	if !ensureMobileNotReg(uc.Ctrl, c, mobile) {
		return
	}

	tu, ok := mustGetTu(uc.Ctrl, c)
	if !ok {
		return
	}
	if !recentToken(uc.Ctrl, c, tu) {
		return
	}
	user, ok := dbUserByID(uc.Ctrl, c, tu.BID)
	if !ok {
		return
	}

	if user.LastMVTime == 0 {
		err := util.ErrMobileNo
		web.BadRes(c, err)
		uc.Log.Error().Msg(err.LStr())
		return
	}
	if user.LastEVTime > 0 {
		err := util.ErrMobileByEmail(user.Email)
		web.BadRes(c, err)
		uc.Log.Error().Msg(err.LStr())
		return
	}
	if user.LastGVTime > 0 {
		err := util.ErrMobileByGa
		web.BadRes(c, err)
		uc.Log.Error().Msg(err.LStr())
		return
	}

	if !mobileCodeExists(uc.Ctrl, c, mobile, cmInput.Code) {
		return
	}

	user.Mobile = mobile
	user.LastMVTime = time.Now().Unix()
	if err := uc.DB.Model(&user).Updates(map[string]interface{}{
		"mobile":       user.Mobile,
		"last_mv_time": user.LastMVTime,
	}).Error; err != nil {
		web.BadRes(c, util.ErrDB)
		uc.Log.Err(err).Send()
		return
	}

	tu2 := user.ToTu(tu)
	AToken(uc.Ctrl, c, tu2)
}

// BindEmailBy godoc
//
//	@Summary		已有手机和/或谷歌认证，要求换电子邮箱。
//	@Description	提交新邮箱和新邮箱验证码，认证成功更换邮箱。前提是刚验证过手机和/或谷歌验证
//	@Tags			user
//	@Accept			json
//	@Product		json
//	@Param			email			body	model.VerifyEmailInput	true	"新邮箱和验证码"
//	@Param			Authorization	header	string					true	"Authentication header"
//	@Security		ApiKeyAuth
//	@Success		200	{object}	web.HttpData[string]
//	@Failure		502	{object}	web.HttpMsg
//	@Router			/user/bindEmailBy [post]
func (uc *UserCtrl) BindEmailBy(c *gin.Context) {
	var beInput model.VerifyEmailInput
	if err := c.ShouldBindJSON(&beInput); err != nil {
		web.BadRes(c, util.ErrInvalidArgument)
		uc.Log.Err(err).Send()
		return
	}
	if err := beInput.Validate(); err != nil {
		web.BadRes(c, err)
		uc.Log.Error().Msg(err.LStr())
		return
	}
	email := beInput.Email

	if !ensureEmailNotReg(uc.Ctrl, c, email) {
		return
	}

	tu, ok := mustGetTu(uc.Ctrl, c)
	if !ok {
		return
	}
	user, ok := dbUserByID(uc.Ctrl, c, tu.BID)
	if !ok {
		return
	}

	if user.LastMVTime == 0 && user.LastGVTime == 0 {
		err := util.ErrMobileGaNo
		web.BadRes(c, err)
		uc.Log.Error().Msg(err.LStr())
		return
	}
	if user.LastMVTime > 0 {
		if time.Now().Unix()-user.LastMVTime > 300 {
			err := util.ErrMobileFirst
			web.BadRes(c, err)
			uc.Log.Error().Msg(err.LStr())
			return
		}
	}
	if user.LastGVTime > 0 {
		if time.Now().Unix()-user.LastGVTime > 300 {
			err := util.ErrGaFirst
			web.BadRes(c, err)
			uc.Log.Error().Msg(err.LStr())
			return
		}
	}

	if !emailCodeExists(uc.Ctrl, c, email, beInput.Code) {
		return
	}

	user.Email = email
	user.LastEVTime = time.Now().Unix()
	if err := uc.DB.Model(&user).Updates(map[string]interface{}{
		"email":        user.Email,
		"last_ev_time": user.LastEVTime,
	}).Error; err != nil {
		web.BadRes(c, util.ErrDB)
		uc.Log.Err(err).Send()
		return
	}
	tu2 := user.ToTu(tu)
	AToken(uc.Ctrl, c, tu2)
}

// BindMobileBy godoc
//
//	@Summary		认证过邮箱和/或谷歌验证的情况下，要求换手机号码。提交新手机和验证码。
//	@Description	认证成功更换手机号。前提是刚验证过邮箱和/或谷歌验证
//	@Tags			user
//	@Accept			json
//	@Product		json
//	@Param			mobile			body	model.VerifyMobileInput	true	"新手机和验证码"
//	@Param			Authorization	header	string					true	"Authentication header"
//	@Security		ApiKeyAuth
//	@Success		200	{object}	web.HttpData[string]
//	@Failure		502	{object}	web.HttpMsg
//	@Router			/user/bindMobileBy [post]
func (uc *UserCtrl) BindMobileBy(c *gin.Context) {
	var bmInput model.VerifyMobileInput
	if err := c.ShouldBindJSON(&bmInput); err != nil {
		web.BadRes(c, util.ErrInvalidArgument)
		uc.Log.Err(err).Send()
		return
	}
	if err := bmInput.Validate(); err != nil {
		web.BadRes(c, err)
		uc.Log.Error().Msg(err.LStr())
		return
	}
	mobile := bmInput.Number()

	if !ensureMobileNotReg(uc.Ctrl, c, mobile) {
		return
	}

	tu, ok := mustGetTu(uc.Ctrl, c)
	if !ok {
		return
	}
	user, ok := dbUserByID(uc.Ctrl, c, tu.BID)
	if !ok {
		return
	}

	if user.LastEVTime == 0 && user.LastGVTime == 0 {
		err := util.ErrEmailGaNo
		web.BadRes(c, err)
		uc.Log.Error().Msg(err.LStr())
		return
	}
	if user.LastEVTime > 0 {
		if time.Now().Unix()-user.LastEVTime > 300 {
			err := util.ErrEmailFirst
			web.BadRes(c, err)
			uc.Log.Error().Msg(err.LStr())
			return
		}
	}
	if user.LastGVTime > 0 {
		if time.Now().Unix()-user.LastGVTime > 300 {
			err := util.ErrGaFirst
			web.BadRes(c, err)
			uc.Log.Error().Msg(err.LStr())
			return
		}
	}

	if !mobileCodeExists(uc.Ctrl, c, mobile, bmInput.Code) {
		return
	}

	user.Mobile = mobile
	user.LastMVTime = time.Now().Unix()

	if err := uc.DB.Model(&user).Updates(map[string]interface{}{
		"mobile":       mobile,
		"last_mv_time": user.LastMVTime,
	}).Error; err != nil {
		web.BadRes(c, util.ErrDB)
		uc.Log.Err(err).Send()
		return
	}
	tu2 := user.ToTu(tu)
	AToken(uc.Ctrl, c, tu2)
}

// UnbindGa godoc
//
//	@Summary		解绑谷歌验证
//	@Description	成功返回一个 token 包含当前用户更新后的信息
//	@Tags			user
//	@Accept			json
//	@Product		json
//	@Param			code			body	model.UnbindGAInput	true	"(email-code & | sms-code)"
//	@Param			Authorization	header	string				true	"Authentication header"
//	@Security		ApiKeyAuth
//	@Success		200	{object}	web.HttpData[string]
//	@Failure		502	{object}	web.HttpMsg
//	@Router			/user/unbindGa [post]
func (uc *UserCtrl) UnbindGa(c *gin.Context) {
	var ubInput model.UnbindGAInput
	if err := c.ShouldBindJSON(&ubInput); err != nil {
		web.BadRes(c, util.ErrInvalidArgument)
		uc.Log.Err(err).Send()
		return
	}
	if err := ubInput.Validate(); err != nil {
		web.BadRes(c, err)
		uc.Log.Error().Msg(err.LStr())
		return
	}

	tu, ok := mustGetTu(uc.Ctrl, c)
	if !ok {
		return
	}
	user, ok := dbUserByID(uc.Ctrl, c, tu.BID)
	if !ok {
		return
	}

	redisMKey := ""
	if user.LastMVTime > 0 {
		if redisMKey, ok = checkUserMobileCode(uc.Ctrl, c, &user, ubInput.MCode, true, true); !ok {
			return
		}
	}
	redisEKey := ""
	if user.LastEVTime > 0 {
		if redisEKey, ok = checkUserEmailCode(uc.Ctrl, c, &user, ubInput.ECode, true, true); !ok {
			return
		}
	}

	user.Ga = ""
	user.LastGVTime = 0
	if err := uc.DB.Model(&user).Updates(map[string]interface{}{
		"ga":           user.Ga,
		"last_gv_time": user.LastGVTime,
		"last_ev_time": user.LastEVTime,
		"last_mv_time": user.LastMVTime,
	}).Error; err != nil {
		web.BadRes(c, util.ErrDB)
		uc.Log.Err(err).Send()
		return
	}
	c0 := context.Background()
	if redisEKey != "" {
		_, _ = uc.RDB.Delete(c0, redisEKey)
	}
	if redisMKey != "" {
		_, _ = uc.RDB.Delete(c0, redisMKey)
	}

	tu2 := user.ToTu(tu)
	AToken(uc.Ctrl, c, tu2)
}

// DoResetPassword godoc
//
//	@Summary		认证完2FA的所有认证，直接提交新密码
//	@Description	刚认证完所有2FA认证，只需要提交新密码即可复位密码
//	@Tags			user
//	@Accept			json
//	@Product		json
//	@Param			newPassword		body	model.ResetPasswordInput	true	"新密码，需要符合密码规则"
//	@Param			Authorization	header	string						true	"Authentication header"
//	@Security		ApiKeyAuth
//	@Success		200	{object}	web.HttpData[string]
//	@Failure		502	{object}	web.HttpMsg
//	@Router			/user/do-reset-password [post]
func (uc *UserCtrl) DoResetPassword(c *gin.Context) {
	var rpInput model.ResetPasswordInput
	if err := c.ShouldBindJSON(&rpInput); err != nil {
		web.BadRes(c, util.ErrInvalidArgument)
		uc.Log.Err(err).Send()
		return
	}
	if data, err := base64.StdEncoding.DecodeString(rpInput.Password); err != nil {
		web.BadRes(c, util.ErrMsgDecode)
		uc.Log.Err(err).Send()
		return
	} else {
		rpInput.Password = string(data)
	}

	if err := rpInput.Validate(); err != nil {
		web.BadRes(c, err)
		uc.Log.Error().Msg(err.LStr())
		return
	}

	// 根据 token 拿到用户
	tu, ok := mustGetTu(uc.Ctrl, c)
	if !ok {
		return
	}
	user, ok := dbUserByID(uc.Ctrl, c, tu.BID)
	if !ok {
		return
	}

	if user.LastMVTime == 0 && user.LastEVTime == 0 {
		err := util.ErrEmailMobileNo
		web.BadRes(c, err)
		uc.Log.Error().Msg(err.LStr())
		return
	}
	if user.LastEVTime > 0 && time.Now().Unix()-user.LastEVTime > 300 {
		err := util.ErrEmailFirst
		web.BadRes(c, err)
		uc.Log.Error().Msg(err.LStr())
		return

	}
	if user.LastMVTime > 0 && time.Now().Unix()-user.LastMVTime > 300 {
		err := util.ErrMobileFirst
		web.BadRes(c, err)
		uc.Log.Error().Msg(err.LStr())
		return

	}
	if user.LastGVTime > 0 && time.Now().Unix()-user.LastGVTime > 300 {
		err := util.ErrGaFirst
		web.BadRes(c, err)
		uc.Log.Error().Msg(err.LStr())
		return
	}

	hashedPass, er := util.HashPassword(rpInput.Password)
	if er != nil {
		web.BadRes(c, er)
		uc.Log.Error().Msg(er.LStr())
		return
	}
	user.Password = hashedPass
	if err := uc.DB.Model(&user).Update("password", hashedPass).Error; err != nil {
		web.BadRes(c, util.ErrDB)
		uc.Log.Err(err).Send()
		return
	}
	web.GoodResp(c, "Ok")
}

// ChangeNick godoc
//
//	@Summary		修改自己的昵称
//	@Description	返回带有新昵称的token
//	@Tags			user
//	@Accept			json
//	@Product		json
//	@Param			Nick			body	model.NickInput	true	"新昵称"
//	@Param			Authorization	header	string			true	"Authentication header"
//	@Security		ApiKeyAuth
//	@Success		200	{object}	web.HttpData[string]
//	@Failure		502	{object}	web.HttpMsg
//	@Router			/user/change-nick [post]
func (uc *UserCtrl) ChangeNick(c *gin.Context) {
	var ni model.NickInput
	if err := c.ShouldBindJSON(&ni); err != nil {
		web.BadRes(c, util.ErrInvalidArgument)
		uc.Log.Err(err).Send()
		return
	}
	if !ensureNickNotReg(uc.Ctrl, c, ni.Nick) {
		return
	}

	// 根据 token 拿到用户
	tu, ok := mustGetTu(uc.Ctrl, c)
	if !ok {
		return
	}
	user, ok := dbUserByID(uc.Ctrl, c, tu.BID)
	if !ok {
		return
	}

	if err := uc.DB.Model(&user).Update("nick", ni.Nick).Error; err != nil {
		web.BadRes(c, util.ErrDB)
		uc.Log.Err(err).Send()
		return
	}

	user.Nick = ni.Nick
	tu2 := user.ToTu(tu)
	tu2.Auth2 = tu.Auth2
	AToken(uc.Ctrl, c, tu2)
}

// QueryFreeNickByID godoc
//
//	@Summary		企业管理员添加成员前根据uid查询用户昵称
//	@Description	如果用户不属于其他企业则返回昵称，否则返回已占用错误
//	@Tags			user
//	@Accept			json
//	@Product		json
//	@Param			uid				body	model.UidInput	true	"uid"
//	@Param			Authorization	header	string			true	"Authentication header"
//	@Security		ApiKeyAuth
//	@Success		200	{object}	web.HttpData[string]
//	@Failure		502	{object}	web.HttpMsg
//	@Router			/user/query-free-nick-by-id [post]
func (uc *UserCtrl) QueryFreeNickByID(c *gin.Context) {
	var ni model.UidInput
	if err := c.ShouldBindJSON(&ni); err != nil {
		web.BadRes(c, util.ErrInvalidArgument)
		uc.Log.Err(err).Send()
		return
	}

	// 根据 token 拿到用户
	tu, ok := mustGetTu(uc.Ctrl, c)
	if !ok {
		return
	}
	if !tu.Admin {
		err := util.ErrUnAuthed
		web.BadRes(c, err)
		uc.Log.Error().Interface("tu", tu).Msg(err.LStr())
		return
	}
	objUser, ok := dbUserByID(uc.Ctrl, c, ni.Uid)
	if !ok {
		return
	}

	if objUser.Fid != "" {
		err := util.ErrKycUser2(objUser.Nick)
		if objUser.Fid == tu.BID {
			err = util.ErrKycUser1(objUser.Nick)
		}
		web.BadRes(c, err)
		uc.Log.Error().Interface("user", objUser).Msg(err.LStr())
		return
	}
	web.GoodResp(c, objUser.Nick)
}

// AddFirmUserByID godoc
//
//	@Summary		企业管理员添加成员
//	@Description	如果用户不属于其他企业则返回Ok，否则返回已占用错误
//	@Tags			user
//	@Accept			json
//	@Product		json
//	@Param			uid				body	model.UidInput	true	"uid"
//	@Param			Authorization	header	string			true	"Authentication header"
//	@Security		ApiKeyAuth
//	@Success		200	{object}	web.HttpData[string]
//	@Failure		502	{object}	web.HttpMsg
//	@Router			/user/add-firm-user-by-id [post]
func (uc *UserCtrl) AddFirmUserByID(c *gin.Context) {
	var ni model.UidInput
	if err := c.ShouldBindJSON(&ni); err != nil {
		web.BadRes(c, util.ErrInvalidArgument)
		uc.Log.Err(err).Send()
		return
	}

	// 根据 token 拿到用户
	tu, ok := mustGetTu(uc.Ctrl, c)
	if !ok {
		return
	}
	if !tu.Admin {
		err := util.ErrUnAuthed
		web.BadRes(c, err)
		uc.Log.Error().Interface("tu", tu).Msg(err.LStr())
		return
	}
	if !dbFirmAddUser(uc.Ctrl, c, ni.Uid, tu.BID) {
		return
	}
	web.GoodResp(c, "Ok")
}

// DelFirmUserByID godoc
//
//	@Summary		企业管理员去除成员
//	@Description	如果去除成功则返回Ok，否则返回错误（已安排在多签业务场景的，必须清理完多签委派，才能
//	@Tags			user
//	@Accept			json
//	@Product		json
//	@Param			uid				body	model.UidInput	true	"uid"
//	@Param			Authorization	header	string			true	"Authentication header"
//	@Security		ApiKeyAuth
//	@Success		200	{object}	web.HttpData[string]
//	@Failure		502	{object}	web.HttpMsg
//	@Router			/user/del-firm-user-by-id [post]
func (uc *UserCtrl) DelFirmUserByID(c *gin.Context) {
	var ni model.UidInput
	if err := c.ShouldBindJSON(&ni); err != nil {
		web.BadRes(c, util.ErrInvalidArgument)
		uc.Log.Err(err).Send()
		return
	}

	// 根据 token 拿到用户
	tu, ok := mustGetTu(uc.Ctrl, c)
	if !ok {
		return
	}
	if !tu.Admin {
		err := util.ErrUnAuthed
		web.BadRes(c, err)
		uc.Log.Error().Interface("tu", tu).Msg(err.LStr())
		return
	}

	// TODO resty 狗哥多签
	if !dbFirmDelUser(uc.Ctrl, c, ni.Uid, tu.BID) {
		return
	}
	web.GoodResp(c, "Ok")
}

// ListFirmUser godoc
//
//	@Summary		企业管理员查看自家成员
//	@Description	返回列表
//	@Tags			user
//	@Accept			json
//	@Product		json
//	@Param			Authorization	header	string	true	"Authentication header"
//	@Security		ApiKeyAuth
//	@Success		200	{object}	web.HttpData[[]model.FirmUser]
//	@Failure		502	{object}	web.HttpMsg
//	@Router			/user/list-firm-user [get]
func (uc *UserCtrl) ListFirmUser(c *gin.Context) {
	// 根据 token 拿到用户
	tu, ok := mustGetTu(uc.Ctrl, c)
	if !ok {
		return
	}
	if !tu.Admin {
		err := util.ErrUnAuthed
		web.BadRes(c, err)
		uc.Log.Error().Interface("tu", tu).Msg(err.LStr())
		return
	}
	users, ok := dbFirmListUser(uc.Ctrl, c, tu.BID)
	if !ok {
		return
	}
	web.GoodResp(c, users)
}

// ValidateAndSign godoc
//
//	@Summary		重要操作之前验证三码，签名操作信息
//	@Description	返回签发的信息
//	@Tags			user
//	@Accept			json
//	@Product		json
//	@Param			preSign			body	model.PreSignInput	true	"code & data"
//	@Param			Authorization	header	string				true	"Authentication header"
//	@Security		ApiKeyAuth
//	@Success		200	{object}	web.HttpData[string]
//	@Failure		502	{object}	web.HttpMsg
//	@Router			/user/validate-n-sign [post]
func (uc *UserCtrl) ValidateAndSign(c *gin.Context) {
	var ps model.PreSignInput
	if err := c.ShouldBindJSON(&ps); err != nil {
		web.BadRes(c, util.ErrInvalidArgument)
		uc.Log.Err(err).Send()
		return
	}

	// 根据 token 拿到用户
	tu, ok := mustGetTu(uc.Ctrl, c)
	if !ok {
		return
	}

	user, ok := dbUserByID(uc.Ctrl, c, tu.BID)
	if !ok {
		return
	}

	if user.Mobile == "" {
		err := util.ErrMobileNo
		web.BadRes(c, err)
		uc.Log.Error().Msg(err.LStr())
		return
	}
	if user.Email == "" {
		err := util.ErrEmailNo
		web.BadRes(c, err)
		uc.Log.Error().Msg(err.LStr())
		return
	}
	if user.Ga == "" {
		err := util.ErrGaNew
		web.BadRes(c, err)
		uc.Log.Error().Msg(err.LStr())
		return
	}
	redisMKey, ok := checkUserMobileCode(uc.Ctrl, c, &user, ps.MCode, true, true)
	if !ok {
		return
	}
	redisEKey, ok := checkUserEmailCode(uc.Ctrl, c, &user, ps.ECode, true, true)
	if !ok {
		return
	}
	v, er := util.ValidateTOTP(ps.GCode, user.Ga)
	if er != nil || !v {
		web.BadRes(c, util.ErrGaInvalid)
		uc.Log.Error().Msg(er.LStr())
		return
	}

	now := time.Now().Unix()
	if err := uc.DB.Model(&user).Updates(map[string]interface{}{
		"last_ev_time": now,
		"last_mv_time": now,
		"last_gv_time": now,
	}).Error; err != nil {
		web.BadRes(c, util.ErrDB)
		uc.Log.Err(err).Send()
		return
	}
	_, _ = uc.RDB.Delete(context.Background(), redisEKey)
	_, _ = uc.RDB.Delete(context.Background(), redisMKey)

	pem, _ := base64.StdEncoding.DecodeString(conf.Conf.AccessTokenPrivateKey)
	if pk, err := rsa0.ParseRsaPrivateKeyFromPemStr(string(pem)); err != nil {
		web.BadRes(c, util.ErrBcryptHash)
		uc.Log.Err(err).Send()
		return
	} else {
		data, err := util.RsaSign(ps.Data, pk)
		if err != nil {
			web.BadRes(c, err)
			uc.Log.Error().Msg(err.LStr())
			return
		}
		web.GoodResp(c, data)
	}
}

type safeResult struct {
	Code    int64
	Data    interface{}
	Msg     string
	Success bool
}

func (uc *UserCtrl) isSafeMember(r *Ctrl, c *gin.Context, uid string) (ret, ok bool) {
	var result safeResult
	resp, err := uc.SafeCli.Cli.R().
		SetHeader("uid", uid).
		SetResult(&result).
		Get("/v1/in/sig/member")
	if err != nil {
		web.BadRes(c, util.ErrSafeSvr)
		r.Log.Err(err).Send()
		return
	}
	if resp.StatusCode() != http.StatusOK || !result.Success {
		web.BadRes(c, util.ErrSafeSvr)
		r.Log.Error().Int("safeSvr.status", resp.StatusCode()).Send()
		return
	}
	ret, ok = result.Data.(bool)
	if !ok {
		web.BadRes(c, util.ErrSafeSvr)
		r.Log.Error().Interface("result", result).Send()
	}
	return
}
