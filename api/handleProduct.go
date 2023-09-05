package api

import (
	"github.com/ethereum/BGService/db"
	"github.com/ethereum/BGService/types"
	"github.com/ethereum/BGService/util"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
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

func (a *ApiService) productList(c *gin.Context) {
	var payload *types.StrategyInput
	if err := c.ShouldBindJSON(&payload); err != nil {
		logrus.Error(err)
		res := util.ResponseMsg(-1, "fail", err.Error())
		c.SecureJSON(http.StatusOK, res)
		return
	}
	db.GetScreenStrategy(a.dbEngine, payload)
	body := make(map[string]interface{})
	res := util.ResponseMsg(1, "success", body)
	c.SecureJSON(http.StatusOK, res)
	return
}
