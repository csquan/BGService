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
	"time"
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

// 得到特定策略的信息--总收益 总收益率 运行时长--查询用户策略收益表，统计这个策略的信息
func (a *ApiService) getStragetyDetail(c *gin.Context) {
	strategyName := c.Query("strategyName")

	strategy, err := db.GetStrategyByName(a.dbEngine, strategyName)
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
	}
	totalBenefit, err := db.GetStrategyTotalBenefits(a.dbEngine, strategy.StrategyID)
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
	}
	totalInvest, err := db.GetStrategyTotalInvests(a.dbEngine, strategy.StrategyID)
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
	}
	dec1 := decimal.NewFromFloat(totalBenefit)
	dec2 := decimal.NewFromFloat(totalInvest)
	ratio := dec1.Div(dec2)

	var strategyStats types.StrategyStats

	strategyStats.TotalBenefit = dec1.String()
	strategyStats.TotalRatio = ratio.String()
	strategyStats.RunTime = time.Now().Sub(strategy.CreateTime).String()

	res := util.ResponseMsg(0, "getStragetyDetail success", strategyStats)
	c.SecureJSON(http.StatusOK, res)
	return
}

// 得到交易账户列表--遍历我的策略产品列表
func (a *ApiService) getTradeList(c *gin.Context) {
	accountTotalAssets := make(map[string]string)
	initAssets := make(map[string]string)
	todayBenefits := make(map[string]string)

	var tradeDetails types.TradeDetails
	var tradeList []types.TradeDetails
	//首先得到我的策略
	uid := c.Query("uid")

	//首先得到我的仓位
	userData, err := util.GetBinanceUMUserData()

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

		dayBefinit := cexTotalProfit.Sub(latestEarning.TotalBenefit)

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
					initAssets["usdt"] = userStrategy.ActualInvest.String()
					todayBenefits["usdt"] = dayBefinit.String()
				}
			}
		}
		if strings.Contains(strings.ToLower(strategyInfo.StrategyName), "usdc") == true {
			for _, asset := range userData.Assets {
				if asset.Asset == "usdc" || asset.Asset == "USDC" {
					accountTotalAssets["usdc"] = asset.MarginBalance
					initAssets["usdc"] = userStrategy.ActualInvest.String()
					todayBenefits["usdc"] = dayBefinit.String()
				}
			}
		}
		if strings.Contains(strings.ToLower(strategyInfo.StrategyName), "busd") == true {
			for _, asset := range userData.Assets {
				if asset.Asset == "busd" || asset.Asset == "BUSD" {
					accountTotalAssets["busd"] = asset.MarginBalance
					initAssets["busd"] = userStrategy.ActualInvest.String()
					todayBenefits["busd"] = dayBefinit.String()
				}
			}
		}
		tradeDetails.AccountTotalAssets = accountTotalAssets
		tradeDetails.InitAssets = initAssets
		tradeDetails.CurBenefit = todayBenefits
		tradeDetails.InDays = time.Now().Sub(userStrategy.JoinTime).String()

		tradeDetails.Source = strategyInfo.Source
		tradeDetails.Type = strategyInfo.Type
		tradeDetails.ShareRatio = strategyInfo.ShareRatio
		tradeDetails.DividePeriod = strategyInfo.DividePeriod
		tradeDetails.AgreementPeriod = strategyInfo.AgreementPeriod

		tradeList = append(tradeList, tradeDetails)
	}

	res := util.ResponseMsg(0, "getTradeList success", tradeList)
	c.SecureJSON(http.StatusOK, res)
	return
}

// 得到交易历史--todo:目前量化那边没有接口可以区分用户自己得交易记录和量化交易记录，等那边提供再增加区分逻辑
func (a *ApiService) getTradeHistory(c *gin.Context) {
	//首先得到我的策略
	pairName := c.Query("pairName")

	symbol := util.RemoveElement(pairName, "/")

	userHistory, err := util.GetBinanceUMUserTxHistory(symbol, 1000)

	if err != nil { //经常报 Timestamp for this request is outside of the recvWindow.
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}

	res := util.ResponseMsg(0, "getTradeHistory success", userHistory)
	c.SecureJSON(http.StatusOK, res)
	return
}

func (a *ApiService) getUser30Beneift(c *gin.Context) {
	var userBenefit30Days types.UserBenefit30Days
	//todo 取出用户每日收益-得到当前日期的前30天内最高和最低的收益
	sid := c.Query("sid")
	uid := c.Query("uid")
	startTime := time.Now().AddDate(0, -1, 0)

	//取出30天按照时间倒序排序的收益
	earnings, err := db.GetStrategyBenefits(a.dbEngine, sid, uid, startTime.String(), time.Now().String())
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}
	//取出用户策略的实际投资额
	userStrategy, err := db.GetExactlyUserStrategy(a.dbEngine, uid, sid)
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}

	sumRatio, err := decimal.NewFromString("")
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}

	win := 0
	for _, earning := range earnings {
		userBenefit30Days.BenefitSum = decimal.Sum(userBenefit30Days.BenefitSum, earning.DayBenefit)
		sumRatio = decimal.Sum(sumRatio, earning.DayRatio)

		if earning.DayRatio.IsPositive() { //收益率为正 胜利次数++
			win = win + 1
		}
	}

	days := decimal.New(int64(len(earnings)), 32)
	userBenefit30Days.BenefitRatio = sumRatio.Div(days).String()

	//计算胜率
	length := len(earnings)
	dec1 := decimal.NewFromInt32(int32(win))
	dec2 := decimal.NewFromInt32(int32(length))

	userBenefit30Days.WinRatio = dec1.Div(dec2).String() //30日胜率

	//开始计算回撤率
	capital := userStrategy.ActualInvest //实际投资额

	maxEarning := earnings[0].DayBenefit        //30日最大收益
	minEarning := earnings[length-1].DayBenefit //30日最小收益

	//净值
	maxNetValue := decimal.Sum(capital, maxEarning)
	//计算回撤率：(最大收益-最小收益)/净值
	userBenefit30Days.Huiche = maxEarning.Sub(minEarning).Div(maxNetValue).String() //30日最大回撤率

	res := util.ResponseMsg(0, "getUser30Beneift success", userBenefit30Days)
	c.SecureJSON(http.StatusOK, res)
	return
}

func (a *ApiService) getUserBeneiftInfo(c *gin.Context) {
	timeType := c.Query("timeType")
	sid := c.Query("sid")
	uid := c.Query("uid")

	startTime := time.Now().String()

	switch timeType {
	case "1":
		startTime = time.Now().AddDate(0, -1, 0).Format("2006-01-02")
	case "2":
		startTime = time.Now().AddDate(0, -3, 0).Format("2006-01-02")
	case "3":
		startTime = time.Now().AddDate(1, 0, 0).Format("2006-01-02")
	default:
		startTime = time.Now().AddDate(100, 0, 0).Format("2006-01-02")
	}
	//todo 取出用户每日收益
	earnings, err := db.GetStrategyBenefits(a.dbEngine, sid, uid, startTime, time.Now().String())
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}
	res := util.ResponseMsg(0, "getUserBeneift success", earnings)
	c.SecureJSON(http.StatusOK, res)
	return
}
