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
	"time"
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
		if len(userBindInfos.ApiKey) > 0 && len(userBindInfos.ApiSecret) > 0 {
			body["apiBinding"] = true
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
	// 查询三个条件-查询数据库
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
	session := a.dbEngine.NewSession()
	err := session.Begin()
	if err != nil {
		return
	}
	receiveSum := map[string]string{
		"1": "30",
		"2": "30",
		"3": "40",
	}

	//首先用户体验金增加一条记录
	userExperience := types.UserExperience{}
	timeNow := time.Now()
	sevenDayAgo := timeNow.AddDate(0, 0, 7)
	userExperience.UId = uidFormatted
	userExperience.Type = "1" // 新人有礼
	userExperience.ReceiveSum = receiveSum[experienceType]
	userExperience.CoinName = "USDT"
	userExperience.ValidTime = timeNow
	userExperience.ValidStartTime = timeNow
	userExperience.ValidStartTime = sevenDayAgo
	userExperience.Status = true

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

	res := util.ResponseMsg(0, "get exp success", body)
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
