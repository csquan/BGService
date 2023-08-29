package ctrl

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"time"
	"user/pkg/conf"
	"user/pkg/db"
	"user/pkg/log"
	"user/pkg/model"
	"user/pkg/msg"
	"user/pkg/redis"
	"user/pkg/util"
	"user/pkg/util/ecies"
	"user/pkg/web"
)

type PubCtrl struct {
	*Ctrl
	Pk     *ecies.PrivateKey
	KeyCli *web.MyInnerCli
	Prd    *kafka.Producer
}

func NewPubCtrl() PubCtrl {
	rdb := redis.NewStore("go:user:controller:")
	logger := log.C.Logger().With().Str("ctrl", "pub").Logger()
	pk, er := ecies.PrivateFromString(conf.Conf.KycPrivateKey)
	if er != nil {
		logger.Error().Msg(er.LStr())
		panic(er)
	}
	keyCli, err := web.NewCli(conf.Conf.KeySvrUrl, conf.Conf.KeySvrPubKey)
	if err != nil {
		logger.Error().Msg(err.LStr())
		panic(err)
	}
	producer := msg.NewKafkaProducer()
	return PubCtrl{
		Ctrl:   &Ctrl{DB: db.DB, RDB: rdb, Log: &logger},
		Pk:     pk,
		KeyCli: keyCli,
		Prd:    producer,
	}
}

// MyEmail godoc
//
//	@Summary		首次提交电子邮箱，要求发起认证真实性
//	@Description	后台生成 6 位随机数字邮件发送并保存在 Redis，5分钟过期
//	@Tags			pub
//	@Accept			json
//	@Product		json
//	@Param			email	body		model.EmailInput	true	"在 email 字段填入你的邮箱"
//	@Success		200		{object}	web.HttpData[string]
//	@Failure		502		{object}	web.HttpMsg
//	@Router			/pub/myEmail [post]
func (pc *PubCtrl) MyEmail(c *gin.Context) {
	var emailInput model.EmailInput
	if err := c.ShouldBindJSON(&emailInput); err != nil {
		web.BadRes(c, util.ErrInvalidArgument)
		pc.Log.Err(err).Send()
		return
	}
	if err := emailInput.Validate(); err != nil {
		web.BadRes(c, err)
		pc.Log.Warn().Msg(err.LStr())
		return
	}
	email := emailInput.Email

	// 首先保证数据库中不重复
	if !ensureEmailNotReg(pc.Ctrl, c, email) {
		return
	}
	sendEmailNStore(pc.Ctrl, c, email)
}

// MyMobile godoc
//
//	@Summary		首次提交手机号，要求发起认证真实性
//	@Description	后台生成 6 位随机数字短信发送给用户并保存在 Redis，5分钟过期
//	@Tags			pub
//	@Accept			json
//	@Product		json
//	@Param			mobile	body		model.MobileInput	true	"在 mobile 字段填入你的手机号，带国家字冠"
//	@Success		200		{object}	web.HttpData[string]
//	@Failure		502		{object}	web.HttpMsg
//	@Router			/pub/myMobile [post]
func (pc *PubCtrl) MyMobile(c *gin.Context) {
	var mobileInput model.MobileInput
	if err := c.ShouldBindJSON(&mobileInput); err != nil {
		web.BadRes(c, util.ErrInvalidArgument)
		pc.Log.Err(err).Send()
		return
	}
	if err := mobileInput.Validate(); err != nil {
		web.BadRes(c, err)
		pc.Log.Warn().Msg(err.LStr())
		return
	}
	mobile := mobileInput.Number()

	// 首先保证不重复
	if !ensureMobileNotReg(pc.Ctrl, c, mobile) {
		return
	}

	sendMobileNStore(pc.Ctrl, c, mobile)
}

// VerifyEmail godoc
//
//	@Summary		校验邮件验证码，成功证明是本人邮箱
//	@Description	成功返回一个 token 用于保存当前用户的所有信息
//	@Tags			pub
//	@Accept			json
//	@Product		json
//	@Param			input	body		model.VerifyEmailInput	true	"输入你的邮箱和收到的验证码"
//	@Success		200		{object}	web.HttpData[string]
//	@Failure		502		{object}	web.HttpMsg
//	@Router			/pub/verifyEmail [post]
func (pc *PubCtrl) VerifyEmail(c *gin.Context) {
	var verifyEmailInput model.VerifyEmailInput
	if err := c.ShouldBindJSON(&verifyEmailInput); err != nil {
		web.BadRes(c, util.ErrInvalidArgument)
		pc.Log.Err(err).Send()
		return
	}
	if err := verifyEmailInput.Validate(); err != nil {
		web.BadRes(c, err)
		pc.Log.Error().Msg(err.LStr())
		return
	}
	email := verifyEmailInput.Email
	code := verifyEmailInput.Code

	if !emailCodeExists(pc.Ctrl, c, email, code) {
		return
	}

	tu := model.TokenUser{
		Email:     email,
		EAuthTime: time.Now().Unix(),
		Auth2:     false,
	}
	aToken(pc.Ctrl, c, &tu)
}

// VerifyMobile godoc
//
//	@Summary		校验手机验证码，成功证明是本人手机
//	@Description	成功返回一个 token 可以用于保存当前用户的所有信息
//	@Tags			pub
//	@Accept			json
//	@Product		json
//	@Param			input	body		model.VerifyMobileInput	true	"输入你的手机号和收到的验证码"
//	@Success		200		{object}	web.HttpData[string]
//	@Failure		502		{object}	web.HttpMsg
//	@Router			/pub/verifyMobile [post]
func (pc *PubCtrl) VerifyMobile(c *gin.Context) {
	var verifyMobileInput model.VerifyMobileInput
	if err := c.ShouldBindJSON(&verifyMobileInput); err != nil {
		web.BadRes(c, util.ErrInvalidArgument)
		pc.Log.Err(err).Send()
		return
	}
	if err := verifyMobileInput.Validate(); err != nil {
		web.BadRes(c, err)
		pc.Log.Warn().Msg(err.LStr())
		return
	}
	mobile := verifyMobileInput.Number()
	code := verifyMobileInput.Code

	if ok := mobileCodeExists(pc.Ctrl, c, mobile, code); !ok {
		return
	}

	tu := model.TokenUser{
		Mobile:    mobile,
		MAuthTime: time.Now().Unix(),
		Auth2:     false,
	}
	aToken(pc.Ctrl, c, &tu)
}

// SignIn godoc
//
//	@Summary		登录系统
//	@Description	成功返回一个 token 可以用于保存当前用户的所有信息
//	@Tags			pub
//	@Accept			json
//	@Product		json
//	@Param			user	body		model.SignInInput	true	"user information"
//	@Success		200		{object}	web.HttpData[string]
//	@Failure		502		{object}	web.HttpMsg
//	@Router			/pub/signIn [post]
func (pc *PubCtrl) SignIn(c *gin.Context) {
	var signInInput model.SignInInput
	if err := c.ShouldBindJSON(&signInInput); err != nil {
		web.BadRes(c, util.ErrInvalidArgument)
		pc.Log.Err(err).Send()
		return
	}
	if data, err := base64.StdEncoding.DecodeString(signInInput.Password); err != nil {
		web.BadRes(c, util.ErrMsgDecode)
		pc.Log.Err(err).Send()
		return
	} else {
		signInInput.Password = string(data)
	}

	var user model.User
	ok := false
	user, ok = dbUser(pc.Ctrl, c, signInInput.User)
	if !ok {
		return
	}

	if user.Status < 0 {
		web.BadRes(c, util.ErrUserBan)
		pc.Log.Error().Int8("user.status", user.Status).Send()
		return
	}

	if err := util.VerifyPassword(user.Password, signInInput.Password); err != nil {
		web.BadRes(c, err)
		pc.Log.Error().Msg(err.LStr())
		return
	}
	aa, ok := dbAppAddr(pc.Ctrl, c, user.BID, signInInput.AppID)
	if !ok {
		return
	}
	tu := user.ToTU(&aa)
	tu.Auth2 = false
	if !rCookie(pc.Ctrl, c, tu) {
		return
	}
	aToken(pc.Ctrl, c, tu)
}

// Refresh godoc
//
//	@Summary		刷新登录
//	@Description	因为token有效期并不是很长，失效后在可刷新周期内，刷新一下可以获得一个新token
//	@Tags			pub
//	@Accept			json
//	@Product		json
//	@Success		200	{object}	web.HttpData[string]
//	@Failure		502	{object}	web.HttpMsg
//	@Router			/pub/refresh [get]
func (pc *PubCtrl) Refresh(c *gin.Context) {
	tu, ok := shouldGetRTu(pc.Ctrl, c)
	if !ok {
		return
	}
	if !tuDbOk(pc.Ctrl, c, tu) {
		return
	}
	if !rCookie(pc.Ctrl, c, tu) {
		return
	}
	AToken(pc.Ctrl, c, tu)
}

// ForgetE godoc
//
//	@Summary		忘记密码，从邮箱开始找回。
//	@Description	需要已经认证过的邮箱。返回可供认证的信息
//	@Tags			pub
//	@Accept			json
//	@Param			email	body	model.EmailInput	true	"在 email 字段填入你的邮箱"
//	@Product		json
//	@Success		200	{object}	web.HttpData[model.ForgetResponse]
//	@Failure		502	{object}	web.HttpMsg
//	@Router			/pub/forget-e [post]
func (pc *PubCtrl) ForgetE(c *gin.Context) {
	var emailInput model.EmailInput
	if err := c.ShouldBindJSON(&emailInput); err != nil {
		web.BadRes(c, util.ErrInvalidArgument)
		pc.Log.Err(err).Send()
		return
	}
	if err := emailInput.Validate(); err != nil {
		web.BadRes(c, err)
		pc.Log.Error().Msg(err.LStr())
		return
	}
	email := emailInput.Email

	user, ok := dbUserByEmail(pc.Ctrl, c, email)
	if !ok {
		return
	}

	res := model.ForgetResponse{
		MaskedEmail:  user.MaskedEmail(),
		MaskedMobile: user.MaskedMobile(),
		E:            user.LastEVTime > 0,
		M:            user.LastMVTime > 0,
		G:            user.LastGVTime > 0,
	}
	web.GoodResp(c, res)
}

// ForgetM godoc
//
//	@Summary		忘记密码，从手机开始找回
//	@Description	需要手机已经认证过。返回可供验证的信息
//	@Tags			pub
//	@Accept			json
//	@Param			mobile	body	model.MobileInput	true	"在 mobile 字段填入你的手机号，带国家字冠"
//	@Product		json
//	@Success		200	{object}	web.HttpData[model.ForgetResponse]
//	@Failure		502	{object}	web.HttpMsg
//	@Router			/pub/forget-m [post]
func (pc *PubCtrl) ForgetM(c *gin.Context) {
	var mobileInput model.MobileInput
	if err := c.ShouldBindJSON(&mobileInput); err != nil {
		web.BadRes(c, util.ErrInvalidArgument)
		pc.Log.Err(err).Send()
		return
	}
	if err := mobileInput.Validate(); err != nil {
		web.BadRes(c, err)
		pc.Log.Error().Msg(err.LStr())
		return
	}
	mobile := mobileInput.Number()

	user, ok := dbUserByMobile(pc.Ctrl, c, mobile)
	if !ok {
		return
	}
	res := model.ForgetResponse{
		MaskedEmail:  user.MaskedEmail(),
		MaskedMobile: user.MaskedMobile(),
		E:            user.LastEVTime > 0,
		M:            user.LastMVTime > 0,
		G:            user.LastGVTime > 0,
	}
	web.GoodResp(c, res)
}

// ForgetSendEmailByEmail godoc
//
//	@Summary		知道邮箱，找回密码时，请求邮箱验证码
//	@Description	需要邮箱已经认证过。返回Ok
//	@Tags			pub
//	@Accept			json
//	@Param			email	body	model.EmailInput	true	"在 email 字段填入你的邮箱"
//	@Product		json
//	@Success		200	{object}	web.HttpData[string]
//	@Failure		502	{object}	web.HttpMsg
//	@Router			/pub/forget-send-email-by-email [post]
func (pc *PubCtrl) ForgetSendEmailByEmail(c *gin.Context) {
	var emailInput model.EmailInput
	if err := c.ShouldBindJSON(&emailInput); err != nil {
		web.BadRes(c, util.ErrInvalidArgument)
		pc.Log.Err(err).Send()
		return
	}
	if err := emailInput.Validate(); err != nil {
		web.BadRes(c, err)
		pc.Log.Error().Msg(err.LStr())
		return
	}
	email := emailInput.Email

	if !ensureEmailExists(pc.Ctrl, c, email) {
		return
	}
	sendEmailNStore(pc.Ctrl, c, email)
}

// ForgetSendSMSByEmail godoc
//
//	@Summary		知道邮箱，找回密码时，请求手机验证码
//	@Description	需要邮箱和手机都已经认证过。返回Ok
//	@Tags			pub
//	@Accept			json
//	@Param			email	body	model.EmailInput	true	"在 email 字段填入你的邮箱"
//	@Product		json
//	@Success		200	{object}	web.HttpData[string]
//	@Failure		502	{object}	web.HttpMsg
//	@Router			/pub/forget-send-sms-by-email [post]
func (pc *PubCtrl) ForgetSendSMSByEmail(c *gin.Context) {
	var emailInput model.EmailInput
	if err := c.ShouldBindJSON(&emailInput); err != nil {
		web.BadRes(c, util.ErrInvalidArgument)
		pc.Log.Err(err).Send()
		return
	}
	if err := emailInput.Validate(); err != nil {
		web.BadRes(c, err)
		pc.Log.Error().Msg(err.LStr())
		return
	}
	email := emailInput.Email

	// 查询用户是否存在
	var user model.User
	ok := false
	if user, ok = dbUserByEmail(pc.Ctrl, c, email); !ok {
		return
	}
	if !ensureUserHasMobile(pc.Ctrl, c, &user) {
		return
	}

	sendMobileNStore(pc.Ctrl, c, user.Mobile)
}

// ForgetSendEmailByMobile godoc
//
//	@Summary		知道手机，找回密码时，请求邮箱验证码
//	@Description	需要手机和邮箱都已认证过。返回Ok
//	@Tags			pub
//	@Accept			json
//	@Param			mobile	body	model.MobileInput	true	"在 mobile 字段填入你的手机号，带国家字冠"
//	@Product		json
//	@Success		200	{object}	web.HttpData[string]
//	@Failure		502	{object}	web.HttpMsg
//	@Router			/pub/forget-send-email-by-mobile [post]
func (pc *PubCtrl) ForgetSendEmailByMobile(c *gin.Context) {
	var mobileInput model.MobileInput
	if err := c.ShouldBindJSON(&mobileInput); err != nil {
		web.BadRes(c, util.ErrInvalidArgument)
		pc.Log.Err(err).Send()
		return
	}
	if err := mobileInput.Validate(); err != nil {
		web.BadRes(c, err)
		pc.Log.Error().Msg(err.LStr())
		return
	}
	mobile := mobileInput.Mobile

	// 查询用户是否存在
	var user model.User
	ok := false
	if user, ok = dbUserByMobile(pc.Ctrl, c, mobile); !ok {
		return
	}
	if !ensureUserHasEmail(pc.Ctrl, c, &user) {
		return
	}
	sendEmailNStore(pc.Ctrl, c, user.Email)
}

// ForgetSendSMSByMobile godoc
//
//	@Summary		知道手机号码，找回密码时，请求手机验证码
//	@Description	需要手机已经认证过。返回Ok
//	@Tags			pub
//	@Accept			json
//	@Param			mobile	body	model.MobileInput	true	"在 mobile 字段填入你的手机号，带国家字冠"
//	@Product		json
//	@Success		200	{object}	web.HttpData[string]
//	@Failure		502	{object}	web.HttpMsg
//	@Router			/pub/forget-send-sms-by-mobile [post]
func (pc *PubCtrl) ForgetSendSMSByMobile(c *gin.Context) {
	var mobileInput model.MobileInput
	if err := c.ShouldBindJSON(&mobileInput); err != nil {
		web.BadRes(c, util.ErrInvalidArgument)
		pc.Log.Err(err).Send()
		return
	}
	if err := mobileInput.Validate(); err != nil {
		web.BadRes(c, err)
		pc.Log.Error().Msg(err.LStr())
		return
	}
	mobile := mobileInput.Mobile

	if !ensureMobileExists(pc.Ctrl, c, mobile) {
		return
	}

	sendMobileNStore(pc.Ctrl, c, mobile)
}

// PresetPassword godoc
//
//	@Summary		通过收到的各种验证吗，提交验证，颁发token
//	@Description	返回 Token
//	@Tags			pub
//	@Accept			json
//	@Param			input	body	model.PresetPasswordInput	true	"复位密码所需的信息"
//	@Product		json
//	@Success		200	{object}	web.HttpData[string]
//	@Failure		502	{object}	web.HttpMsg
//	@Router			/pub/preset-password [post]
func (pc *PubCtrl) PresetPassword(c *gin.Context) {
	var resetInput model.PresetPasswordInput
	if err := c.ShouldBindJSON(&resetInput); err != nil {
		web.BadRes(c, util.ErrInvalidArgument)
		pc.Log.Err(err).Send()
		return
	}
	if err := resetInput.Validate(); err != nil {
		web.BadRes(c, err)
		pc.Log.Error().Msg(err.LStr())
		return
	}

	var user model.User
	ok := false
	if user, ok = dbUser(pc.Ctrl, c, resetInput.User); !ok {
		return
	}

	redisEKey := ""
	redisMKey := ""
	c0 := context.Background()
	if user.LastEVTime > 0 {
		if redisEKey, ok = checkUserEmailCode(pc.Ctrl, c, &user, resetInput.ECode, true, false); !ok {
			return
		}
	}
	if user.LastMVTime > 0 {
		if redisMKey, ok = checkUserMobileCode(pc.Ctrl, c, &user, resetInput.MCode, true, false); !ok {
			return
		}
	}
	if user.LastGVTime > 0 {
		if !checkUserGCode(pc.Ctrl, c, &user, resetInput.GCode) {
			return
		}
	}
	err := pc.DB.Model(&user).Updates(map[string]interface{}{
		"last_ev_time": user.LastEVTime,
		"last_mv_time": user.LastMVTime,
		"last_gv_time": user.LastGVTime,
	}).Error
	if err != nil {
		web.BadRes(c, util.ErrDB)
		pc.Log.Err(err).Send()
		return
	}
	if redisEKey != "" {
		_, _ = pc.RDB.Delete(c0, redisEKey)
	}
	if redisMKey != "" {
		_, _ = pc.RDB.Delete(c0, redisMKey)
	}
	aa, ok := dbAppAddr(pc.Ctrl, c, user.BID, resetInput.AppID)
	if !ok {
		return
	}
	tu := user.ToTU(&aa)
	tu2 := web.GetTokenUser(c)
	if tu2 != nil && tu2.BID == tu.BID && tu2.AppID == resetInput.AppID {
		tu.Auth2 = tu2.Auth2
		tu.SK = tu2.SK
	} else {
		tu.Auth2 = false
		tu.ReSk()
	}

	AToken(pc.Ctrl, c, tu)
}

// ValidateEmail godoc
//
//	@Summary		提交SignUp之前，验证输入邮箱是否可用
//	@Description	只验证邮箱状态，不发送验证码
//	@Tags			pub
//	@Accept			json
//	@Product		json
//	@Param			EmailInput	body		model.EmailInput	true	"邮箱"
//	@Success		200			{object}	web.HttpData[string]
//	@Failure		502			{object}	web.HttpMsg
//	@Router			/pub/validate-email [post]
func (pc *PubCtrl) ValidateEmail(c *gin.Context) {
	var ei model.EmailInput
	if err := c.ShouldBindJSON(&ei); err != nil {
		web.BadRes(c, util.ErrInvalidArgument)
		pc.Log.Err(err).Send()
		return
	}
	if err := ei.Validate(); err != nil {
		web.BadRes(c, err)
		pc.Log.Error().Msg(err.LStr())
		return
	}

	if !ensureEmailNotReg(pc.Ctrl, c, ei.Email) {
		return
	}

	web.GoodResp(c, "Ok")
}

// ValidateMobile godoc
//
//	@Summary		提交SignUp之前，已验证邮箱的情况下，验证输入手机号是否可用
//	@Description	只验手机状态，不发送验证码
//	@Tags			pub
//	@Accept			json
//	@Product		json
//	@Param			MobileInput	body		model.MobileInput	true	"手机"
//	@Success		200			{object}	web.HttpData[string]
//	@Failure		502			{object}	web.HttpMsg
//	@Router			/pub/validate-mobile [post]
func (pc *PubCtrl) ValidateMobile(c *gin.Context) {
	var mi model.MobileInput
	if err := c.ShouldBindJSON(&mi); err != nil {
		web.BadRes(c, util.ErrInvalidArgument)
		pc.Log.Err(err).Send()
		return
	}
	if err := mi.Validate(); err != nil {
		web.BadRes(c, err)
		pc.Log.Error().Msg(err.LStr())
		return
	}
	mobile := mi.Number()

	if !ensureMobileNotReg(pc.Ctrl, c, mobile) {
		return
	}

	web.GoodResp(c, "Ok")
}

// SignOut godoc
//
//	@Summary		退出系统
//	@Description	退出系统
//	@Tags			pub
//	@Accept			json
//	@Product		json
//	@Param			Authorization	header	string	true	"Authentication header"
//	@Security		ApiKeyAuth
//	@Success		200	{object}	web.HttpData[string]
//	@Failure		502	{object}	web.HttpMsg
//	@Router			/pub/signOut [post]
func (pc *PubCtrl) SignOut(c *gin.Context) {
	c.SetCookie("refresh_token", "", -1, "/", "", true, true)
	web.GoodResp(c, "Ok")
}

// PubKey godoc
//
//	@Summary		获取 PublicKey
//	@Description	用于校验 token 的有效性
//	@Tags			pub
//	@Accept			json
//	@Product		json
//	@Success		200	{object}	web.HttpData[string]
//	@Failure		502	{object}	web.HttpMsg
//	@Router			/pub/pubKey [get]
func (pc *PubCtrl) PubKey(c *gin.Context) {
	web.GoodResp(c, conf.Conf.AccessTokenPublicKey)
}

// ResEmail godoc
//
//	@Summary		已登录用户再次验证电子邮箱，发起认证是否本人操作
//	@Description	带token操作，无需主动输入参数。后台生成 6 位随机数字邮件发送并保存在 Redis，5分钟过期
//	@Tags			pub
//	@Accept			json
//	@Product		json
//	@Param			Authorization	header	string	true	"Authentication header"
//	@Security		ApiKeyAuth
//	@Success		200	{object}	web.HttpData[string]
//	@Failure		502	{object}	web.HttpMsg
//	@Router			/pub/resEmail [post]
func (pc *PubCtrl) ResEmail(c *gin.Context) {
	tu, ok := shouldGetTu(pc.Ctrl, c)
	if !ok {
		return
	}

	user, ok := dbUserByID(pc.Ctrl, c, tu.BID)
	if !ok {
		return
	}

	if user.Email == "" {
		err := util.ErrEmailNo
		web.BadRes(c, err)
		pc.Log.Error().Msg(err.LStr())
		return
	}

	sendEmailNStore(pc.Ctrl, c, user.Email)
}

// ResMobile godoc
//
//	@Summary		提交验证手机号，验证是否本人操作
//	@Description	已经登录情况下操作，无需其他参数。后台生成 6 位随机数字短信发送给用户并保存在 Redis，5分钟过期
//	@Tags			pub
//	@Accept			json
//	@Product		json
//	@Param			Authorization	header	string	true	"Authentication header"
//	@Security		ApiKeyAuth
//	@Success		200	{object}	web.HttpData[string]
//	@Failure		502	{object}	web.HttpMsg
//	@Router			/pub/resMobile [post]
func (pc *PubCtrl) ResMobile(c *gin.Context) {
	tu, ok := shouldGetTu(pc.Ctrl, c)
	if !ok {
		return
	}

	user, ok := dbUserByID(pc.Ctrl, c, tu.BID)
	if !ok {
		return
	}

	if user.Mobile == "" {
		err := util.ErrMobileNo
		web.BadRes(c, err)
		pc.Log.Error().Msg(err.LStr())
		return
	}

	sendMobileNStore(pc.Ctrl, c, user.Mobile)
}

// RevEmail godoc
//
//	@Summary		已登录状态下校验邮件验证码，证明是本人操作
//	@Description	成功返回一个新 token
//	@Tags			pub
//	@Accept			json
//	@Product		json
//	@Param			Authorization	header	string			true	"Authentication header"
//	@Param			input			body	model.CodeInput	true	"输入你收到的验证码"
//	@Security		ApiKeyAuth
//	@Success		200	{object}	web.HttpData[string]
//	@Failure		502	{object}	web.HttpMsg
//	@Router			/pub/revEmail [post]
func (pc *PubCtrl) RevEmail(c *gin.Context) {
	var reVerifyInput model.CodeInput
	if err := c.ShouldBindJSON(&reVerifyInput); err != nil {
		web.BadRes(c, util.ErrInvalidArgument)
		pc.Log.Err(err).Send()
		return
	}
	code := reVerifyInput.Code

	tu, ok := shouldGetTu(pc.Ctrl, c)
	if !ok {
		return
	}

	user, ok := dbUserByID(pc.Ctrl, c, tu.BID)
	if !ok {
		return
	}

	if user.Email == "" {
		err := util.ErrEmailNo
		web.BadRes(c, err)
		pc.Log.Error().Msg(err.LStr())
		return
	}

	redisKey, ok := checkUserEmailCode(pc.Ctrl, c, &user, code, false, true)
	if !ok {
		return
	}

	if err := pc.DB.Model(&user).Update("last_ev_time", user.LastEVTime).Error; err != nil {
		web.BadRes(c, util.ErrDB)
		pc.Log.Err(err).Send()
		return
	}

	_, _ = pc.RDB.Delete(context.Background(), redisKey)

	tu2 := user.ToTu(tu)
	if tu.SK == "" {
		tu2.ReSk()
	}
	AToken(pc.Ctrl, c, tu2)
}

// RevMobile godoc
//
//	@Summary		校验手机验证码，成功证明是本人手机
//	@Description	成功返回一个 token 可以用于保存当前用户的所有信息
//	@Tags			pub
//	@Accept			json
//	@Product		json
//	@Param			Authorization	header	string			true	"Authentication header"
//	@Param			input			body	model.CodeInput	true	"输入你收到的验证码"
//	@Security		ApiKeyAuth
//	@Success		200	{object}	web.HttpData[string]
//	@Failure		502	{object}	web.HttpMsg
//	@Router			/pub/revMobile [post]
func (pc *PubCtrl) RevMobile(c *gin.Context) {
	var reVerifyInput model.CodeInput
	if err := c.ShouldBindJSON(&reVerifyInput); err != nil {
		web.BadRes(c, util.ErrInvalidArgument)
		pc.Log.Err(err).Send()
		return
	}
	code := reVerifyInput.Code

	tu, ok := shouldGetTu(pc.Ctrl, c)
	if !ok {
		return
	}

	user, ok := dbUserByID(pc.Ctrl, c, tu.BID)
	if !ok {
		return
	}

	if user.Mobile == "" {
		err := util.ErrMobileNo
		web.BadRes(c, err)
		pc.Log.Error().Msg(err.LStr())
		return
	}

	redisKey, ok := checkUserMobileCode(pc.Ctrl, c, &user, code, false, true)
	if !ok {
		return
	}

	if err := pc.DB.Model(&user).Update("last_mv_time", user.LastMVTime).Error; err != nil {
		web.BadRes(c, util.ErrDB)
		pc.Log.Err(err).Send()
		return
	}
	_, _ = pc.RDB.Delete(context.Background(), redisKey)

	tu2 := user.ToTu(tu)
	if tu.SK == "" {
		tu2.ReSk()
	}

	AToken(pc.Ctrl, c, tu2)
}

// VerifyGa godoc
//
//	@Summary		单独验证谷歌验证
//	@Description	成功返回一个 token 保存当前用户的所有信息
//	@Tags			pub
//	@Accept			json
//	@Product		json
//	@Param			code			body	model.CodeInput	true	"ga code"
//	@Param			Authorization	header	string			true	"Authentication header"
//	@Security		ApiKeyAuth
//	@Success		200	{object}	web.HttpData[string]
//	@Failure		502	{object}	web.HttpMsg
//	@Router			/pub/verifyGa [post]
func (pc *PubCtrl) VerifyGa(c *gin.Context) {
	var gaInput model.CodeInput
	if err := c.ShouldBindJSON(&gaInput); err != nil {
		web.BadRes(c, util.ErrInvalidArgument)
		pc.Log.Err(err).Send()
		return
	}

	tu, ok := shouldGetTu(pc.Ctrl, c)
	if !ok {
		return
	}
	user, ok := dbUserByID(pc.Ctrl, c, tu.BID)
	if !ok {
		return
	}

	v, er := util.ValidateTOTP(gaInput.Code, user.Ga)
	if er != nil || !v {
		web.BadRes(c, util.ErrGaInvalid)
		pc.Log.Error().Msg(er.LStr())
		return
	}
	user.LastGVTime = time.Now().Unix()
	if err := pc.DB.Model(&user).Update("last_gv_time", user.LastGVTime).Error; err != nil {
		web.BadRes(c, util.ErrDB)
		pc.Log.Err(err).Send()
		return
	}
	tu2 := user.ToTu(tu)
	if tu.SK == "" {
		tu2.ReSk()
	}

	AToken(pc.Ctrl, c, tu2)
}

// SignUp godoc
//
//	@Summary		提交用户信息表单
//	@Description	此时用户已经验证完Email或者Mobile中的至少一项，已经具备token才能提交信息
//	@Tags			user
//	@Accept			json
//	@Product		json
//	@Param			Authorization	header	string				true	"Authentication header"
//	@Param			user			body	model.SignUpInput	true	"user information"
//	@Security		ApiKeyAuth
//	@Success		200	{object}	web.HttpData[string]
//	@Failure		502	{object}	web.HttpMsg
//	@Router			/pub/signUp [post]
func (pc *PubCtrl) SignUp(c *gin.Context) {
	tu, ok := shouldGetTu(pc.Ctrl, c)
	if !ok {
		return
	}
	var signUpInput model.SignUpInput
	if err := c.ShouldBindJSON(&signUpInput); err != nil {
		web.BadRes(c, util.ErrInvalidArgument)
		pc.Log.Err(err).Send()
		return
	}
	if data, err := base64.StdEncoding.DecodeString(signUpInput.Password); err != nil {
		web.BadRes(c, util.ErrMsgDecode)
		pc.Log.Err(err).Send()
		return
	} else {
		signUpInput.Password = string(data)
	}

	if err := signUpInput.Validate(); err != nil {
		web.BadRes(c, err)
		pc.Log.Error().Msg(err.LStr())
		return
	}
	if tu.Email != "" && signUpInput.Email != tu.Email {
		err := util.ErrEmailNotEq(tu.Email, signUpInput.Email)
		web.BadRes(c, err)
		pc.Log.Error().Msg(err.LStr())
		return
	}
	if tu.Mobile != "" && signUpInput.Mobile != tu.Mobile {
		err := util.ErrMobileNotEq(tu.Mobile, signUpInput.Mobile)
		web.BadRes(c, err)
		pc.Log.Error().Msg(err.LStr())
		return
	}

	var user model.User
	email := ""
	if tu.Email != "" {
		email = tu.Email
	} else if signUpInput.Email != "" {
		email = signUpInput.Email
	}
	if email != "" {
		if !ensureEmailNotReg(pc.Ctrl, c, email) {
			return
		}
	}
	mobile := ""
	if tu.Mobile != "" {
		mobile = tu.Mobile
	} else if signUpInput.Mobile != "" {
		mi := model.MobileInput{Mobile: signUpInput.Mobile}
		mobile = mi.Number()
	}
	if mobile != "" {
		if !ensureMobileNotReg(pc.Ctrl, c, mobile) {
			return
		}
	}

	bid := pc.mustBid()
	hashedPass, err := util.HashPassword(signUpInput.Password)
	if err != nil {
		web.BadRes(c, err)
		pc.Log.Error().Msg(err.LStr())
		return
	}

	nick := pc.mustNick()
	user = model.User{
		BID:          bid,
		Nick:         nick,
		Email:        signUpInput.Email,
		Mobile:       mobile,
		Password:     hashedPass,
		FirmName:     signUpInput.FirmName,
		FirmType:     signUpInput.FirmType,
		Country:      signUpInput.Country,
		FirmVerified: 0,
		Status:       0,
		LastMVTime:   tu.MAuthTime,
		LastEVTime:   tu.EAuthTime,
		LastGVTime:   tu.GAuthTime,
	}

	uceMap, err := pc.getAddrFromKeySvr(bid, signUpInput.AppID)
	if err != nil {
		web.BadRes(c, err)
		pc.Log.Error().Msg(err.LStr())
		return
	}
	tx := pc.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if er := tx.Model(&user).Create(&user).Error; er != nil {
		tx.Rollback()
		web.BadRes(c, util.ErrDB)
		pc.Log.Err(er).Interface("user", user).Send()
		return
	}
	aa := model.AppIDAddress{
		BID:   user.BID,
		AppID: signUpInput.AppID,
		Eth:   uceMap["eth"],
	}
	if er := tx.Model(&aa).Create(&aa).Error; er != nil {
		tx.Rollback()
		web.BadRes(c, util.ErrDB)
		pc.Log.Err(er).Interface("user", user).Send()
		return
	}
	if er := tx.Commit().Error; er != nil {
		web.BadRes(c, util.ErrDB)
		pc.Log.Err(er).Interface("user", user).Send()
		return
	}
	tu1 := user.ToTU(&aa)
	tu2 := user.ToTu(tu1)
	tu2.ReSk()
	AToken(pc.Ctrl, c, tu2)

	uceMap["uid"] = bid
	uceMap["app_id"] = signUpInput.AppID
	k, _ := json.Marshal(uceMap)
	msg.Send(pc.Prd, "registrar-user-created", user.BID, k)
}

func (pc *PubCtrl) mustBid() string {
	bid, _ := util.CryptoRandomNumerical(model.BidLen)
	var user model.User
	err := pc.DB.Model(&user).First(&user, "b_id=?", bid).Error
	if err == nil {
		return pc.mustBid()
	} else {
		if strings.Contains(err.Error(), "not found") {
			return bid
		} else {
			return pc.mustBid()
		}
	}
}

func (pc *PubCtrl) mustNick() string {
	nick, _ := util.CryptoRandomString(6)
	nick = "Anonymous-User-" + nick
	var user model.User
	err := pc.DB.Model(&user).First(&user, "nick=?", nick).Error
	if err == nil {
		return pc.mustNick()
	} else {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nick
		} else {
			return pc.mustNick()
		}
	}
}

// getAddrFromKeySvr get addresses from key svr concurrently
// now we support eth &
// will support btc/trx asap
func (pc *PubCtrl) getAddrFromKeySvr(uid, appID string) (map[string]string, util.Err) {
	errChan := make(chan util.Err)
	keyMap := make(map[string]string, len(model.AllChainType))
	get := func(coin string) {
		addr, err := pc.doGetAddrFromKeySvr(uid+appID, coin)
		if err != nil {
			errChan <- err
		} else {
			keyMap[coin] = addr
			errChan <- nil
		}
	}
	for _, coin := range model.AllChainType {
		go get(coin)
	}
	for i := 0; i < len(model.AllChainType); i++ {
		select {
		case err := <-errChan:
			if err != nil {
				return nil, err
			}
		}
	}
	return keyMap, nil
}

type BindData struct {
	UID   string `json:"uid"`
	Chain string `json:"chain"`
}

// doGetAddrFromKeySvr uid=BID+appID
func (pc *PubCtrl) doGetAddrFromKeySvr(uid, typ string) (string, util.Err) {
	bd := BindData{
		UID:   uid,
		Chain: typ,
	}
	bdMsg, _ := json.Marshal(bd)
	enc, err := ecies.Encrypt(pc.KeyCli.Pub, bdMsg)
	if err != nil {
		pc.Log.Error().Msg(err.LStr())
		return "", err
	}
	formData := map[string]string{"data": hex.EncodeToString(enc)}

	var result string
	resp, er := pc.KeyCli.Cli.R().SetFormData(formData).SetResult(&result).Post("/bind")
	if er != nil || resp.StatusCode() != http.StatusOK {
		pc.Log.Err(er).Send()
		return "", util.ErrWalletSvr
	}
	return result, nil
}
