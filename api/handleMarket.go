package api

import (
	"fmt"
	"github.com/LinkinStars/go-scaffold/contrib/cryptor"
	"github.com/ethereum/BGService/db"
	"github.com/ethereum/BGService/types"
	"github.com/ethereum/BGService/util"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var base_binance_url = "https://api.binance.com/"

// var base_ok_url = "https://api.binance.com/"

//var base_future_binance_url = "https://fapi.binance.com"

//var base_cmc_url = "https://pro-api.coinmarketcap.com/"
//
//var test_cmc_key = "b6f2d5f6-21c5-4a54-a61a-e85853d8a043"

// 将币对添加/移除个人自选
func (a *ApiService) addConcern(c *gin.Context) {
	var userConcern types.UserConcern
	Uid, _ := c.Get("Uid")
	// 根据uid查询用户信息
	uidFormatted := fmt.Sprintf("%s", Uid)

	err := c.BindJSON(&userConcern)
	if err != nil {
		logrus.Info("not valid json")

		res := util.ResponseMsg(-1, "not valid json", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}
	uid := uidFormatted
	coinPair := strings.ToLower(userConcern.CoinPair)
	method := strings.ToLower(userConcern.Method)

	//参数校验
	if method != "add" && method != "remove" {
		logrus.Info("method can not support")

		res := util.ResponseMsg(-1, "method can not support", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}

	if strings.HasPrefix(coinPair, ",") || strings.HasSuffix(coinPair, ",") {
		logrus.Info("coinPair can not start or end with comma")

		res := util.ResponseMsg(-1, "coinPair can not start or end with comma", nil)
		c.SecureJSON(http.StatusOK, res)
		return
	}

	user, err := db.GetUser(a.dbEngine, uid)

	if err != nil {
		logrus.Info("query db error:", err)

		res := util.ResponseMsg(-1, "query db error", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}

	if user == nil {
		logrus.Info("no user record:", uid)

		res := util.ResponseMsg(-1, "no user record:", uid)
		c.SecureJSON(http.StatusOK, res)
		return
	}
	var concern []string

	if user.ConcernCoinList != "{}" && len(user.ConcernCoinList) != 0 {
		concern = strings.Split(user.ConcernCoinList[1:len(user.ConcernCoinList)-1], ",")
		logrus.Info(concern)
		if method == "add" {
			//判断，该币种列表是否已经存在
			if strings.Contains(user.ConcernCoinList, coinPair) == true {
				res := util.ResponseMsg(-1, "coinPair is already exist", nil)
				c.SecureJSON(http.StatusOK, res)
				return
			}

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
				res := util.ResponseMsg(-1, "can not find remove record", nil)
				c.SecureJSON(http.StatusOK, res)
				return
			}
		}
	} else {
		if method == "add" {
			concern = append(concern, coinPair)
		} else {
			res := util.ResponseMsg(-1, "null list can not remove anything", nil)
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

		res := util.ResponseMsg(-1, "update user concern", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}

	res := util.ResponseMsg(0, "add or remove concern success", nil)
	c.SecureJSON(http.StatusOK, res)
	return
}

// 得到我得自选
func (a *ApiService) getCoinInfo(c *gin.Context) {
	res := util.ResponseMsg(0, "getCoinInfo nil for now", nil)
	c.SecureJSON(http.StatusOK, res)
	return
}

// 得到我得自选
func (a *ApiService) getConcern(c *gin.Context) {
	Uid, _ := c.Get("Uid")
	// 根据uid查询用户信息
	uidFormatted := fmt.Sprintf("%s", Uid)

	user, err := db.GetUser(a.dbEngine, uidFormatted)

	if err != nil {
		logrus.Info("query db error:", err)

		res := util.ResponseMsg(-1, "query db error", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}

	if user == nil {
		logrus.Info("no user record:", uidFormatted)

		res := util.ResponseMsg(-1, "no user record:", nil)
		c.SecureJSON(http.StatusOK, res)
		return
	}
	var concern []string

	if user.ConcernCoinList != "{}" && len(user.ConcernCoinList) != 0 {
		concern = strings.Split(user.ConcernCoinList[1:len(user.ConcernCoinList)-1], ",")
		logrus.Info(concern)
	}
	body := make(map[string]interface{})
	body["list"] = concern
	res := util.ResponseMsg(0, "getConcern success", body)
	c.SecureJSON(http.StatusOK, res)
	return
}

// 得到特定策略的信息--总收益 总收益率 运行时长--查询用户策略收益表，统计这个策略的信息
func (a *ApiService) getKlinesHistory(c *gin.Context) {
	interval := c.Query("interval")
	startTimeParam := c.Query("startTime")
	endTimeParam := c.Query("endTime")
	KlineTypeParam := c.Query("KlineType")
	symbol := c.Query("symbol")

	symbol = strings.ToUpper(symbol)

	startTime, err := strconv.ParseInt(startTimeParam, 10, 64)
	if err != nil {
		res := util.ResponseMsg(-1, "startTime ParseInt fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}

	endTime, err := strconv.ParseInt(endTimeParam, 10, 64)
	if err != nil {
		res := util.ResponseMsg(-1, "endTime ParseInt fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}

	KlineType, err := strconv.Atoi(KlineTypeParam)
	if err != nil {
		res := util.ResponseMsg(-1, "KlineType ParseInt fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}

	klines, err := util.GetBinanceKlinesHistory(interval, startTime, endTime, KlineType, symbol)

	if err != nil {
		res := util.ResponseMsg(-1, "GetKlinesHistory fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}

	res := util.ResponseMsg(0, "GetKlinesHistory success", klines)
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

// 得到交易账户列表--遍历我的策略产品列表--
func (a *ApiService) getTradeList(c *gin.Context) {
	var tradeDetails types.TradeDetails
	var tradeList []types.TradeDetails
	//首先得到我的策略
	uid, _ := c.Get("Uid")
	uidFormatted := fmt.Sprintf("%s", uid)

	//由UID得到用户的APIKEY

	//查询用户策略表得到用户对应得所有策略
	userBind, err := db.GetUserBindInfos(a.dbEngine, uidFormatted)
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}
	//这里再进行解密
	//先解密再使用
	apiKey := cryptor.AesSimpleDecrypt(userBind.ApiKey, types.AesKey)
	apiSecret := cryptor.AesSimpleDecrypt(userBind.ApiSecret, types.AesKey)

	//首先得到我的仓位
	userData, err := util.GetBinanceUMUserData(apiKey, apiSecret)

	if err != nil { //经常报 Timestamp for this request is outside of the recvWindow.
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}

	//查询用户策略表得到用户对应得所有策略
	userStrategys, err := db.GetUserStrategys(a.dbEngine, uidFormatted)
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}

	//一个稳定币只可能存在一个策略
	for _, userStrategy := range userStrategys {
		//查询量化收益表
		latestEarning, err := db.GetUserStrategyLatestEarnings(a.dbEngine, uidFormatted, userStrategy.StrategyID)

		if latestEarning == nil {
			res := util.ResponseMsg(-1, "no earning ", nil)
			c.SecureJSON(http.StatusOK, res)
			return
		}
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
		dec, err := decimal.NewFromString(latestEarning.TotalBenefit)
		if err != nil {
			res := util.ResponseMsg(-1, "fail", err)
			c.SecureJSON(http.StatusOK, res)
			return
		}

		dayBefinit := cexTotalProfit.Sub(dec)

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
					tradeDetails.AccountTotalAssets = asset.MarginBalance
					tradeDetails.InitAssets = userStrategy.ActualInvest
					tradeDetails.CurBenefit = dayBefinit.String()
					tradeDetails.Name = strategyInfo.StrategyName
				}
			}
		}
		if strings.Contains(strings.ToLower(strategyInfo.StrategyName), "usdc") == true {
			for _, asset := range userData.Assets {
				if asset.Asset == "usdc" || asset.Asset == "USDC" {
					tradeDetails.AccountTotalAssets = asset.MarginBalance
					tradeDetails.InitAssets = userStrategy.ActualInvest
					tradeDetails.CurBenefit = dayBefinit.String()
					tradeDetails.Name = strategyInfo.StrategyName
				}
			}
		}
		if strings.Contains(strings.ToLower(strategyInfo.StrategyName), "busd") == true {
			for _, asset := range userData.Assets {
				if asset.Asset == "busd" || asset.Asset == "BUSD" {
					tradeDetails.AccountTotalAssets = asset.MarginBalance
					tradeDetails.InitAssets = userStrategy.ActualInvest
					tradeDetails.CurBenefit = dayBefinit.String()
					tradeDetails.Name = strategyInfo.StrategyName
				}
			}
		}
		tradeDetails.ProductID = userStrategy.StrategyID
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

// 得到特定产品的详情
func (a *ApiService) getTradeDetail(c *gin.Context) {
	var tradeDetails types.TradeDetails
	//首先得到我的策略
	uid, _ := c.Get("Uid")
	uidFormatted := fmt.Sprintf("%s", uid)

	productID := c.Query("productID")
	//一期先不处理status

	userBind, err := db.GetUserBindInfos(a.dbEngine, uidFormatted)
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}
	//先解密再使用
	apiKey := cryptor.AesSimpleDecrypt(userBind.ApiKey, types.AesKey)
	apiSecret := cryptor.AesSimpleDecrypt(userBind.ApiSecret, types.AesKey)

	//首先得到我的仓位
	userData, err := util.GetBinanceUMUserData(apiKey, apiSecret)

	if err != nil { //经常报 Timestamp for this request is outside of the recvWindow.
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}

	//查询用户策略表得到具体的策略
	userStrategy, err := db.GetExactlyUserStrategy(a.dbEngine, uidFormatted, productID)
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}

	//一个稳定币只可能存在一个策略
	//查询量化收益表
	latestEarning, err := db.GetUserStrategyLatestEarnings(a.dbEngine, uidFormatted, userStrategy.StrategyID)

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
	dec, err := decimal.NewFromString(latestEarning.TotalBenefit)

	dayBefinit := cexTotalProfit.Sub(dec)

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
				tradeDetails.AccountTotalAssets = asset.MarginBalance
				tradeDetails.InitAssets = userStrategy.ActualInvest
				tradeDetails.CurBenefit = dayBefinit.String()
				tradeDetails.Name = strategyInfo.StrategyName
			}
		}
	}
	if strings.Contains(strings.ToLower(strategyInfo.StrategyName), "usdc") == true {
		for _, asset := range userData.Assets {
			if asset.Asset == "usdc" || asset.Asset == "USDC" {
				tradeDetails.AccountTotalAssets = asset.MarginBalance
				tradeDetails.InitAssets = userStrategy.ActualInvest
				tradeDetails.CurBenefit = dayBefinit.String()
				tradeDetails.Name = strategyInfo.StrategyName
			}
		}
	}
	if strings.Contains(strings.ToLower(strategyInfo.StrategyName), "busd") == true {
		for _, asset := range userData.Assets {
			if asset.Asset == "busd" || asset.Asset == "BUSD" {
				tradeDetails.AccountTotalAssets = asset.MarginBalance
				tradeDetails.InitAssets = userStrategy.ActualInvest
				tradeDetails.CurBenefit = dayBefinit.String()
				tradeDetails.Name = strategyInfo.StrategyName
			}
		}
	}
	tradeDetails.ProductID = userStrategy.StrategyID
	tradeDetails.InDays = time.Now().Sub(userStrategy.JoinTime).String()

	tradeDetails.Source = strategyInfo.Source
	tradeDetails.Type = strategyInfo.Type
	tradeDetails.ShareRatio = strategyInfo.ShareRatio
	tradeDetails.DividePeriod = strategyInfo.DividePeriod
	tradeDetails.AgreementPeriod = strategyInfo.AgreementPeriod

	res := util.ResponseMsg(0, "TradeDetails success", tradeDetails)
	c.SecureJSON(http.StatusOK, res)
	return
}

// 得到交易历史--todo:目前量化那边没有接口可以区分用户自己得交易记录和量化交易记录，等那边提供再增加区分逻辑
func (a *ApiService) getTradeHistory(c *gin.Context) {
	var transactionRecord types.TransactionRecord
	var transactionRecords []types.TransactionRecord
	//首先得到我的策略
	productID := c.Query("productID")

	uid, _ := c.Get("Uid")
	uidFormatted := fmt.Sprintf("%s", uid)

	strategy, err := db.GetStrategy(a.dbEngine, productID)
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}
	//
	//查询用户策略表得到用户对应得所有策略
	userBind, err := db.GetUserBindInfos(a.dbEngine, uidFormatted)
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}
	//这里再进行解密
	//先解密再使用
	apiKey := cryptor.AesSimpleDecrypt(userBind.ApiKey, types.AesKey)
	apiSecret := cryptor.AesSimpleDecrypt(userBind.ApiSecret, types.AesKey)

	symbol := util.RemoveElement(strategy.StrategyName, "/")

	userHistorys, err := util.GetBinanceUMUserTxHistory(apiKey, apiSecret, symbol, 1000)

	if err != nil { //经常报 Timestamp for this request is outside of the recvWindow.
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}

	for index, record := range userHistorys {
		transactionRecord.ID = index

		t := time.Unix(record.Time/1000, 0)
		transactionRecord.Time = t.Format("2006-01-02 15:04:05")

		if record.Buyer == true {
			transactionRecord.Action = "buy"
		} else {
			transactionRecord.Action = "sell"
		}
		transactionRecord.Behavior = "行为"
	}
	transactionRecords = append(transactionRecords, transactionRecord)

	res := util.ResponseMsg(0, "getTradeHistory success", transactionRecords)
	c.SecureJSON(http.StatusOK, res)
	return
}

func (a *ApiService) getUserDaysBenefit(c *gin.Context) {
	var userBenefitNDays types.UserBenefitNDays

	var userBenefits []types.UserBenefits

	sid := c.Query("sid")

	uidSession, _ := c.Get("Uid")
	uid := fmt.Sprintf("%s", uidSession)

	timeType := c.Query("timeType")

	startTime := time.Now()

	switch timeType {
	case "1":
		startTime = time.Now().AddDate(0, -1, 0)
	case "3":
		startTime = time.Now().AddDate(0, -3, 0)
	case "12":
		startTime = time.Now().AddDate(1, 0, 0)
	default:
		startTime = time.Now().AddDate(100, 0, 0)
	}

	earnings, err := db.GetStrategyBenefits(a.dbEngine, sid, uid, startTime.Format("2006-01-02 15:04:05"), time.Now().Format("2006-01-02 15:04:05"))
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

	sumRatio, err := decimal.NewFromString("0")
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}

	win := 0
	for _, earning := range earnings {
		var userBenefit types.UserBenefits
		darDec, err := decimal.NewFromString(earning.DayBenefit)
		if err != nil {
			res := util.ResponseMsg(-1, "fail", err)
			c.SecureJSON(http.StatusOK, res)
		}
		userBenefitNDays.BenefitSum = decimal.Sum(userBenefitNDays.BenefitSum, darDec)
		ratioDec, err := decimal.NewFromString(earning.DayRatio)
		if err != nil {
			res := util.ResponseMsg(-1, "fail", err)
			c.SecureJSON(http.StatusOK, res)
		}
		sumRatio = decimal.Sum(sumRatio, ratioDec)
		if ratioDec.IsPositive() { //收益率为正 胜利次数++
			win = win + 1
		}
		userBenefit.Date = earning.CreateTime.String()
		userBenefit.Benefit = earning.DayBenefit
		userBenefit.Ratio = earning.DayRatio
		userBenefits = append(userBenefits, userBenefit)
	}

	days := decimal.New(int64(len(earnings)), 32)
	userBenefitNDays.BenefitRatio = sumRatio.Div(days).String()

	//计算胜率
	length := len(earnings)
	dec1 := decimal.NewFromInt32(int32(win))
	dec2 := decimal.NewFromInt32(int32(length))

	userBenefitNDays.WinRatio = dec1.Div(dec2).String() //30日胜率

	//开始计算回撤率
	capital := userStrategy.ActualInvest //实际投资额

	maxEarning := earnings[0].DayBenefit        //30日最大收益
	minEarning := earnings[length-1].DayBenefit //30日最小收益

	maxDec, err := decimal.NewFromString(maxEarning)
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}
	minDec, err := decimal.NewFromString(minEarning)
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}
	capitalDec, err := decimal.NewFromString(capital)
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}

	//净值
	maxNetValue := decimal.Sum(capitalDec, maxDec)
	//计算回撤率：(最大收益-最小收益)/净值

	userBenefitNDays.Huiche = maxDec.Sub(minDec).Div(maxNetValue).String() //最大回撤率
	userBenefitNDays.Benefitlist = userBenefits

	res := util.ResponseMsg(0, "getUserDaysBenefit success", userBenefitNDays)
	c.SecureJSON(http.StatusOK, res)
	return
}

func (a *ApiService) getUserBeneiftInfo(c *gin.Context) {
	timeType := c.Query("timeType")
	sid := c.Query("sid")

	uidSession, _ := c.Get("Uid")
	uid := fmt.Sprintf("%s", uidSession)

	startTime := time.Now()

	switch timeType {
	case "1":
		startTime = time.Now().AddDate(0, -1, 0)
	case "3":
		startTime = time.Now().AddDate(0, -3, 0)
	case "12":
		startTime = time.Now().AddDate(1, 0, 0)
	default:
		startTime = time.Now().AddDate(100, 0, 0)
	}

	earnings, err := db.GetStrategyBenefits(a.dbEngine, sid, uid, startTime.Format("2006-01-02 15:04:05"), time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}
	res := util.ResponseMsg(0, "getUserBeneift success", earnings)
	c.SecureJSON(http.StatusOK, res)
	return
}

func (a *ApiService) getBinanceHighPercent(c *gin.Context) {
	ret, err := util.GetBinanceHighPercent()
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}
	res := util.ResponseMsg(0, "getUserBeneift success", ret)
	c.SecureJSON(http.StatusOK, res)
	return
}
