package api

import (
	"fmt"
	"github.com/ethereum/BGService/db"
	"github.com/ethereum/BGService/types"
	"github.com/ethereum/BGService/util"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
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

func (a *ApiService) investHandle(c *gin.Context, uidFormatted string, id string, ProductId string, Balance string) (error, *types.Strategy, float64, float64, float64) {
	userBindInfos, err := db.GetIdUserBindInfos(a.dbEngine, uidFormatted, id)
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return err, nil, 0, 0, 0
	}
	if userBindInfos == nil {
		res := util.ResponseMsg(-1, "fail", "apiKey is not exist")
		c.SecureJSON(http.StatusOK, res)
		return err, nil, 0, 0, 0
	}
	// 获取具体产品
	strategyInfo, err := db.GetStrategy(a.dbEngine, ProductId)
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return err, nil, 0, 0, 0
	}
	//principalGuaranteeDepositDrop, err := strconv.ParseFloat(strategyInfo.PrincipalGuaranteeDepositDrop, 64)
	//if err != nil {
	//	logrus.Error(err)
	//	return err, nil, 0, 0, 0
	//}
	//shareBonusDrop, err := strconv.ParseFloat(strategyInfo.ShareBonusDrop, 64)
	//if err != nil {
	//	logrus.Error(err)
	//	return err, nil, 0, 0, 0
	//}
	//managementFeesDrop, err := strconv.ParseFloat(strategyInfo.ManagementFeesDrop, 64)
	//if err != nil {
	//	logrus.Error(err)
	//	return err, nil, 0, 0, 0
	//}

	//shareBonus := Balance * 0 / 100
	//managementFees := Balance * 0 / 100
	//principalGuaranteeDeposit := Balance * 0 / 100
	//return nil, strategyInfo, shareBonus, managementFees, principalGuaranteeDeposit
	return nil, strategyInfo, 0, 0, 0
}

func (a *ApiService) invest(c *gin.Context) {
	uid, _ := c.Get("Uid")
	// 根据uid查询用户信息
	uidFormatted := fmt.Sprintf("%s", uid)
	// 交易所id
	id := "2"
	ProductId := "1"
	//id, ok := c.GetQuery("id")
	//if !ok {
	//	logrus.Error("id not exist.")
	//	res := util.ResponseMsg(-1, "fail", "id not exist.")
	//	c.SecureJSON(http.StatusOK, res)
	//	return
	//}
	//// 投入产品id
	//ProductId, ok := c.GetQuery("productId")
	//if !ok {
	//	logrus.Error("productId not exist.")
	//	res := util.ResponseMsg(-1, "fail", "productId not exist.")
	//	c.SecureJSON(http.StatusOK, res)
	//	return
	//}
	var balance string
	// 根据产品属性 取响应的 现货 U本位 币本位 获取余额
	strategy, err := db.GetProduct(a.dbEngine, ProductId)
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}
	//还要根据策略名字解析得到具体交易币种
	array := strings.Split(strategy.StrategyName, "/")

	switch strategy.CoinName {
	case "SPOT":
		//取现货余额
		userData, err := util.GetBinanceSpotUserData()
		if err != nil {

		}
		for {
			if err == nil {
				break
			}
			userData, err = util.GetBinanceSpotUserData()
		}
		for _, data := range userData {
			if strings.ToLower(data.Coin) == strings.ToLower(array[1]) {
				balance = data.Free
			}
		}
	case "CM":
		//取币本位余额
		userData, err := util.GetBinanceCMUserData()
		if err != nil {

		}
		for {
			if err == nil {
				break
			}
			userData, err = util.GetBinanceCMUserData()
		}

		for _, asset := range userData.Assets {
			if strings.ToLower(asset.Asset) == strings.ToLower(array[0]) {
				balance = asset.MarginBalance
			}
		}
	case "UM":
		//取U本位余额
		userData, err := util.GetBinanceUMUserData()
		if err != nil {

		}
		for {
			if err == nil {
				break
			}
			userData, err = util.GetBinanceUMUserData()
		}
		for _, asset := range userData.Assets {
			if strings.ToLower(asset.Asset) == strings.ToLower(array[1]) {
				balance = asset.MarginBalance
			}
		}
	}

	err, _, shareBonus, managementFees, principalGuaranteeDeposit := a.investHandle(c, uidFormatted, id, ProductId, balance)
	if err != nil {
		logrus.Error(err)
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}
	body := make(map[string]interface{})
	body["usableBalance"] = balance
	body["investBudget"] = balance
	body["shareBonusDrop"] = 0
	body["managementFeesDrop"] = 0
	body["principalGuaranteeDepositDrop"] = 0
	body["shareBonus"] = shareBonus
	body["managementFees"] = managementFees
	body["principalGuaranteeDeposit"] = principalGuaranteeDeposit

	res := util.ResponseMsg(1, "success", body)
	c.SecureJSON(http.StatusOK, res)
	return
}

func (a *ApiService) executeStrategy(c *gin.Context) {
	uid, _ := c.Get("Uid")
	// 根据uid查询用户信息
	uidFormatted := fmt.Sprintf("%s", uid)
	var payload *types.ExecuteStrategyInput
	if err := c.ShouldBindJSON(&payload); err != nil {
		logrus.Error(err)
		res := util.ResponseMsg(-1, "fail", err.Error())
		c.SecureJSON(http.StatusOK, res)
		return
	}
	// TODO 根据用户绑定的交易所获取余额
	balance := ""
	err, _, _, _, _ := a.investHandle(c, uidFormatted, payload.ID, payload.ProductId, balance)
	if err != nil {
		logrus.Error(err)
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}
	// todo 扣除费用
	actualInvest, err := decimal.NewFromString(balance)
	if err != nil {
		logrus.Error(err)
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}
	UserStrategy := types.UserStrategy{
		Uid:          uidFormatted,
		StrategyID:   payload.ProductId,
		JoinTime:     time.Now(), //.Format("2006-01-02"),
		ActualInvest: actualInvest.String(),
	}
	err = db.InsertUserStrategy(a.dbEngine, &UserStrategy)
	if err != nil {
		logrus.Error(err)
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}
	body := make(map[string]interface{})
	res := util.ResponseMsg(1, "success", body)
	c.SecureJSON(http.StatusOK, res)
	return
}

func ProductRevenue(a *ApiService, Revenue []map[string]string) (error, []map[string]interface{}) {
	// 总收益排行
	var ProductRevenueList []map[string]interface{}
	for i := 0; i < len(Revenue); i++ {
		UserRevenue := make(map[string]interface{})
		strategy, err := db.GetStrategy(a.dbEngine, Revenue[i]["f_stragetyID"])
		if err != nil {
			return err, ProductRevenueList
		}
		fRevenue, err := strconv.ParseFloat(Revenue[i]["totalBenefit"], 64)
		if err != nil {
			return err, ProductRevenueList
		}
		UserRevenue["stragetyName"] = strategy.StrategyName
		UserRevenue["revenue"] = fRevenue
		ProductRevenueList = append(ProductRevenueList, UserRevenue)
	}
	sort.Slice(ProductRevenueList, func(i, j int) bool {
		return getRevenue(ProductRevenueList[i], "revenue") > getRevenue(ProductRevenueList[j], "revenue")
	})
	for key, value := range ProductRevenueList {
		value["placed"] = key + 1
	}
	return nil, ProductRevenueList
}

func ProductRevenueRatio(a *ApiService, AllRevenue []map[string]string, AllInvest []map[string]string) (error, []map[string]interface{}) {
	var RevenueRatioRanking []map[string]interface{}
	for _, RevenueValue := range AllRevenue {
		for _, InvestValue := range AllInvest {
			if RevenueValue["f_stragetyID"] == InvestValue["f_strategyID"] {
				RevenueRatio := make(map[string]interface{})
				strategy, err := db.GetStrategy(a.dbEngine, RevenueValue["f_stragetyID"])
				if err != nil {
					return err, RevenueRatioRanking
				}
				fRevenue, err := strconv.ParseFloat(RevenueValue["totalBenefit"], 64)
				if err != nil {
					return err, RevenueRatioRanking
				}
				fInvest, err := strconv.ParseFloat(InvestValue["totalInvest"], 64)
				if err != nil {
					return err, RevenueRatioRanking
				}
				revenueRatio := fRevenue / fInvest
				RevenueRatio["revenueRatio"] = revenueRatio
				RevenueRatio["stragetyName"] = strategy.StrategyName
				RevenueRatioRanking = append(RevenueRatioRanking, RevenueRatio)
			}
		}
	}
	sort.Slice(RevenueRatioRanking, func(i, j int) bool {
		return getRevenue(RevenueRatioRanking[i], "revenueRatio") > getRevenue(RevenueRatioRanking[j], "revenueRatio")
	})
	for key, value := range RevenueRatioRanking {
		value["placed"] = key + 1
	}
	return nil, RevenueRatioRanking
}

func (a *ApiService) productRanking(c *gin.Context) {
	total := c.Query("total")
	totalInt, err := strconv.Atoi(total)
	if err != nil {
		logrus.Error(err)
	}
	err, Revenue := db.ProductRevenue(a.dbEngine, totalInt)
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}
	// 产品总收益排行
	err, ProductRevenueList := ProductRevenue(a, Revenue)
	if err != nil {
		return
	}

	// 总收益率排行
	// 总收益
	err, AllRevenue := db.ProductAllRevenue(a.dbEngine)
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}
	// 总投资
	err, AllInvest := db.ProductAllInvest(a.dbEngine)
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}
	err, ProductRevenueRatioList := ProductRevenueRatio(a, AllRevenue, AllInvest)
	if err != nil {
		return
	}
	var RevenueRatio []map[string]interface{}
	if len(ProductRevenueRatioList) > totalInt {
		RevenueRatio = ProductRevenueRatioList[:totalInt]
	} else {
		RevenueRatio = ProductRevenueRatioList
	}
	body := make(map[string]interface{})
	body["revenueAmount"] = ProductRevenueList
	body["revenueAmountRatio"] = RevenueRatio
	res := util.ResponseMsg(1, "success", body)
	c.SecureJSON(http.StatusOK, res)
	return
}
