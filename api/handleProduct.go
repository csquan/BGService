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
	"strconv"
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

func strToInt(strList []string) []int {
	intList := make([]int, len(strList))
	for i, str := range strList {
		num, err := strconv.Atoi(str)
		if err != nil {
			logrus.Error("无法将字符串转换为整数：", str)
			return []int{}
		}
		intList[i] = num
	}
	return intList
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
	var CollectStragetyListInt []int
	session := sessions.Default(c)
	uid := session.Get("Uid")
	if uid != nil {
		// 登录状态
		uidFormatted := fmt.Sprintf("%s", uid)
		user, err := db.GetUser(a.dbEngine, uidFormatted)
		if err != nil {
			res := util.ResponseMsg(-1, "fail", err)
			c.SecureJSON(http.StatusOK, res)
			return
		}
		CollectStragetyList = strings.Split(user.CollectStragetyList[1:len(user.CollectStragetyList)-1], ",")
		CollectStragetyListInt = strToInt(CollectStragetyList)
	} else if payload.Strategy == "1" {
		res := util.ResponseMsg(-1, "fail", "Please log in")
		c.SecureJSON(http.StatusOK, res)
		return
	}

	var ScreenStrategys []types.Strategy
	if payload.Keywords != "" {
		// 模糊搜索
		var err error
		ScreenStrategys, err = db.GetSearchScreenStrategy(a.dbEngine, payload, CollectStragetyListInt)
		if err != nil {
			logrus.Error(err)
			res := util.ResponseMsg(-1, "fail", err.Error())
			c.SecureJSON(http.StatusOK, res)
			return
		}
	} else {
		// 筛选
		var err error
		ScreenStrategys, err = db.GetScreenStrategy(a.dbEngine, payload, CollectStragetyListInt)
		if err != nil {
			logrus.Error(err)
			res := util.ResponseMsg(-1, "fail", err.Error())
			c.SecureJSON(http.StatusOK, res)
			return
		}
	}

	var ScreenStrategyList []interface{}

	var isCollect = false
	for _, value := range ScreenStrategys {
		ScreenStrategy := make(map[string]interface{})
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

func (a *ApiService) collect(c *gin.Context) {
	uid, _ := c.Get("Uid")
	// 根据uid查询用户信息
	uidFormatted := fmt.Sprintf("%s", uid)
	id, ok := c.GetQuery("id")
	if !ok {
		logrus.Error("id not exist.")
		res := util.ResponseMsg(-1, "fail", "id not exist.")
		c.SecureJSON(http.StatusOK, res)
		return
	}
	collect, ok := c.GetQuery("collect")
	if !ok {
		logrus.Error("collect not exist.")
		res := util.ResponseMsg(-1, "fail", "collect not exist.")
		c.SecureJSON(http.StatusOK, res)
		return
	}
	boolcollect, err := strconv.ParseBool(collect)
	if err != nil {
		logrus.Error(err)
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}
	if boolcollect {
		err = db.UpdateAddCollectProduct(a.dbEngine, id, uidFormatted)
		if err != nil {
			logrus.Info("update secret err:", err)
			res := util.ResponseMsg(-1, "fail", err)
			c.SecureJSON(http.StatusOK, res)
			return
		}
	} else {
		err, user := db.QuerySecret(a.dbEngine, uidFormatted)
		if err != nil {
			res := util.ResponseMsg(-1, "fail", err)
			c.SecureJSON(http.StatusOK, res)
			return
		}
		product := user.CollectStragetyList
		oldId := fmt.Sprintf(",%s", id)
		product = strings.Replace(product, oldId, "", -1)
		err = db.UpdateDelCollectProduct(a.dbEngine, product, uidFormatted)
		if err != nil {
			logrus.Info("update secret err:", err)
			res := util.ResponseMsg(-1, "fail", err)
			c.SecureJSON(http.StatusOK, res)
			return
		}
	}
	body := make(map[string]interface{})
	res := util.ResponseMsg(1, "success", body)
	c.SecureJSON(http.StatusOK, res)
	return
}

func (a *ApiService) productInfo(c *gin.Context) {
	id, ok := c.GetQuery("id")
	if !ok {
		logrus.Error("id not exist.")
		res := util.ResponseMsg(-1, "fail", "id not exist.")
		c.SecureJSON(http.StatusOK, res)
		return
	}
	strategyInfo, err := db.GetStrategy(a.dbEngine, id)
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}
	var CollectStragetyList []string
	session := sessions.Default(c)
	uid := session.Get("Uid")
	if uid != nil {
		// 登录状态
		uidFormatted := fmt.Sprintf("%s", uid)
		user, err := db.GetUser(a.dbEngine, uidFormatted)
		if err != nil {
			res := util.ResponseMsg(-1, "fail", err)
			c.SecureJSON(http.StatusOK, res)
			return
		}
		CollectStragetyList = strings.Split(user.CollectStragetyList[1:len(user.CollectStragetyList)-1], ",")
	}
	body := make(map[string]interface{})
	isCollect := isInCollectStrategyList(id, CollectStragetyList)
	body["id"] = strategyInfo.StrategyID
	body["name"] = strategyInfo.StrategyID
	body["recommendRate"] = strategyInfo.RecommendRate
	body["strategySource"] = strategyInfo.Source
	body["productCategory"] = strategyInfo.Type
	body["isCollect"] = isCollect
	body["collectionsNum"] = strategyInfo.ParticipateNum
	body["totalRevenue"] = strategyInfo.TotalRevenue
	body["totalYield"] = strategyInfo.TotalYield
	body["runTime"] = strategyInfo.CreateTime
	body["strategyDesc"] = strategyInfo.Describe
	body["expectedYield"] = strategyInfo.ExpectedBefenit
	body["winRate"] = strategyInfo.WinChance
	body["maxWithdrawalRate"] = strategyInfo.MaxDrawDown
	body["sharpeRatio"] = strategyInfo.SharpRatio
	body["controlLine"] = strategyInfo.ControlLine
	body["leverageRatio"] = strategyInfo.LeverageRatio
	body["minimumInvestmentAmount"] = strategyInfo.MinInvest
	body["policyCapacity"] = strategyInfo.Cap
	body["tradableAssets"] = strategyInfo.TradableAssets
	body["transactionCurrency"] = strategyInfo.CoinName
	body["shareRatio"] = strategyInfo.ShareRatio
	body["divideIntoPeriods"] = strategyInfo.DividePeriod
	body["protocolPeriod"] = strategyInfo.AgreementPeriod
	body["hostingPlatform"] = strategyInfo.HostPlatform
	res := util.ResponseMsg(1, "success", body)
	c.SecureJSON(http.StatusOK, res)
	return
}

func (a *ApiService) transactionRecords(c *gin.Context) {
	id, ok := c.GetQuery("id")
	if !ok {
		logrus.Error("id not exist.")
		res := util.ResponseMsg(-1, "fail", "id not exist.")
		c.SecureJSON(http.StatusOK, res)
		return
	}
	pageSize, ok := c.GetQuery("pageSize")
	if !ok {
		logrus.Error("pageSize not exist.")
		res := util.ResponseMsg(-1, "fail", "id not exist.")
		c.SecureJSON(http.StatusOK, res)
		return
	}
	pageIndex, ok := c.GetQuery("pageIndex")
	if !ok {
		logrus.Error("pageIndex not exist.")
		res := util.ResponseMsg(-1, "fail", "id not exist.")
		c.SecureJSON(http.StatusOK, res)
		return
	}

	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		logrus.Error(err)
		return
	}
	pageIndexInt, err := strconv.Atoi(pageIndex)
	if err != nil {
		logrus.Error(err)
		return
	}
	Records, err := db.TransactionRecords(a.dbEngine, pageSizeInt, pageIndexInt, id)
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}
	var RecordsList []interface{}
	for _, value := range Records {
		RecordsInfo := make(map[string]interface{})
		RecordsInfo["id"] = value.ID
		RecordsInfo["time"] = value.Time
		RecordsInfo["action"] = value.Action
		RecordsInfo["behavior"] = value.Behavior
		RecordsList = append(RecordsList, RecordsInfo)
	}
	body := make(map[string]interface{})
	body["total"] = len(Records)
	body["list"] = RecordsList
	res := util.ResponseMsg(1, "success", body)
	c.SecureJSON(http.StatusOK, res)
	return
}
