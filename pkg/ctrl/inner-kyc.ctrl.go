package ctrl

import (
	"context"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
	"user/pkg/model"
	"user/pkg/util"
	"user/pkg/web"
)

// QueryUserByAddr godoc
//
//	@Summary		内部查询，归集模块查询企业用户信息
//	@Description	成功返回BID+AppID
//	@Tags			pub
//	@Accept			json
//	@Param			input	body	model.AddrInput	true	"用户链上地址"
//	@Product		json
//	@Success		200	{object}	web.HttpData[model.IDResp]
//	@Failure		502	{object}	web.HttpMsg
//	@Router			/pub/i-q-user-by-addr [post]
func (pc *PubCtrl) QueryUserByAddr(c *gin.Context) {
	var ai model.AddrInput
	if err := c.ShouldBindJSON(&ai); err != nil {
		web.BadRes(c, util.ErrInvalidArgument)
		pc.Log.Err(err).Send()
		return
	}
	if err := ai.Validate(pc.Pk, pc.Log); err != nil {
		web.BadRes(c, err)
		pc.Log.Error().Msg(err.LStr())
		return
	}

	if ids, ok := iDbUserByAddr(pc.Ctrl, c, ai.Addr, ai.Type); !ok {
		return
	} else {
		web.GoodResp(c, ids)
	}
}

// KYC godoc
//
//	@Summary		企业进行后台认证，验证成功后返回认证状态给用户模块
//	@Description	成功返回Ok
//	@Tags			pub
//	@Accept			json
//	@Param			input	body	model.FirmConfirmed	true	"用户ID和认证成功的企业信息+校验字段"
//	@Product		json
//	@Success		200	{object}	web.HttpData[string]
//	@Failure		502	{object}	web.HttpMsg
//	@Router			/pub/kyc-ok [post]
func (pc *PubCtrl) KYC(c *gin.Context) {
	var fc model.FirmConfirmed
	if err := c.ShouldBindJSON(&fc); err != nil {
		web.BadRes(c, util.ErrInvalidArgument)
		pc.Log.Err(err).Send()
		return
	}
	if err := fc.Validate(pc.Pk, pc.Log); err != nil {
		web.BadRes(c, err)
		pc.Log.Error().Msg(err.LStr())
		return
	}
	if ok := dbFirmNotReg(pc.Ctrl, c, fc.FirmName, fc.Country); !ok {
		return
	}
	user, ok := dbUserByID(pc.Ctrl, c, fc.BID)
	if !ok {
		return
	}
	if user.FirmVerified != 0 && (user.FirmName != fc.FirmName || user.Country != fc.Country) {
		err := util.ErrKycUser(user.FirmName)
		web.BadRes(c, err)
		pc.Log.Error().Interface("user", user).Msg(err.LStr())
		return
	}
	user.FirmName = fc.FirmName
	user.FirmType = fc.FirmType
	user.Country = fc.Country
	user.FirmVerified = time.Now().Unix()
	user.Fid = user.BID
	updateMap := mapUpdateFirmUser(&user)
	updateMap["admin"] = true
	if err := pc.DB.Model(&user).Where("b_id=?", fc.BID).Updates(updateMap).Error; err != nil {
		web.BadRes(c, util.ErrDB)
		pc.Log.Err(err).Send()
		return
	}
	web.GoodResp(c, "Ok")
}

// KycQuery godoc
//
//	@Summary		后台管理员查询企业用户信息
//	@Description	成功返回用户信息
//	@Tags			pub
//	@Accept			json
//	@Param			input	body	model.FirmQuery	true	"用户ID/邮箱/手机+校验字段"
//	@Product		json
//	@Success		200	{object}	web.HttpData[[]model.FirmResp]
//	@Failure		502	{object}	web.HttpMsg
//	@Router			/pub/kyc-query [post]
func (pc *PubCtrl) KycQuery(c *gin.Context) {
	var fq model.FirmQuery
	if err := c.ShouldBindJSON(&fq); err != nil {
		web.BadRes(c, util.ErrInvalidArgument)
		pc.Log.Err(err).Send()
		return
	}
	if err := fq.Validate(pc.Pk, pc.Log); err != nil {
		web.BadRes(c, err)
		pc.Log.Error().Msg(err.LStr())
		return
	}

	ok := false
	var rows []model.FirmResp
	if fq.Uid != "" {
		if rows, ok = iDbUserByID(pc.Ctrl, c, fq.Uid); !ok {
			return
		}
	}
	if fq.Email != "" {
		if rows, ok = iDbUserByEmail(pc.Ctrl, c, fq.Email); !ok {
			return
		}
	}
	if fq.Mobile != "" {
		m := model.MobileInput{Mobile: fq.Mobile}
		if err := m.Validate(); err != nil {
			web.BadRes(c, err)
			return
		}
		if rows, ok = iDbUserByMobile(pc.Ctrl, c, fq.Number()); !ok {
			return
		}
	}
	web.GoodResp(c, rows)
}

// KycSetUserStatus godoc
//
//	@Summary		后台管理员设置企业用户状态
//	@Description	成功返回Ok
//	@Tags			pub
//	@Accept			json
//	@Param			input	body	model.KycUserStatusInput	true	"用户ID和设置的状态+校验字段"
//	@Product		json
//	@Success		200	{object}	web.HttpData[string]
//	@Failure		502	{object}	web.HttpMsg
//	@Router			/pub/kyc-set-user-status [post]
func (pc *PubCtrl) KycSetUserStatus(c *gin.Context) {
	var usi model.KycUserStatusInput
	if err := c.ShouldBindJSON(&usi); err != nil {
		web.BadRes(c, util.ErrInvalidArgument)
		pc.Log.Err(err).Send()
		return
	}
	if err := usi.Validate(pc.Pk, pc.Log); err != nil {
		web.BadRes(c, err)
		pc.Log.Error().Msg(err.LStr())
		return
	}
	var user model.User
	if err := pc.DB.Model(&user).Where("b_id=?", usi.Uid).Update("status", usi.Status).Error; err != nil {
		web.BadRes(c, util.ErrDB)
		pc.Log.Err(err).Send()
		return
	}
	if usi.Status < 0 {
		redisKey := "singleton:" + usi.Uid
		if ok, err := pc.RDB.Delete(context.Background(), redisKey); err != nil {
			web.BadRes(c, util.ErrDB)
			pc.Log.Error().Str("when deleting singleton:", usi.Uid).Msg(err.LStr())
			return
		} else {
			pc.Log.Info().Str("key", redisKey).Bool("del", ok).Msg("redis.del")
		}
	}
	web.GoodResp(c, "Ok")
}

// KycUserList godoc
//
//	@Summary		后台管理员查询企业用户列表
//	@Description	成功返回用户列表
//	@Tags			pub
//	@Accept			json
//	@Param			input	body	model.KycUserListInput	true	"列表参数+校验字段"
//	@Product		json
//	@Success		200	{object}	web.HttpData[[]model.FirmResp]
//	@Failure		502	{object}	web.HttpMsg
//	@Router			/pub/kyc-user-list [post]
func (pc *PubCtrl) KycUserList(c *gin.Context) {
	var uli model.KycUserListInput
	if err := c.ShouldBindJSON(&uli); err != nil {
		web.BadRes(c, util.ErrInvalidArgument)
		pc.Log.Err(err).Send()
		return
	}
	if err := uli.Validate(pc.Pk, pc.Log); err != nil {
		web.BadRes(c, err)
		pc.Log.Error().Msg(err.LStr())
		return
	}
	tx := appUserSelect(pc.Ctrl).
		Where("u.created_at between ? and ?", time.Unix(uli.Start, 0), time.Unix(uli.End, 0))
	if uli.Status != -99 {
		tx = tx.Where("u.status=?", uli.Status)
	}
	limit := uli.GetLimit()
	var rowCount int64
	if err := tx.
		Count(&rowCount).Error; err != nil {
		web.BadRes(c, util.ErrDB)
		pc.Log.Err(err).Send()
		return
	}

	var rows []model.FirmResp
	if err := tx.
		Offset(uli.GetFrom()).Limit(limit).Order("u.id asc").
		Find(&rows).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			web.BadRes(c, util.ErrNoRec)
		} else {
			web.BadRes(c, util.ErrDB)
		}
		pc.Log.Err(err).Send()
		return
	}

	web.GoodResp(c, model.UserPage{
		Total: rowCount,
		Page:  uli.Page,
		Limit: limit,
		Rows:  rows,
	})
}

// InnerUserQueryMulti godoc
//
//	@Summary		内部接口通过AppID/[]UID查询用户信息
//	@Description	成功返回用户列表
//	@Tags			pub
//	@Accept			json
//	@Param			input	body	model.IUserQueryMultiInput	true	"列表参数+校验字段"
//	@Product		json
//	@Success		200	{object}	web.HttpData[[]model.FirmResp]
//	@Failure		502	{object}	web.HttpMsg
//	@Router			/pub/i-q-user-multi [post]
func (pc *PubCtrl) InnerUserQueryMulti(c *gin.Context) {
	var uqm model.IUserQueryMultiInput
	if err := c.ShouldBindJSON(&uqm); err != nil {
		web.BadRes(c, util.ErrInvalidArgument)
		pc.Log.Err(err).Send()
		return
	}
	if err := uqm.Validate(pc.Pk, pc.Log); err != nil {
		web.BadRes(c, err)
		pc.Log.Error().Msg(err.LStr())
		return
	}
	tx := appUserSelect(pc.Ctrl).
		Where("u.b_id in ? and a.app_id=?", uqm.UidList, uqm.AppID)

	var rows []model.FirmResp
	if err := tx.Find(&rows).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			web.BadRes(c, util.ErrNoRec)
		} else {
			web.BadRes(c, util.ErrDB)
		}
		pc.Log.Err(err).Send()
		return
	}

	web.GoodResp(c, rows)
}
