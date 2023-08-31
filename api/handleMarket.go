package api

import (
	"github.com/ethereum/BGService/types"
	"github.com/ethereum/BGService/util"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

var base_binance_url = "https://api.binance.com/"

//var base_ok_url = "https://api.binance.com/"

// 默认展示币安交易所的行情
// 这里交给交易所直接校验
func (a *ApiService) getBinancePrice(c *gin.Context) {
	symbols := c.Query("symbols")

	res := types.HttpRes{}

	url := base_biance_url + "api/v3/ticker/price?symbols=" + symbols

	data, err := util.Get(url)
	if err != nil {
		logrus.Info("获取币价失败", err)

		res.Code = 0
		res.Message = "成功获取价格"
		res.Data = err

		c.SecureJSON(http.StatusOK, res)
		return
	}
	res.Code = 0
	res.Message = "成功获取价格"
	res.Data = data

	c.SecureJSON(http.StatusOK, res)
	return
}

// 这里交给交易所直接校验
func (a *ApiService) getBinance24hInfos(c *gin.Context) {
	symbols := c.Query("symbols")

	res := types.HttpRes{}

	url := base_biance_url + "/api/v3/ticker/24hr?symbols=" + symbols

	data, err := util.Get(url)
	if err != nil {
		logrus.Info("获取24小时涨跌失败", err)

		res.Code = 0
		res.Message = "获取24小时涨跌失败"
		res.Data = err

		c.SecureJSON(http.StatusOK, res)
		return
	}
	res.Code = 0
	res.Message = "获取24小时涨跌成功"
	res.Data = data

	c.SecureJSON(http.StatusOK, res)
	return
}
