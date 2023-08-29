package ctrl

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/pquerna/otp/totp"
	"github.com/rs/zerolog"
	"github.com/suiguo/hwlib/smtp"
	"gorm.io/gorm"
	"image/png"
	"strings"
	"time"
	"user/pkg/conf"
	"user/pkg/model"
	"user/pkg/redis"
	"user/pkg/util"
	"user/pkg/web"
	"user/pkg/web/auth"
)

const (
	DefaultRecent = 300
	DefaultExpire = DefaultRecent * time.Second
)

var (
	skExpire = time.Minute * time.Duration(conf.Conf.RefreshTokenAge)
	json     = jsoniter.ConfigCompatibleWithStandardLibrary
)

type Ctrl struct {
	DB  *gorm.DB
	RDB *redis.Store
	Log *zerolog.Logger
}

func rCookie(r *Ctrl, c *gin.Context, tu *model.TokenUser) bool {
	rToken, err := auth.GenRefToken(tu)
	if err != nil {
		web.BadRes(c, err)
		r.Log.Error().Msg(err.LStr())
		return false
	}
	c.SetCookie("refresh_token", rToken, conf.Conf.RefreshTokenAge*60, "/", "", true, true)
	return true
}

// AToken gen & send token to client, return true if ok
func aToken(r *Ctrl, c *gin.Context, tu *model.TokenUser) bool {
	token, err := auth.GenAToken(tu)
	if err != nil {
		web.BadRes(c, err)
		r.Log.Error().Msg(err.LStr())
		return false
	}
	web.GoodResp(c, token)
	return true
}

// AToken gen & send token to client, return true if ok
func AToken(r *Ctrl, c *gin.Context, tu *model.TokenUser) bool {
	c0 := context.Background()
	redisKey := "singleton:" + tu.BID
	if err := r.RDB.Set(c0, redisKey, tu.SK, skExpire); err != nil {
		web.BadRes(c, err)
		r.Log.Error().Msg(err.LStr())
		return false
	}

	return aToken(r, c, tu)
}

func recentToken(r *Ctrl, c *gin.Context, tu *model.TokenUser) bool {
	now := time.Now().Unix()
	if (tu.EAuthTime > 0 && now-tu.EAuthTime > DefaultRecent) ||
		(tu.MAuthTime > 0 && now-tu.MAuthTime > DefaultRecent) ||
		(tu.GAuthTime > 0 && now-tu.GAuthTime > DefaultRecent) {
		err := util.Err2FaExpire
		web.BadRes(c, err)
		r.Log.Error().Interface("tu", tu).Msg(err.LStr())
		return false
	}
	return true
}

func sendMobileNStore(r *Ctrl, c *gin.Context, mobile string) bool {
	// 如果 Redis 中有，就用旧的
	c0 := context.Background()
	redisKey := "mobile:" + mobile
	codeSix, er := r.RDB.Get(c0, redisKey)
	if er != nil {
		web.BadRes(c, er)
		r.Log.Error().Msg(er.LStr())
		return false
	}
	if codeSix == "" {
		codeSix, er = util.CryptoRandomNumerical2(6)
		if er != nil {
			web.BadRes(c, er)
			r.Log.Error().Msg(er.LStr())
			return false
		}
	}
	// FIXME send sms really
	// and put to redis, 5 minutes live time.
	er = r.RDB.Set(c0, redisKey, codeSix, DefaultExpire)
	if er != nil {
		web.BadRes(c, er)
		r.Log.Error().Msg(er.LStr())
		return false
	}
	// FIXME 现在直接返回给前端为了调试
	web.GoodResp(c, codeSix)
	return true
}

func sendEmailNStore(r *Ctrl, c *gin.Context, email string) bool {
	// redis 中如果已经有了就用已有的
	c0 := context.Background()
	redisKey := "email:" + email
	codeSix, er := r.RDB.Get(c0, redisKey)
	if er != nil {
		web.BadRes(c, er)
		r.Log.Error().Msg(er.LStr())
		return false
	}
	if codeSix == "" {
		codeSix, er = util.CryptoRandomNumerical(6)
		if er != nil {
			web.BadRes(c, er)
			r.Log.Error().Msg(er.LStr())
			return false
		}
	}

	cli := smtp.GetClient("smtp.gmail.com", 465, "suiguo3564@gmail.com", "hzdradxhxswufoxs")
	err := cli.SendMail(
		smtp.WithFrom("suiguo3564@gmail.com"),
		smtp.WithTo(email),
		smtp.WithTitle("验证码"),
		smtp.WithBodyReg(codeSix),
	)
	if err != nil {
		r.Log.Err(err).Send()
		er := util.ErrEmailSend
		web.BadRes(c, er)
		return false
	}

	// and put to redis, 5 minutes live time.
	er = r.RDB.Set(c0, redisKey, codeSix, DefaultExpire)
	if er != nil {
		web.BadRes(c, er)
		r.Log.Error().Msg(er.LStr())
		return false
	}
	web.GoodResp(c, "Ok")
	return true
}

func emailCodeExists(r *Ctrl, c *gin.Context, email, code string) bool {
	// 查看 Redis 中有没有
	c0 := context.Background()
	redisKey := "email:" + email
	codeSix, err := r.RDB.Get(c0, redisKey)
	if err != nil {
		web.BadRes(c, err)
		r.Log.Error().Msg(err.LStr())
		return false
	}
	if code != codeSix {
		err = util.ErrEmailCode("")
		web.BadRes(c, err)
		r.Log.Error().Msg(err.LStr())
		return false
	}
	_, _ = r.RDB.Delete(c0, redisKey)
	return true
}

func mobileCodeExists(r *Ctrl, c *gin.Context, mobile, code string) bool {
	// 查看 Redis 中有没有
	c0 := context.Background()
	redisKey := "mobile:" + mobile
	codeSix, err := r.RDB.Get(c0, redisKey)
	if err != nil {
		web.BadRes(c, err)
		r.Log.Error().Msg(err.LStr())
		return false
	}
	if code != codeSix {
		err = util.ErrMobileCode("")
		web.BadRes(c, err)
		r.Log.Error().Msg(err.LStr())
		return false
	}
	_, _ = r.RDB.Delete(c0, redisKey)
	return true
}

func checkUserMobileCode(r *Ctrl, c *gin.Context, u *model.User, code string, must, show bool) (redisMKey string, ok bool) {
	redisMKey = ""
	ok = false
	c0 := context.Background()

	if must && code == "" {
		err := util.ErrMobileFirst
		web.BadRes(c, err)
		r.Log.Error().Msg(err.LStr())
		return
	}

	// 查看 Redis 中有没有 MCode
	redisMKey = "mobile:" + u.Mobile
	codeSix, err := r.RDB.Get(c0, redisMKey)
	if err != nil {
		web.BadRes(c, err)
		r.Log.Error().Msg(err.LStr())
		return
	}
	if code != codeSix {
		if show {
			err = util.ErrMobileCode(u.Mobile)
		} else {
			err = util.ErrMobileCode("")
		}
		web.BadRes(c, err)
		r.Log.Error().Msg(err.LStr())
		return
	}
	u.LastMVTime = time.Now().Unix()
	ok = true
	return
}

func checkUserEmailCode(r *Ctrl, c *gin.Context, u *model.User, code string, must, show bool) (redisEKey string, ok bool) {
	c0 := context.Background()
	redisEKey = ""
	ok = false
	if must && code == "" {
		err := util.ErrEmailFirst
		web.BadRes(c, err)
		r.Log.Error().Msg(err.LStr())
		return
	}

	// 查看 Redis 中有没有 ECode
	redisEKey = "email:" + u.Email
	codeSix, err := r.RDB.Get(c0, redisEKey)
	if err != nil {
		web.BadRes(c, err)
		r.Log.Error().Msg(err.LStr())
		return
	}
	if code != codeSix {
		if show {
			err = util.ErrEmailCode(u.Email)
		} else {
			err = util.ErrEmailCode("")
		}
		web.BadRes(c, err)
		r.Log.Error().Msg(err.LStr())
		return
	}
	u.LastEVTime = time.Now().Unix()
	ok = true
	return
}

func checkUserGCode(r *Ctrl, c *gin.Context, u *model.User, code string) bool {
	if code == "" {
		err := util.ErrGaFirst
		web.BadRes(c, err)
		r.Log.Error().Msg(err.LStr())
		return false
	}
	v, err := util.ValidateTOTP(code, u.Ga)
	if err != nil || !v {
		web.BadRes(c, util.ErrGaInvalid)
		r.Log.Error().Msg(err.LStr())
		return false
	}
	u.LastGVTime = time.Now().Unix()
	return true
}

func mapUpdateFirmUser(user *model.User) map[string]interface{} {
	return map[string]interface{}{
		"firm_name":     user.FirmName,
		"firm_type":     user.FirmType,
		"country":       user.Country,
		"firm_verified": user.FirmVerified,
		"fid":           user.Fid,
	}
}

func dbUser(r *Ctrl, c *gin.Context, name string) (user model.User, ok bool) {
	ok = false
	if strings.Contains(name, "@") {
		err := r.DB.Model(&user).First(&user, "last_ev_time>0 and email=?", name).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				web.BadRes(c, util.ErrNoRec)
			} else {
				web.BadRes(c, util.ErrDB)
			}
			r.Log.Err(err).Send()
			return
		}
	} else {
		err := r.DB.Model(&user).First(&user, "last_mv_time>0 and mobile=?", name).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				web.BadRes(c, util.ErrNoRec)
			} else {
				web.BadRes(c, util.ErrDB)
			}
			r.Log.Err(err).Send()
			return
		}
	}
	ok = true
	return
}

func dbUserByEmail(r *Ctrl, c *gin.Context, email string) (user model.User, ok bool) {
	ok = false
	err := r.DB.Model(&user).First(&user, "last_ev_time>0 and email=?", email).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			web.BadRes(c, util.ErrNoRec)
		} else {
			web.BadRes(c, util.ErrDB)
		}
		r.Log.Err(err).Send()
		return
	}
	ok = true
	return
}

func dbUserByMobile(r *Ctrl, c *gin.Context, mobile string) (user model.User, ok bool) {
	ok = false
	err := r.DB.Model(&user).First(&user, "last_mv_time>0 and mobile=?", mobile).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			web.BadRes(c, util.ErrNoRec)
		} else {
			web.BadRes(c, util.ErrDB)
		}
		r.Log.Err(err).Send()
		return
	}
	ok = true
	return
}

func iDbUserByAddr(r *Ctrl, c *gin.Context, addr, typ string) (user model.IDResp, ok bool) {
	ok = false
	var aia model.AppIDAddress
	cond := fmt.Sprintf("%s=?", typ)
	err := r.DB.Model(&aia).First(&user, cond, addr).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			web.BadRes(c, util.ErrNoRec)
		} else {
			web.BadRes(c, util.ErrDB)
		}
		r.Log.Err(err).Send()
		return
	}
	ok = true
	return
}

func dbUserByID(r *Ctrl, c *gin.Context, bid string) (user model.User, ok bool) {
	ok = false
	err := r.DB.Model(&user).First(&user, "b_id=?", bid).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			web.BadRes(c, util.ErrNoRec)
		} else {
			web.BadRes(c, util.ErrDB)
		}
		r.Log.Err(err).Send()
		return
	}
	ok = true
	return
}

func dbFirmNotReg(r *Ctrl, c *gin.Context, firmName, country string) bool {
	var user model.User
	err := r.DB.Model(&user).First(&user, "firm_name=? and country=? and firm_verified>0", firmName, country).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return true
		} else {
			web.BadRes(c, util.ErrDB)
			r.Log.Err(err).Send()
			return false
		}
	}
	er := util.ErrKycFirm(user.Fid)
	web.BadRes(c, er)
	r.Log.Error().Interface("user", user).Msg(er.LStr())
	return false
}

func dbFirmAddUser(r *Ctrl, c *gin.Context, bid, boss string) bool {
	var mb model.User
	if err := r.DB.Model(&mb).First(&mb, "b_id=?", boss).Error; err != nil {
		web.BadRes(c, util.ErrDB)
		r.Log.Err(err).Send()
		return false
	}
	var user model.User
	if err := r.DB.Model(&user).
		Where("b_id=? and fid=''", bid).
		Updates(mapUpdateFirmUser(&mb)).Error; err != nil {
		web.BadRes(c, util.ErrDB)
		r.Log.Err(err).Send()
		return false
	}
	return true
}

func dbFirmDelUser(r *Ctrl, c *gin.Context, bid, boss string) bool {
	var user model.User
	if err := r.DB.Model(&user).
		Where("b_id=? and fid=?", bid, boss).
		Updates(map[string]interface{}{
			"fid":           "",
			"firm_verified": 0,
		}).Error; err != nil {
		web.BadRes(c, util.ErrDB)
		r.Log.Err(err).Send()
		return false
	}
	return true
}

func dbFirmListUser(r *Ctrl, c *gin.Context, boss string) (users []model.FirmUser, ok bool) {
	var user model.User
	err := r.DB.Model(&user).Find(&users, "fid=?", boss).Error
	if err != nil {
		web.BadRes(c, util.ErrDB)
		r.Log.Err(err).Send()
		return
	}
	ok = true
	return
}

func appUserSelect(r *Ctrl) *gorm.DB {
	return r.DB.
		Table("users as u, app_id_addresses as a").
		Select("u.b_id as uid, u.nick, u.email, u.mobile, u.firm_name, u.country, " +
			"u.firm_verified, u.admin, u.status, floor(extract(epoch from u.created_at)) as created," +
			"u.fid, a.app_id, a.eth").Where("u.b_id=a.b_id")
}

func dbAppAddr(r *Ctrl, c *gin.Context, bid, appID string) (aa model.AppIDAddress, ok bool) {
	err := r.DB.Model(&aa).First(&aa, "b_id=? and app_id=?", bid, appID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			web.BadRes(c, util.ErrNoRec)
		} else {
			web.BadRes(c, util.ErrDB)
		}
		r.Log.Err(err).Send()
		return
	}
	ok = true
	return
}

func iDbUserByID(r *Ctrl, c *gin.Context, bid string) (users []model.FirmResp, ok bool) {
	ok = false
	err := appUserSelect(r).Where("u.b_id=?", bid).Find(&users).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			web.BadRes(c, util.ErrNoRec)
		} else {
			web.BadRes(c, util.ErrDB)
		}
		r.Log.Err(err).Send()
		return
	}
	ok = true
	return
}

func iDbUserByEmail(r *Ctrl, c *gin.Context, email string) (users []model.FirmResp, ok bool) {
	ok = false
	err := appUserSelect(r).Where("u.email=?", email).Find(&users).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			web.BadRes(c, util.ErrNoRec)
		} else {
			web.BadRes(c, util.ErrDB)
		}
		r.Log.Err(err).Send()
		return
	}
	ok = true
	return
}

func iDbUserByMobile(r *Ctrl, c *gin.Context, mobile string) (users []model.FirmResp, ok bool) {
	ok = false
	err := appUserSelect(r).Where("u.mobile=?", mobile).Find(&users).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			web.BadRes(c, util.ErrNoRec)
		} else {
			web.BadRes(c, util.ErrDB)
		}
		r.Log.Err(err).Send()
		return
	}
	ok = true
	return
}
func ensureUserHasEmail(r *Ctrl, c *gin.Context, user *model.User) bool {
	if user.LastEVTime == 0 || user.Email == "" {
		err := util.ErrEmailNo
		web.BadRes(c, err)
		r.Log.Error().Msg(err.LStr())
		return false
	}
	return true
}

func ensureUserHasMobile(r *Ctrl, c *gin.Context, user *model.User) bool {
	if user.LastMVTime == 0 || user.Mobile == "" {
		er := util.ErrMobileNo
		web.BadRes(c, er)
		r.Log.Error().Msg(er.LStr())
		return false
	}
	return true
}

func ensureEmailExists(r *Ctrl, c *gin.Context, email string) bool {
	// 查询用户是否存在
	var user model.User
	err := r.DB.Model(&user).First(&user, "last_ev_time>0 and email=?", email).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			web.BadRes(c, util.ErrNoRec)
		} else {
			web.BadRes(c, util.ErrDB)
		}
		r.Log.Err(err).Send()
		return false
	}
	return true
}

func ensureMobileExists(r *Ctrl, c *gin.Context, mobile string) bool {
	// 查询用户是否存在
	var user model.User
	err := r.DB.Model(&user).First(&user, "last_mv_time>0 and mobile=?", mobile).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			web.BadRes(c, util.ErrNoRec)
		} else {
			web.BadRes(c, util.ErrDB)
		}
		r.Log.Err(err).Send()
		return false
	}
	return true
}

func ensureEmailNotReg(r *Ctrl, c *gin.Context, email string) bool {
	var user model.User
	err := r.DB.Model(&user).First(&user, "last_ev_time>0 and email=?", email).Error
	if err == nil {
		er := util.ErrEmailExists(email)
		web.BadRes(c, er)
		r.Log.Error().Msg(er.LStr())
		return false
	}
	return true
}

func ensureNickNotReg(r *Ctrl, c *gin.Context, nick string) bool {
	var user model.User
	err := r.DB.Model(&user).First(&user, "nick=?", nick).Error
	if err == nil {
		er := util.ErrNickExists(nick)
		web.BadRes(c, er)
		r.Log.Error().Msg(er.LStr())
		return false
	}
	return true
}

func ensureMobileNotReg(r *Ctrl, c *gin.Context, mobile string) bool {
	var user model.User
	err := r.DB.Model(&user).First(&user, "last_mv_time>0 and mobile=?", mobile).Error
	if err == nil {
		er := util.ErrMobileExists(mobile)
		web.BadRes(c, er)
		r.Log.Warn().Msg(er.LStr())
		return false
	}
	return true
}

func tuDbOk(r *Ctrl, c *gin.Context, tu *model.TokenUser) bool {
	var user model.User
	if tu.EAuthTime > 0 {
		err := r.DB.Model(&user).First(&user, "last_ev_time>0 and email=?", tu.Email).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				web.BadRes(c, util.ErrNoRec)
			} else {
				web.BadRes(c, util.ErrDB)
			}
			r.Log.Err(err).Send()
			return false
		}
	} else if tu.MAuthTime > 0 {
		err := r.DB.Model(&user).First(&user, "last_mv_time>0 and mobile=?", tu.Mobile).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				web.BadRes(c, util.ErrNoRec)
			} else {
				web.BadRes(c, util.ErrDB)
			}
			r.Log.Err(err).Send()
			return false
		}
	} else {
		err := util.ErrTokenInvalid
		web.BadRes(c, err)
		r.Log.Error().Msg(err.LStr())
		return false
	}
	if user.Status < 0 {
		err := util.ErrTokenInvalid
		web.BadRes(c, err)
		r.Log.Error().Msg(err.LStr())
		return false
	}
	tu.Admin = user.Admin
	tu.FirmName = user.FirmName
	tu.KycTime = user.FirmVerified
	tu.EAuthTime = user.LastEVTime
	tu.MAuthTime = user.LastMVTime
	tu.GAuthTime = user.LastGVTime
	return true
}

func getRTu(r *Ctrl, c *gin.Context) *model.TokenUser {
	tu := web.GetTokenUser(c)
	if tu == nil {
		rtStr, err := c.Cookie("refresh_token")
		if err != nil {
			rtStr = c.Request.Header.Get("Refresh-Token")
		}
		if rtStr != "" {
			tu2, er := auth.ValidateToken(rtStr, conf.Conf.RefreshTokenPublicKey)
			if er != nil {
				web.BadRes(c, er)
				r.Log.Error().Msg(er.LStr())
				return nil
			}
			return &tu2
		} else {
			return nil
		}
	} else {
		return tu
	}
}

func shouldGetRTu(r *Ctrl, c *gin.Context) (*model.TokenUser, bool) {
	tu := getRTu(r, c)
	if tu == nil {
		err := util.ErrTokenInvalid
		web.BadRes(c, err)
		r.Log.Error().Msg(err.LStr())
		return nil, false
	} else {
		return sku(r, c, tu)
	}
}
func shouldGetTu(r *Ctrl, c *gin.Context) (*model.TokenUser, bool) {
	tu := web.GetTokenUser(c)
	if tu == nil {
		err := util.ErrTokenInvalid
		web.BadRes(c, err)
		r.Log.Error().Msg(err.LStr())
		return nil, false
	}
	return tu, true
}

func mustGetTu(r *Ctrl, c *gin.Context) (*model.TokenUser, bool) {
	tu := web.GetTokenUser(c)
	if tu == nil {
		err := util.ErrTokenInvalid
		web.BadRes(c, err)
		r.Log.Error().Msg(err.LStr())
		return nil, false
	} else {
		return sku(r, c, tu)
	}
}

func sku(r *Ctrl, c *gin.Context, tu *model.TokenUser) (*model.TokenUser, bool) {
	redisKey := "singleton:" + tu.BID
	if sk, err := r.RDB.Get(context.Background(), redisKey); err != nil {
		web.BadRes(c, err)
		r.Log.Error().Msg(err.LStr())
		return nil, false
	} else {
		if sk != tu.SK {
			err := util.ErrOldLogin
			web.BadRes(c, err)
			r.Log.Error().Msg(err.LStr())
			return nil, false
		}
	}
	return tu, true
}

func genGa(r *Ctrl, c *gin.Context, tu *model.TokenUser) (*model.GenGaResponse, bool) {
	tk, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "FatChain",
		AccountName: tu.Name(),
	})
	if err != nil {
		web.BadRes(c, util.ErrGaGen)
		r.Log.Err(err).Send()
		return nil, false
	}
	// Convert TOTP key into a PNG
	var buf bytes.Buffer
	img, err := tk.Image(200, 200)
	if err != nil {
		web.BadRes(c, util.ErrGaGen)
		r.Log.Err(err).Send()
		return nil, false
	}
	err = png.Encode(&buf, img)
	if err != nil {
		web.BadRes(c, util.ErrGaGen)
		r.Log.Err(err).Send()
		return nil, false
	}
	// 存
	gaSec := tk.Secret()
	hashedGa, er := util.HashSecret(gaSec)
	if er != nil {
		web.BadRes(c, util.ErrGaGen)
		r.Log.Error().Msg(er.LStr())
		return nil, false
	}
	// 改存 redis
	redisGKey := "ga:" + tu.BID
	c0 := context.Background()
	er = r.RDB.Set(c0, redisGKey, hashedGa, DefaultExpire)
	if er != nil {
		web.BadRes(c, er)
		r.Log.Error().Interface("tu", tu).Msg(er.LStr())
		return nil, false
	}
	imgTxt := base64.StdEncoding.EncodeToString(buf.Bytes())
	return &model.GenGaResponse{Image: imgTxt, Text: gaSec}, true
}

func bindGa(r *Ctrl, c *gin.Context, user *model.User, code string) (redisGKey string, ok bool) {
	c0 := context.Background()
	redisGKey = "ga:" + user.BID
	ok = false
	hashedGa, er := r.RDB.Get(c0, redisGKey)
	if er != nil {
		web.BadRes(c, er)
		r.Log.Error().Msg(er.LStr())
		return
	}
	v, er := util.ValidateTOTP(code, hashedGa)
	if er != nil {
		web.BadRes(c, util.ErrGaInvalid)
		r.Log.Error().Msg(er.LStr())
		return
	}
	if !v {
		web.BadRes(c, util.ErrGaInvalid)
		r.Log.Error().Str("bid", user.BID).Msg("bindGa")
		return
	}
	user.Ga = hashedGa
	user.LastGVTime = time.Now().Unix()
	ok = true
	return
}
