package api

import (
	"github.com/ethereum/BGService/db"
	"github.com/ethereum/BGService/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (a *ApiService) overview(c *gin.Context) {
	allStrategy, err := db.GetAllStrategy(a.dbEngine)
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err.Error())
		c.SecureJSON(http.StatusOK, res)
		return
	}
	totalAssets, err := db.GetStrategyTotalAssets(a.dbEngine)
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err.Error())
		c.SecureJSON(http.StatusOK, res)
		return
	}
	UserCount, err := db.GetStrategyUserCount(a.dbEngine)
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err.Error())
		c.SecureJSON(http.StatusOK, res)
		return
	}
	UserIncome, err := db.GetUserIncome(a.dbEngine)
	if err != nil {
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
