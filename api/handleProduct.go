package api

import (
	"fmt"
	"github.com/ethereum/BGService/db"
	"github.com/ethereum/BGService/types"
	"github.com/ethereum/BGService/util"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func (a *ApiService) overview(c *gin.Context) {
	// 运行中的策略
	allStrategy, err := db.GetAllStrategy(a.dbEngine)
	if err != nil {
		logrus.Error(err)
		res := util.ResponseMsg(-1, "fail", err.Error())
		c.SecureJSON(http.StatusOK, res)
		return
	}
	// 总资产
	totalAssets, err := db.GetStrategyTotalAssets(a.dbEngine)
	if err != nil {
		logrus.Error(err)
		res := util.ResponseMsg(-1, "fail", err.Error())
		c.SecureJSON(http.StatusOK, res)
		return
	}
	// 量化交易全球用户数
	UserCount, err := db.GetStrategyUserCount(a.dbEngine)

	if err != nil {
		logrus.Error(err)
		res := util.ResponseMsg(-1, "fail", err.Error())
		c.SecureJSON(http.StatusOK, res)
		return
	}
	// 量化用户累计收益
	UserIncome, err := db.GetUserIncome(a.dbEngine)
	if err != nil {
		logrus.Error(err)
		res := util.ResponseMsg(-1, "fail", err.Error())
		c.SecureJSON(http.StatusOK, res)
		return
	}
	body := make(map[string]interface{})
	body["runStrategy"] = len(allStrategy)
	body["totalAssets"] = totalAssets
	body["globalUserCount"] = UserCount
	body["globalUserIncome"] = UserIncome

	res := util.ResponseMsg(1, "success", body)
	c.SecureJSON(http.StatusOK, res)
	return
}

func isInCollectStrategyList(element string, collectStrategyList []string) bool {
	for _, item := range collectStrategyList {
		if item == element {
			return true
		}
	}
	return false
}

func (a *ApiService) productList(c *gin.Context) {
	var payload *types.StrategyInput
	if err := c.ShouldBindJSON(&payload); err != nil {
		logrus.Error(err)
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}
	var CollectStragetyList []string
	if payload.Currency == "1" {
		session := sessions.Default(c)
		uid := session.Get("Uid")
		if uid == nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			return
		}
		uidFormatted := fmt.Sprintf("%s", uid)
		user, err := db.GetUser(a.dbEngine, uidFormatted)
		if err != nil {
			res := util.ResponseMsg(-1, "fail", err)
			c.SecureJSON(http.StatusOK, res)
			return
		}
		CollectStragetyList = strings.Split(user.CollectStragetyList[1:len(user.CollectStragetyList)-1], ",")
	}
	var ScreenStrategys []types.Strategy
	if payload.Keywords != "" {
		// 模糊搜索
		var err error
		ScreenStrategys, err = db.GetSearchScreenStrategy(a.dbEngine, payload)
		if err != nil {
			logrus.Error(err)
			res := util.ResponseMsg(-1, "fail", err.Error())
			c.SecureJSON(http.StatusOK, res)
			return
		}
	} else {
		// 筛选
		var err error
		ScreenStrategys, err = db.GetScreenStrategy(a.dbEngine, payload, CollectStragetyList)
		if err != nil {
			logrus.Error(err)
			res := util.ResponseMsg(-1, "fail", err.Error())
			c.SecureJSON(http.StatusOK, res)
			return
		}
	}

	var ScreenStrategyList []interface{}

	ScreenStrategy := make(map[string]interface{})
	var isCollect = false
	for _, value := range ScreenStrategys {
		ScreenStrategy["id"] = value.StrategyID
		ScreenStrategy["name"] = value.StrategyName
		ScreenStrategy["productCategory"] = value.Type
		ScreenStrategy["recommendRate"] = value.RecommendRate
		if payload.Currency == "1" {
			isCollect = isInCollectStrategyList(value.StrategyID, CollectStragetyList)
		}
		ScreenStrategy["isCollect"] = isCollect
		ScreenStrategy["participateNum"] = value.ParticipateNum
		ScreenStrategy["totalYield"] = value.TotalYield
		ScreenStrategy["runTime"] = value.CreateTime
		ScreenStrategy["maxWithdrawalRate"] = value.MaxDrawDown
		ScreenStrategy["minimumInvestmentAmount"] = value.MinInvest
		ScreenStrategy["strategySource"] = value.Source
		ScreenStrategyList = append(ScreenStrategyList, ScreenStrategy)
	}
	body := make(map[string]interface{})
	body["list"] = ScreenStrategyList
	body["total"] = len(ScreenStrategyList)
	res := util.ResponseMsg(1, "success", body)
	c.SecureJSON(http.StatusOK, res)
	return
}
