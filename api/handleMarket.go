package api

import (
	"github.com/ethereum/BGService/db"
	"github.com/ethereum/BGService/types"
	"github.com/ethereum/BGService/util"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"strings"
)

var base_binance_url = "https://api.binance.com/"

// var base_ok_url = "https://api.binance.com/"

//var base_future_binance_url = "https://fapi.binance.com"

var base_cmc_url = "https://pro-api.coinmarketcap.com/"

var test_cmc_key = "b6f2d5f6-21c5-4a54-a61a-e85853d8a043"

// 默认展示币安交易所的行情
// 这里交给交易所直接校验
func (a *ApiService) getBinancePrice(c *gin.Context) {
	symbols := c.Query("symbols")

	res := types.HttpRes{}

	url := base_binance_url + "api/v3/ticker/price?symbols=" + symbols

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

	url := base_binance_url + "/api/v3/ticker/24hr?symbols=" + symbols

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

// 这里交给CMC直接校验
func (a *ApiService) getCoinInfos(c *gin.Context) {
	//symbols := c.Query("symbols")

	res := types.HttpRes{}

	cmcUrl := base_cmc_url + "v1/cryptocurrency/map"

	params := url.Values{}
	params.Add("symbol", "BTC")

	data, err := util.GetWithDataHeader(cmcUrl, params, test_cmc_key)
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

// 将币对添加/移除个人自选
func (a *ApiService) addConcern(c *gin.Context) {
	var userConcern types.UserConcern
	res := types.HttpRes{}

	err := c.BindJSON(&userConcern)
	if err != nil {
		logrus.Info("传递的不是合法的json参数")

		res.Code = -1
		res.Message = "传递的不是合法的json参数"
		res.Data = err
		c.SecureJSON(http.StatusOK, res)
		return
	}
	uid := userConcern.Uid
	coinPair := userConcern.CoinPair
	method := userConcern.Method

	//参数校验
	if method != "add" && method != "remove" {
		res.Code = -1
		res.Message = "method can not support"
		c.SecureJSON(http.StatusOK, res)
		return
	}

	if strings.HasPrefix(coinPair, ",") || strings.HasSuffix(coinPair, ",") {
		res.Code = -1
		res.Message = "coinPair can not start or end with comma"
		c.SecureJSON(http.StatusOK, res)
		return
	}

	user, err := db.GetUser(a.dbEngine, uid)

	if err != nil {
		logrus.Info("query db error:", err)

		res.Code = -1
		res.Message = "query db error"
		res.Data = err
		c.SecureJSON(http.StatusOK, res)
		return
	}

	if user == nil {
		logrus.Info("no user record:", uid)

		res.Code = -1
		res.Message = "no user record"
		c.SecureJSON(http.StatusOK, res)
		return
	}
	var concern []string

	if user.ConcernCoinList != "{}" && len(user.ConcernCoinList) != 0 {
		concern = strings.Split(user.ConcernCoinList[1:len(user.ConcernCoinList)-1], ",")
		logrus.Info(concern)
		if method == "add" {
			concern = append(concern, coinPair)
			logrus.Info(concern)
		} else {
			//首先找到这个remove位置，找不到返回错误，找到按照这个位置remove
			find := false
			for index, value := range concern {
				if value == coinPair {
					concern = append(concern[:index], concern[index+1:]...)
					find = true
					break
				}
			}
			if find == false {
				res.Code = -1
				res.Message = "can not find remove record"
				c.SecureJSON(http.StatusOK, res)
				return
			}
		}
	} else {
		if method == "add" {
			concern = append(concern, coinPair)
		} else {
			res.Code = -1
			res.Message = "null list can not remove anything"
			c.SecureJSON(http.StatusOK, res)
			return
		}
	}

	concernStr := "{"
	length := len(concern)
	//再将这个数组转化为字串，更新数据库
	for index, value := range concern {
		concernStr = concernStr + value

		if index+1 < length {
			concernStr = concernStr + ","
		}
	}
	concernStr = concernStr + "}"

	user.ConcernCoinList = concernStr

	err = db.UpdateUser(a.dbEngine, uid, user)
	if err != nil {
		logrus.Info("update user concern:", err)

		res.Code = 0
		res.Message = "update user concern"
		res.Data = err

		c.SecureJSON(http.StatusOK, res)
		return
	}

	res.Code = 0
	res.Message = "add or remove concern success"

	c.SecureJSON(http.StatusOK, res)
	return
}

// 得到交易账户详情
func (a *ApiService) getTradeAccountDetail(c *gin.Context) {
	accountTotalAssets := make(map[string]string)
	initAssets := make(map[string]string)
	todayBenefits := make(map[string]string)

	var tradeDetails types.TradeDetails

	//首先得到我的策略
	uid := c.Query("uid")

	//首先得到我的仓位
	userData, err := util.GetBinanceUserData()

	if err != nil { //经常报 Timestamp for this request is outside of the recvWindow.
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}

	//查询用户策略表得到用户对应得所有策略
	userStrategys, err := db.GetUserStrategys(a.dbEngine, uid)
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}

	//一个稳定币只可能存在一个策略
	for _, userStrategy := range userStrategys {
		//查询量化收益表
		latestEarning, err := db.GetUserStrategyLatestEarnings(a.dbEngine, uid, userStrategy.StrategyID)

		if err != nil {
			res := util.ResponseMsg(-1, "fail", err)
			c.SecureJSON(http.StatusOK, res)
			return
		}

		cexTotalProfit, err := decimal.NewFromString(userData.TotalUnrealizedProfit)
		if err != nil {
			res := util.ResponseMsg(-1, "fail", err)
			c.SecureJSON(http.StatusOK, res)
			return
		}
		dbTotalBenefit, err := decimal.NewFromString(latestEarning.TotalBenefit)
		if err != nil {
			res := util.ResponseMsg(-1, "fail", err)
			c.SecureJSON(http.StatusOK, res)
			return
		}

		dayBefinit := cexTotalProfit.Sub(dbTotalBenefit)

		//查询策略表
		strategyInfo, err := db.GetStrategy(a.dbEngine, userStrategy.StrategyID)
		if err != nil {
			res := util.ResponseMsg(-1, "fail", err)
			c.SecureJSON(http.StatusOK, res)
			return
		}

		if strings.Contains(strings.ToLower(strategyInfo.StrategyName), "usdt") == true {
			for _, asset := range userData.Assets {
				if asset.Asset == "usdt" || asset.Asset == "USDT" {
					accountTotalAssets["usdt"] = asset.MarginBalance
					initAssets["usdt"] = userStrategy.ActualInvest
					todayBenefits["usdt"] = dayBefinit.String()
				}
			}
		}
		if strings.Contains(strings.ToLower(strategyInfo.StrategyName), "usdc") == true {
			for _, asset := range userData.Assets {
				if asset.Asset == "usdc" || asset.Asset == "USDC" {
					accountTotalAssets["usdc"] = asset.MarginBalance
					initAssets["usdc"] = userStrategy.ActualInvest
					todayBenefits["usdc"] = dayBefinit.String()
				}
			}
		}
		if strings.Contains(strings.ToLower(strategyInfo.StrategyName), "busd") == true {
			for _, asset := range userData.Assets {
				if asset.Asset == "busd" || asset.Asset == "BUSD" {
					accountTotalAssets["busd"] = asset.MarginBalance
					initAssets["busd"] = userStrategy.ActualInvest
					todayBenefits["busd"] = dayBefinit.String()
				}
			}
		}
	}

	tradeDetails.AccountTotalAssets = accountTotalAssets
	tradeDetails.InitAssets = initAssets
	tradeDetails.CurBenefit = todayBenefits

	res := util.ResponseMsg(0, "getTradeDetails success", tradeDetails)
	c.SecureJSON(http.StatusOK, res)
	return
}
