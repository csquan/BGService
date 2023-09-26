package api

import (
	"errors"
	"fmt"
	"github.com/ethereum/BGService/db"
	"github.com/ethereum/BGService/types"
	"github.com/ethereum/BGService/util"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

// 活动资格--绑定google 绑定APIKEY
func (a *ApiService) checkoutQualification(c *gin.Context) {
	uid, _ := c.Get("Uid")
	uidFormatted := fmt.Sprintf("%s", uid)

	body := make(map[string]interface{})
	body["apiBinding"] = false
	body["isBindGoogle"] = false

	body["balance"] = 0
	body["RegisterStatus"] = false
	body["apiBindingStatus"] = false
	body["isBindGoogleStatus"] = false
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
		if len(userBindInfos.ApiKey) > 0 && len(userBindInfos.ApiSecret) > 0 {
			body["apiBinding"] = true
		}
	}
	//取剩余份数
	platformExp, err := db.GetPlatformExperience(a.dbEngine)
	if err != nil {
		logrus.Info("query db error", err)
		res := util.ResponseMsg(-1, "query db error", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}
	if platformExp != nil {
		body["balance"] = platformExp.MaxPersons
	}
	//取注册领取状态
	userExps, err := db.GetUserExperience(a.dbEngine, uidFormatted)
	if err != nil {
		logrus.Info("query db error", err)
		res := util.ResponseMsg(-1, "query db error", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}
	for _, userExp := range userExps {
		switch userExp.ExpType {
		case "1":
			body["RegisterStatus"] = true
		case "2":
			body["isBindGoogleStatus"] = true
		case "3":
			body["apiBindingStatus"] = true
		}
	}

	res := util.ResponseMsg(0, "checkoutQualification success", body)
	c.SecureJSON(http.StatusOK, res)
	return
}

func userRegisterFund(a *ApiService, uidFormatted string) (error, *types.Users) {
	user, err := db.GetUser(a.dbEngine, uidFormatted)
	if err != nil {
		logrus.Info("query db error", err)
		return err, nil
	}
	if user == nil {
		logrus.Info("find no user")
		return errors.New("find no user"), nil
	}
	return nil, user
}

func userBind(a *ApiService, uidFormatted string) error {
	userBindInfos, err := db.GetUserBindInfos(a.dbEngine, uidFormatted)
	if err != nil {
		logrus.Info("query db error", err)
		return err
	}
	if len(userBindInfos.ApiKey) == 0 || len(userBindInfos.ApiSecret) == 0 {
		logrus.Info("find user bind info,but apikey or apiSecret is null", uidFormatted)
		return errors.New("find user bind info,but apikey or apiSecret is null")
	}
	return nil
}

func (a *ApiService) getExperienceFund(c *gin.Context) {
	uid, _ := c.Get("Uid")
	uidFormatted := fmt.Sprintf("%s", uid)

	experience := c.Query("type")
	experienceType := fmt.Sprintf("%s", experience)

	body := make(map[string]interface{})
	var user *types.Users
	if experienceType == "1" {
		// 注册校验
		var err error
		err, user = userRegisterFund(a, uidFormatted)
		if err != nil {
			res := util.ResponseMsg(-1, "fail", err)
			c.SecureJSON(http.StatusOK, res)
			return
		}
	} else if experienceType == "2" {
		//谷歌绑定校验
		var err error
		err, user = db.QuerySecret(a.dbEngine, uidFormatted)
		if err != nil {
			res := util.ResponseMsg(-1, "fail", err)
			c.SecureJSON(http.StatusOK, res)
			return
		}
		if user.IsBindGoogle == false {
			logrus.Info("google is not bind", user.IsBindGoogle)
			res := util.ResponseMsg(-1, "google is not bind", body)
			c.SecureJSON(http.StatusOK, res)
			return
		}
	} else if experienceType == "3" {
		err := userBind(a, uidFormatted)
		if err != nil {
			res := util.ResponseMsg(-1, "api is not bind", err)
			c.SecureJSON(http.StatusOK, res)
			return
		}
	}

	logrus.Info("condition is satisfied,can get money", uid)

	receiveSum := map[string]string{
		"1": "30",
		"2": "30",
		"3": "40",
	}
	sum, err := strconv.ParseInt(receiveSum[experienceType], 10, 64)
	if err != nil {
		res := util.ResponseMsg(-1, "ParseInt error", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}

	session := a.dbEngine.NewSession()
	err = session.Begin()
	if err != nil {
		logrus.Info(err)

		res := util.ResponseMsg(-1, "session Begin error", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}
	////平台体验金资金池减少相应的数额
	//TotalRevenueInfo, err := db.GetPlatformExperience(a.dbEngine)
	//
	//if err != nil {
	//	logrus.Info("query db error", err)
	//
	//	res := util.ResponseMsg(-1, "query db error", err)
	//	c.SecureJSON(http.StatusOK, res)
	//	return
	//}
	//
	//if TotalRevenueInfo == nil {
	//	logrus.Info("platform exp info not exist")
	//
	//	res := util.ResponseMsg(-1, "platform exp info not exist", nil)
	//	c.SecureJSON(http.StatusOK, res)
	//	return
	//}
	//
	////总金额减少
	//TotalRevenueInfo.TotalSum = TotalRevenueInfo.TotalSum - sum
	////接收人数+1
	//TotalRevenueInfo.ReceivePersons = TotalRevenueInfo.ReceivePersons + 1
	////更新
	//_, err = session.Table("platformExperience").Update(TotalRevenueInfo)
	//if err != nil {
	//	err := session.Rollback()
	//	if err != nil {
	//		logrus.Error(err)
	//
	//		res := util.ResponseMsg(0, "internal db session rollback error", err)
	//		c.SecureJSON(http.StatusOK, res)
	//		return
	//	}
	//}

	//首先用户体验金增加一条记录
	userExperience := types.UserExperience{}
	timeNow := time.Now()
	sevenDayAgo := timeNow.AddDate(0, 0, 7)
	userExperience.UId = uidFormatted
	userExperience.Type = "1" // 新人有礼
	userExperience.ExpType = experienceType
	userExperience.ReceiveSum = sum
	userExperience.CoinName = "usdt"
	userExperience.ValidTime = timeNow
	userExperience.ValidStartTime = timeNow
	userExperience.ValidStartTime = sevenDayAgo
	userExperience.Status = false

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

func (a *ApiService) getExperience(c *gin.Context) {
	uid, _ := c.Get("Uid")
	uidFormatted := fmt.Sprintf("%s", uid)

	userExps, err := db.GetUserExperience(a.dbEngine, uidFormatted)

	if err != nil {
		logrus.Info("query db error", err)

		res := util.ResponseMsg(-1, "query db error", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}
	var receiveSum int64 = 0
	for _, userExp := range userExps {
		now := time.Now()
		if userExp.Status == false && now.Before(userExp.ValidEndTime) {
			receiveSum = receiveSum + userExp.ReceiveSum
		}
	}

	body := make(map[string]interface{})
	body["receiveSum"] = receiveSum
	res := util.ResponseMsg(0, "query exp success", body)

	c.SecureJSON(http.StatusOK, res)
	return
}

// 执行系统策略-1.首先检查条件 2.将平台剩余份数-1 3.将用户表的使用状态更新过来
func (a *ApiService) extcuteSystemStrategy(c *gin.Context) {
	uid, _ := c.Get("Uid")
	uidFormatted := fmt.Sprintf("%s", uid)

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
	if user.IsBindGoogle == false {
		res := util.ResponseMsg(-1, "no condition:IsBindGoogle error", user.IsBindGoogle)
		c.SecureJSON(http.StatusOK, res)
		return
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
		if len(userBindInfos.ApiKey) <= 0 || len(userBindInfos.ApiSecret) <= 0 {
			res := util.ResponseMsg(-1, "no condition:BindApi error", nil)
			c.SecureJSON(http.StatusOK, res)
			return
		}
	}
	logrus.Info("meet condition")

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

	session := a.dbEngine.NewSession()
	err = session.Begin()
	if err != nil {
		logrus.Info(err)

		res := util.ResponseMsg(-1, "session Begin error", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}

	TotalRevenueInfo.MaxPersons = TotalRevenueInfo.MaxPersons - 1

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

	//2.将用户表的使用状态更新过来
	userExpUpdate := types.UserExpUpdate{}
	userExpUpdate.Status = "t"

	_, err = session.Table("userExperience").Where("f_uid=?", uidFormatted).Update(userExpUpdate)
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

	res := util.ResponseMsg(0, "UpdatePlatformExp success", nil)
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
