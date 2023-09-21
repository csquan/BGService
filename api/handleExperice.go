package api

import (
	"fmt"
	"github.com/ethereum/BGService/db"
	"github.com/ethereum/BGService/types"
	"github.com/ethereum/BGService/util"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

// 活动资格-注册-绑定google 绑定APIKEY
func (a *ApiService) checkoutQualification(c *gin.Context) {
	uid, _ := c.Get("Uid")
	uidFormatted := fmt.Sprintf("%s", uid)
	body := make(map[string]interface{})
	body["apiBinding"] = false
	body["isBindGoogle"] = false
	// 查询两个条件-查询数据库
	// 谷歌绑定检查
	user, err := db.GetUser(a.dbEngine, uidFormatted)
	if err != nil {
		logrus.Info("query db error", err)
		res := util.ResponseMsg(-1, "query db error", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}
	if user == nil {
		logrus.Info("find no user", uid)
		res := util.ResponseMsg(-1, "find no user", nil)
		c.SecureJSON(http.StatusOK, res)
		return
	}
	if user.IsBindGoogle == true {
		body["isBindGoogle"] = true
	}
	// api绑定检查
	userBindInfos, err := db.GetUserBindInfos(a.dbEngine, uidFormatted)
	if err != nil {
		logrus.Info("query db error", err)
		res := util.ResponseMsg(-1, "query db error", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}
	if userBindInfos != nil {
		if len(userBindInfos.ApiKey) > 0 || len(userBindInfos.ApiSecret) > 0 {
			body["apiBinding"] = true
		}
	}
	res := util.ResponseMsg(-1, "checkoutQualification success", body)
	c.SecureJSON(http.StatusOK, res)
	return
}

func (a *ApiService) getExperienceFund(c *gin.Context) {
	uid, _ := c.Get("Uid")
	uidFormatted := fmt.Sprintf("%s", uid)

	// 查询三个条件-查询数据库
	user, err := db.GetUser(a.dbEngine, uidFormatted)

	if err != nil {
		logrus.Info("query db error", err)

		res := util.ResponseMsg(-1, "query db error", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}
	if user == nil {
		logrus.Info("find no user")

		res := util.ResponseMsg(-1, "find no user", nil)
		c.SecureJSON(http.StatusOK, res)
		return
	}

	//这里首先也要检查资格
	if user.IsBindGoogle == false {
		logrus.Info("google is not bind", user.IsBindGoogle)

		res := util.ResponseMsg(-1, "google is not bind", nil)
		c.SecureJSON(http.StatusOK, res)
		return
	}
	userBindInfos, err := db.GetUserBindInfos(a.dbEngine, uidFormatted)

	if err != nil {
		logrus.Info("query db error", err)

		res := util.ResponseMsg(-1, "query db error", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}

	if len(userBindInfos.ApiKey) == 0 || len(userBindInfos.ApiSecret) == 0 {
		logrus.Info("find user bind info,but apikey or apiSecret is null", uid)

		res := util.ResponseMsg(-1, "find user bind info,but apikey or apiSecret is null", nil)
		c.SecureJSON(http.StatusOK, res)
		return
	}
	if user.InviteNumber < 1 {
		logrus.Info("condition is not satisfied,no invite person", user.InviteNumber)

		res := util.ResponseMsg(-1, "condition is not satisfied,no invite person", nil)
		c.SecureJSON(http.StatusOK, res)
		return
	}
	logrus.Info("condition is satisfied,can get money", uid)

	session := a.dbEngine.NewSession()
	err = session.Begin()
	if err != nil {
		return
	}
	//平台体验金资金池减少相应的数额
	TotalRevenueInfo, err := db.GetPlatformExperience(a.dbEngine)

	if err != nil {
		logrus.Info("query db error", err)

		res := util.ResponseMsg(-1, "query db error", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}

	if TotalRevenueInfo == nil {
		logrus.Info("platform exp info not exist")

		res := util.ResponseMsg(-1, "platform exp info not exist", nil)
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
			logrus.Error(err)

			res := util.ResponseMsg(0, "internal db session rollback error", err)
			c.SecureJSON(http.StatusOK, res)
			return
		}
	}

	//首先用户体验金增加一条记录
	userExperience := types.UserExperience{}

	userExperience.UId = uidFormatted
	userExperience.ReceiveDays = 1
	userExperience.ReceiveSum = TotalRevenueInfo.PerSum

	_, err = session.Table("userExperience").Insert(userExperience)
	if err != nil {
		err := session.Rollback()
		if err != nil {
			logrus.Error(err)

			res := util.ResponseMsg(0, "internal db session rollback error", err)
			c.SecureJSON(http.StatusOK, res)
			return
		}
	}

	err = session.Commit()
	if err != nil {
		logrus.Error(err)

		res := util.ResponseMsg(0, "internal db session commit  error", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}

	res := util.ResponseMsg(0, "get exp success", nil)
	c.SecureJSON(http.StatusOK, res)
	return
}

func (a *ApiService) getUserExperienceRatio(c *gin.Context) {
	uid, _ := c.Get("Uid")
	uidFormatted := fmt.Sprintf("%s", uid)

	userExperience, err := db.GetUserExperience(a.dbEngine, uidFormatted)

	if err != nil {
		logrus.Info("query db error", err)

		res := util.ResponseMsg(-1, "query db error", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}

	if userExperience == nil {
		logrus.Info("user exp info not exist")

		res := util.ResponseMsg(-1, "user exp info not exist", nil)
		c.SecureJSON(http.StatusOK, res)
		return
	}

	res := util.ResponseMsg(0, "user get exp success", userExperience)
	c.SecureJSON(http.StatusOK, res)
	return
}

func (a *ApiService) getPlatformExperienceRatio(c *gin.Context) {
	platformExperience, err := db.GetPlatformExperience(a.dbEngine)

	if err != nil {
		logrus.Info("query db error", err)

		res := util.ResponseMsg(-1, "query db error", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}
	if platformExperience == nil {
		logrus.Info("platform exp not exist")

		res := util.ResponseMsg(-1, "platform exp not exist", nil)
		c.SecureJSON(http.StatusOK, res)
		return
	}

	res := util.ResponseMsg(0, "get platform exp success", platformExperience)
	c.SecureJSON(http.StatusOK, res)
	return
}
