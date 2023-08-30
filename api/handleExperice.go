package api

import (
	"fmt"
	"github.com/ethereum/BGService/db"
	"github.com/ethereum/BGService/types"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (a *ApiService) checkoutQualification(c *gin.Context) {
	uid := c.Query("uid")
	res := types.HttpRes{}
	// 查询三个条件-查询数据库
	user, err := db.GetUser(a.dbEngine, uid)

	if err != nil {
		logrus.Info("查询db发生错误", err)

		res.Code = -1
		res.Message = "查询db发生错误"
		res.Data = err
		c.SecureJSON(http.StatusOK, res)
		return
	}

	if user == nil {
		logrus.Info("未找到用户记录", uid)

		res.Code = -1
		res.Message = "未找到用户记录"
		c.SecureJSON(http.StatusOK, res)
		return
	}

	if user.IsBindGoogle == false {
		logrus.Info("条件不满足，google未绑定", user.IsBindGoogle)

		res.Code = -1
		res.Message = "条件不满足，google未绑定"
		c.SecureJSON(http.StatusOK, res)
		return
	}

	userBindInfos, err := db.GetUserBindInfos(a.dbEngine, uid)

	if err != nil {
		logrus.Info("查询db发生错误", err)

		res.Code = -1
		res.Message = "查询db发生错误"
		res.Data = err
		c.SecureJSON(http.StatusOK, res)
		return
	}

	if userBindInfos == nil {
		logrus.Info("未找到用户绑定记录")

		res.Code = -1
		res.Message = "未找到用户绑定记录"
		res.Data = err
		c.SecureJSON(http.StatusOK, res)
		return
	}

	if len(userBindInfos.ApiKey) == 0 || len(userBindInfos.ApiSecret) == 0 {
		logrus.Info("找到该用户的绑定记录，但是其中有一项apikey或者apiSecret为空", uid)

		res.Code = -1
		res.Message = "找到该用户的绑定记录，但是其中有一项apikey或者apiSecret为空"
		c.SecureJSON(http.StatusOK, res)
		return
	}
	if user.InviteNumber < 1 {
		logrus.Info("条件不满足，当前未邀请人", user.InviteNumber)

		res.Code = -1
		res.Message = "条件不满足，当前未邀请人"
		c.SecureJSON(http.StatusOK, res)
		return
	}

	res.Code = 1
	res.Message = "当前符合领取体验金资格"
	res.Data = ""
	c.SecureJSON(http.StatusOK, res)
	return
}

func (a *ApiService) getExperienceFund(c *gin.Context) {
	uid := c.Query("uid")

	res := types.HttpRes{}

	// 查询三个条件-查询数据库
	user, err := db.GetUser(a.dbEngine, uid)

	if err != nil {
		logrus.Info("查询db发生错误", err)

		res.Code = -1
		res.Message = "查询db发生错误"
		res.Data = err
		c.SecureJSON(http.StatusOK, res)
		return
	}
	if user == nil {
		logrus.Info("未找到用户记录")

		res.Code = -1
		res.Message = "未找到用户记录"
		res.Data = err
		c.SecureJSON(http.StatusOK, res)
		return
	}

	//这里首先也要检查资格
	if user.IsBindGoogle == false {
		logrus.Info("条件不满足，google未绑定", user.IsBindGoogle)
		res.Code = -1
		res.Message = "领取体验金失败：未绑定google"
		c.SecureJSON(http.StatusOK, res)
	}
	userBindInfos, err := db.GetUserBindInfos(a.dbEngine, uid)

	if err != nil {
		logrus.Info("查询db发生错误", err)

		res.Code = -1
		res.Message = "查询db发生错误"
		res.Data = err
		c.SecureJSON(http.StatusOK, res)
		return
	}

	if len(userBindInfos.ApiKey) == 0 || len(userBindInfos.ApiSecret) == 0 {
		logrus.Info("找到该用户的绑定记录，但是其中有一项apikey或者apiSecret为空", uid)

		res.Code = -1
		res.Message = "领取体验金失败：找到该用户的绑定记录，但是其中有一项apikey或者apiSecret为空"
		c.SecureJSON(http.StatusOK, res)
	}

	if user.InviteNumber < 1 {
		logrus.Info("条件不满足，当前未邀请人", user.InviteNumber)
		res.Code = -1
		res.Message = "领取体验金失败：当前未邀请人"
		c.SecureJSON(http.StatusOK, res)
	}
	logrus.Info("条件符合。可以领取体验金", uid)

	session := a.dbEngine.NewSession()
	err = session.Begin()
	if err != nil {
		return
	}
	//平台体验金资金池减少相应的数额
	TotalRevenueInfo, err := db.GetPlatformExperience(a.dbEngine)

	if err != nil {
		logrus.Info("查询db发生错误", err)

		res.Code = -1
		res.Message = "查询db发生错误"
		res.Data = err
		c.SecureJSON(http.StatusOK, res)
		return
	}

	if TotalRevenueInfo == nil {
		logrus.Info("平台体验金信息不存在，请核对")

		res.Code = -1
		res.Message = "平台体验金信息不存在，请核对"
		res.Data = err
		c.SecureJSON(http.StatusOK, res)
		return
	}
	//总金额减少
	TotalRevenueInfo.TotalSum = TotalRevenueInfo.TotalSum - TotalRevenueInfo.PerSum
	//接收人数+1
	TotalRevenueInfo.ReceivePersons = TotalRevenueInfo.ReceivePersons + 1
	//更新
	_, err = session.Table("platformExperience").Update(TotalRevenueInfo)
	if err != nil {
		err := session.Rollback()
		if err != nil {
			return
		}
		logrus.Fatal(err)
	}

	//首先用户体验金增加一条记录
	userExperience := types.UserExperience{}

	userExperience.Uid = uid
	//userExperience.BenefitRatio = 0
	//userExperience.BenefitSum = 0
	userExperience.ReceiveDays = 1
	userExperience.ReceiveSum = TotalRevenueInfo.PerSum

	_, err = session.Table("userExperience").Insert(userExperience)
	if err != nil {
		err := session.Rollback()
		if err != nil {
			return
		}
		logrus.Fatal(err)
	}

	err = session.Commit()
	if err != nil {
		logrus.Fatal(err)
	}

	res.Code = 0
	res.Message = "领取体验金成功"
	c.SecureJSON(http.StatusOK, res)
	return
}

func (a *ApiService) getUserExperience(c *gin.Context) {
	uid := c.Query("uid")
	fmt.Println(uid)

	res := types.HttpRes{}

	userExperience, err := db.GetUserExperience(a.dbEngine, uid)

	if err != nil {
		logrus.Info("查询db发生错误", err)

		res.Code = -1
		res.Message = "查询db发生错误"
		res.Data = err
		c.SecureJSON(http.StatusOK, res)
		return
	}

	if userExperience == nil {
		logrus.Info("用户体验信息记录不存在")

		res.Code = -1
		res.Message = "用户体验信息记录不存在"
		res.Data = err
		c.SecureJSON(http.StatusOK, res)
		return
	}

	res.Code = 0
	res.Message = "获取用户体验金信息成功"
	res.Data = userExperience

	c.SecureJSON(http.StatusOK, res)
	return
}
func (a *ApiService) getPlatformExperience(c *gin.Context) {
	res := types.HttpRes{}

	platformExperience, err := db.GetPlatformExperience(a.dbEngine)

	if err != nil {
		logrus.Info("查询db发生错误", err)

		res.Code = -1
		res.Message = "查询db发生错误"
		res.Data = err
		c.SecureJSON(http.StatusOK, res)
		return
	}
	if platformExperience == nil {
		logrus.Info("平台体验信息不存在，请核对")

		res.Code = -1
		res.Message = "平台体验信息不存在，请核对"
		res.Data = err
		c.SecureJSON(http.StatusOK, res)
		return
	}
	res.Code = 0
	res.Message = "获取平台体验金信息成功"
	res.Data = platformExperience

	c.SecureJSON(http.StatusOK, res)
	return
}
