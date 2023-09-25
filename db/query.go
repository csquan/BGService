package db

import (
	"fmt"
	"github.com/ethereum/BGService/types"
	"github.com/go-xorm/xorm"
	"github.com/sirupsen/logrus"
	"time"
)

func QuerySecret(engine *xorm.Engine, uid string) (error, *types.Users) {
	var user types.Users
	has, err := engine.Where("f_uid=?", uid).Get(&user)
	if err != nil {
		logrus.Error(err)
		return err, nil
	}
	if has {
		return nil, &user // 返回指向 user 的指针
	}
	return nil, nil
}

func QueryInvite(engine *xorm.Engine, uid string) (error, *types.Invitation) {
	var invitation types.Invitation
	has, err := engine.Where("`f_sonUid`=?", uid).Get(&invitation)
	if err != nil {
		logrus.Error(err)
		return err, nil
	}
	if has {
		return nil, &invitation // 返回指向 user 的指针
	}
	return nil, nil
}

func QueryEmail(engine *xorm.Engine, email string) (error, *types.Users) {
	var user types.Users
	has, err := engine.Where("`f_mailBox`=?", email).Get(&user)
	if err != nil {
		logrus.Error(err)
		return err, nil
	}
	if has {
		return nil, &user // 返回指向 user 的指针
	}
	return nil, nil
}

func UserRevenue(engine *xorm.Engine, total int) (error, []map[string]string) {
	sql := fmt.Sprintf("SELECT f_uid, SUM(`f_totalBenefit`) AS `totalBenefit` "+
		"FROM `userStrategyEarnings` GROUP BY f_uid ORDER BY `totalBenefit` DESC limit %d", total)
	result, err := engine.QueryString(sql)
	if err != nil {
		logrus.Error(err)
		return err, nil
	}
	return nil, result
}

func ProductRevenue(engine *xorm.Engine, total int) (error, []map[string]string) {
	sql := fmt.Sprintf("SELECT `f_stragetyID`, SUM(`f_totalBenefit`) AS `totalBenefit` "+
		"FROM `userStrategyEarnings` GROUP BY `f_stragetyID` ORDER BY `totalBenefit` DESC limit %d", total)
	result, err := engine.QueryString(sql)
	if err != nil {
		logrus.Error(err)
		return err, nil
	}
	return nil, result
}

func UserAllRevenue(engine *xorm.Engine) (error, []map[string]string) {
	sql := fmt.Sprintf("SELECT f_uid, SUM(`f_totalBenefit`) AS `totalBenefit` " +
		"FROM `userStrategyEarnings` GROUP BY f_uid ORDER BY `totalBenefit` DESC")
	result, err := engine.QueryString(sql)
	if err != nil {
		logrus.Error(err)
		return err, nil
	}
	return nil, result
}

func UserAllInvest(engine *xorm.Engine) (error, []map[string]string) {
	sql := fmt.Sprintf("select f_uid, SUM(`f_actualInvest`) AS `totalInvest` " +
		"from `userStrategy` GROUP BY f_uid ORDER BY `totalInvest` DESC")
	result, err := engine.QueryString(sql)
	if err != nil {
		logrus.Error(err)
		return err, nil
	}
	return nil, result
}

func ProductAllRevenue(engine *xorm.Engine) (error, []map[string]string) {
	sql := fmt.Sprintf("SELECT `f_stragetyID`, SUM(`f_totalBenefit`) AS `totalBenefit` " +
		"FROM `userStrategyEarnings` GROUP BY `f_stragetyID` ORDER BY `totalBenefit` DESC")
	result, err := engine.QueryString(sql)
	if err != nil {
		logrus.Error(err)
		return err, nil
	}
	return nil, result
}

func ProductAllInvest(engine *xorm.Engine) (error, []map[string]string) {
	sql := fmt.Sprintf("select `f_strategyID`, SUM(`f_actualInvest`) AS `totalInvest` " +
		"from `userStrategy` GROUP BY `f_strategyID` ORDER BY `totalInvest` DESC")
	result, err := engine.QueryString(sql)
	if err != nil {
		logrus.Error(err)
		return err, nil
	}
	return nil, result
}

func QueryInviteCode(engine *xorm.Engine, InviteCode string) (error, *types.Users) {
	var user types.Users
	has, err := engine.Table("users").Where("`f_invitationCode`=?", InviteCode).Get(&user)
	if err != nil {
		logrus.Error(err)
		return err, nil
	}
	if has {
		return nil, &user // 返回指向 user 的指针
	}
	return nil, nil
}

func QueryInviteNum(engine *xorm.Engine, InviteCode string) (error, []types.Invitation) {
	var invitation []types.Invitation
	err := engine.Table("invitation").Where("`f_uid`=?", InviteCode).Find(&invitation)
	if err != nil {
		logrus.Error(err)
		return err, nil
	}
	return nil, invitation
}

func QueryInviteNumLimit(engine *xorm.Engine, InviteCode string, total int) (error, []types.Users) {
	var users []types.Users
	err := engine.Table("users").Where("`f_invitatedCode`=?", InviteCode).Limit(total).Find(&users)
	if err != nil {
		logrus.Error(err)
		return err, nil
	}
	return nil, users
}

func QueryClaimRewardNumber(engine *xorm.Engine) (error, []map[string]string) {
	sql := fmt.Sprintf(`select f_uid, count("f_uid") as "countUid" from invitation GROUP BY "f_uid" ORDER BY "countUid" DESC`)
	result, err := engine.QueryString(sql)
	if err != nil {
		logrus.Error(err)
		return err, nil
	}
	return nil, result
}

func GetUser(engine *xorm.Engine, uid string) (*types.Users, error) {
	var user types.Users
	has, err := engine.Where("f_uid=?", uid).Get(&user)
	if err != nil {
		return nil, err
	}
	if has {
		return &user, nil
	}
	return nil, nil
}

func GetProduct(engine *xorm.Engine, productID string) (*types.Strategy, error) {
	var strategy types.Strategy
	has, err := engine.Where("`f_strategyID`=?", productID).Get(&strategy)
	if err != nil {
		return nil, err
	}
	if has {
		return &strategy, nil
	}
	return nil, nil
}

func GetUserAsset(engine *xorm.Engine, uid string) (*types.UserAsset, error) {
	var userAsset types.UserAsset
	coinName := "trx"
	has, err := engine.Table("userAsset").Where("f_uid=? and `f_coinName`=?", uid, coinName).Get(&userAsset)
	if err != nil {
		return nil, err
	}
	if has {
		return &userAsset, nil
	}
	return nil, nil
}

func GetUserAllAsset(engine *xorm.Engine, uid string) (error, []types.UserAsset) {
	var userAssets []types.UserAsset
	err := engine.Table("userAsset").Where("`f_uid`=?", uid).Find(&userAssets)
	if err != nil {
		logrus.Error(err)
		return err, nil
	}
	return nil, userAssets
}

func GetUserAddrs(engine *xorm.Engine, uid string) (error, []types.UserAddr) {
	var userAddrs []types.UserAddr
	err := engine.Table("userAddr").Where("f_uid=?", uid).Find(&userAddrs)
	if err != nil {
		logrus.Error(err)
		return err, nil
	}
	return nil, userAddrs
}

func GetUserAddr(engine *xorm.Engine, uid string, network string) (*types.UserAddr, error) {
	var userAddr types.UserAddr
	has, err := engine.Table("userAddr").Where("f_uid=? and f_network=?", uid, network).Get(&userAddr)
	if err != nil {
		return nil, err
	}
	if has {
		return &userAddr, nil
	}
	return nil, nil
}

func GetUserFundIn(engine *xorm.Engine, uid string, network string) (*types.UserFundIn, error) {
	var userFundIn types.UserFundIn
	has, err := engine.Table("fundIn").Where("f_uid=? and f_network=?", uid, network).OrderBy("f_id desc").Limit(1).Get(&userFundIn)
	if err != nil {
		return nil, err
	}
	if has {
		return &userFundIn, nil
	}
	return nil, nil
}

func GetUserAllFundIn(engine *xorm.Engine, uid string) (*[]types.UserFundIn, error) {
	var userFundIns []types.UserFundIn
	err := engine.Table("fundIn").Where("`f_uid`=?", uid).Desc("f_id").Find(&userFundIns)
	if err != nil {
		return nil, err
	}
	return &userFundIns, nil
}

func GetUserAllFundOut(engine *xorm.Engine, userAddr string) (*[]types.UserFundOut, error) {
	var userFundOuts []types.UserFundOut
	err := engine.Table("fundOut").Where("f_from= ?", userAddr).Desc("f_id").Find(&userFundOuts)
	if err != nil {
		return nil, err
	}
	return &userFundOuts, nil
}

func GetUserAllShare(engine *xorm.Engine, uid string) (*[]types.UserShare, error) {
	var userShares []types.UserShare
	err := engine.Table("shareRecords").Where("`f_uid`=?", uid).Desc("f_id").Find(&userShares)
	if err != nil {
		return nil, err
	}
	return &userShares, nil
}

func GetUserAllExperience(engine *xorm.Engine, uid string) (*[]types.UserExperience, error) {
	var userExperiences []types.UserExperience
	err := engine.Table("userExperience").Where("`f_uid`=?", uid).Desc("f_id").Find(&userExperiences)
	if err != nil {
		return nil, err
	}
	return &userExperiences, nil
}

func GetUserKey(engine *xorm.Engine, addr string) (*types.UserKey, error) {
	var userKey types.UserKey
	has, err := engine.Table("userKey").Where("f_addr=?", addr).Get(&userKey)
	if err != nil {
		return nil, err
	}
	if has {
		return &userKey, nil
	}
	return nil, nil
}

func GetUserBindInfos(engine *xorm.Engine, uid string) (*types.UserBindInfos, error) {
	var userBindInfos types.UserBindInfos
	has, err := engine.Table("userBindInfos").Where("f_uid=?", uid).Get(&userBindInfos)
	if err != nil {
		return nil, err
	}
	if has {
		return &userBindInfos, nil
	}
	return nil, nil
}

func GetUserBindInfoByUidCex(engine *xorm.Engine, uid string, cex string) (*types.UserBindInfos, error) {
	var userBindInfos types.UserBindInfos
	has, err := engine.Table("userBindInfos").Where("f_uid=? and f_cex=?", uid, cex).Get(&userBindInfos)
	if err != nil {
		return nil, err
	}
	if has {
		return &userBindInfos, nil
	}
	return nil, nil
}

func GetApiKeyUserBindInfos(engine *xorm.Engine, apiKey string) (*types.UserBindInfos, error) {
	var userBindInfos types.UserBindInfos
	has, err := engine.Table("userBindInfos").Where("`f_apiKey`=?", apiKey).Get(&userBindInfos)
	if err != nil {
		return nil, err
	}
	if has {
		return &userBindInfos, nil
	}
	return nil, nil
}

func GetIdUserBindInfos(engine *xorm.Engine, uid string, apiId string) (*types.UserBindInfos, error) {
	var userBindInfos types.UserBindInfos
	has, err := engine.Table("userBindInfos").
		Where("`f_uid`=? and `f_id`=?", uid, apiId).Get(&userBindInfos)
	if err != nil {
		return nil, err
	}
	if has {
		return &userBindInfos, nil
	}
	return nil, nil
}

func GetAllUserBindInfos(engine *xorm.Engine, uid string) ([]types.UserBindInfos, error) {
	var userBindInfos []types.UserBindInfos
	err := engine.Table("userBindInfos").Where("f_uid=?", uid).Find(&userBindInfos)
	if err != nil {
		return nil, err
	}
	return userBindInfos, nil
}

func GetUserExperience(engine *xorm.Engine, uid string) (*types.UserExperience, error) {
	var userExperience types.UserExperience
	has, err := engine.Table("userExperience").Where("f_uid=?", uid).Get(&userExperience)
	if err != nil {
		return nil, err
	}
	if has {
		return &userExperience, nil
	}
	return nil, nil
}

func GetPlatformExperience(engine *xorm.Engine) (*types.PlatformExperience, error) {
	var platformExperience types.PlatformExperience
	has, err := engine.Table("platformExperience").Get(&platformExperience)
	if err != nil {
		return nil, err
	}
	if has {
		return &platformExperience, nil
	}
	return nil, nil
}

func GetMsg(engine *xorm.Engine, pageSizeInt int, pageIndexInt int, Type string) ([]types.News, error) {
	var news []types.News
	err := engine.Table("news").Where("f_type=?", Type).Limit(pageSizeInt, (pageIndexInt-1)*pageSizeInt).Find(&news)
	if err != nil {
		return nil, err
	}
	return news, nil
}

func GetHotMsg(engine *xorm.Engine, total int, Type string) ([]types.News, error) {
	var news []types.News
	err := engine.Table("news").Where("f_type=? and f_hotspot='t'", Type).Limit(total).Find(&news)
	if err != nil {
		return nil, err
	}
	return news, nil
}

func GetMsgDetail(engine *xorm.Engine, msgId string) (*types.News, error) {
	var news types.News
	has, err := engine.Table("news").Where("f_id=?", msgId).Get(&news)
	if err != nil {
		return nil, err
	}
	if has {
		return &news, nil
	}
	return nil, nil
}

func GetConcernList(engine *xorm.Engine, uid string) (tags []string) {
	// the select query, returning 1 column of array type
	url := "SELECT 'f_concernCoinList' FROM users WHERE f_uid=$1"

	ret, err := engine.Query(url, uid)
	// wrap the output parameter in pq.Array for receiving into it
	//has, err := engine.SQL(url, uid).Get(pq.Array(&foo.concernCoinList))

	if err != nil {
		logrus.Info(ret)
		logrus.Info(err)
	}
	logrus.Info(ret)
	return
}

func GetStrategy(engine *xorm.Engine, sid string) (*types.Strategy, error) {
	var strategy types.Strategy
	has, err := engine.Table("strategys").Where("`f_strategyID`=? and `f_isValid`=?", sid, true).Get(&strategy)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	if has {
		return &strategy, nil
	}
	return nil, nil
}

func GetStrategyByName(engine *xorm.Engine, sName string) (*types.Strategy, error) {
	var strategy types.Strategy
	has, err := engine.Table("strategy").Where("`f_strategyName`=?", sName).Get(&strategy)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	if has {
		return &strategy, nil
	}
	return nil, nil
}

func GetExactlyUserStrategy(engine *xorm.Engine, uid string, sid string) (*types.UserStrategy, error) {
	var userStrategy types.UserStrategy
	has, err := engine.Table("userStrategy").Where("f_uid=? and `f_strategyID`=?", uid, sid).Get(&userStrategy)
	if err != nil {
		return nil, err
	}
	if has {
		return &userStrategy, nil
	}
	return nil, nil
}

// 查询天级产品总收益
func GetAllStrategyBenefits(engine *xorm.Engine, sid string, startTime string, endTime string) ([]map[string]string, error) {
	sql := fmt.Sprintf(`select to_char("f_createTime"::DATE, 'YYYY-MM-DD') as day, sum("f_totalBenefit")  as "f_totalBenefit"
from "userStrategyEarnings" where "f_stragetyID"='%s' and "f_createTime">= '%s' and "f_createTime"<= '%s'  GROUP BY "day" ORDER BY "f_totalBenefit" DESC`, sid, startTime, endTime)
	result, err := engine.QueryString(sql)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return result, nil
}

// 取出天极产品的实际投资额
func GetExactlyStrategy(engine *xorm.Engine, sid string, startTime string, endTime string) ([]map[string]string, error) {
	sql := fmt.Sprintf(`select to_char("f_joinTime"::DATE, 'YYYY-MM-DD') as day, sum("f_actualInvest")  as "f_actualInvest"
from "userStrategy" where "f_strategyID"='%s'  and "f_joinTime">= '%s' and "f_joinTime"<= '%s' GROUP BY "day"  ORDER BY "f_actualInvest" DESC`, sid, startTime, endTime)
	result, err := engine.QueryString(sql)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return result, nil
}

func GetUserStrategys(engine *xorm.Engine, uid string) ([]*types.UserStrategy, error) {
	var userStrategys []*types.UserStrategy
	err := engine.Table("userStrategy").Where("f_uid=?", uid).Find(&userStrategys)
	if err != nil {
		return nil, err
	}
	return userStrategys, nil
}

func GetUserStrategy(engine *xorm.Engine, uid string) ([]*types.UserStrategy, error) {
	var userStrategys []*types.UserStrategy
	err := engine.Table("userStrategy").Where("f_uid=?", uid).Find(&userStrategys)
	if err != nil {
		return nil, err
	}
	return userStrategys, nil
}

func GetStrategyTotalAssets(engine *xorm.Engine) (float64, error) {
	var userStrategy types.UserStrategy
	total, err := engine.Table("`userStrategy`").Sum(userStrategy, "`f_actualInvest`")
	if err != nil {
		return 0, err
	}
	return total, nil
}
func GetStrategyTotalBenefits(engine *xorm.Engine, sid string) (float64, error) {
	var userStrategyEarnings types.UserStrategyEarnings
	total, err := engine.Table("`userStrategyEarnings`").Where("`f_strategyID`=?", sid).Sum(userStrategyEarnings, "`f_totalBenefit`")
	if err != nil {
		return 0, err
	}
	return total, nil
}

func GetStrategyBenefits(engine *xorm.Engine, sid, uid string, startTime string, endTime string) ([]*types.UserStrategyEarnings, error) {
	var userStrategyEarnings []*types.UserStrategyEarnings
	err := engine.Table("userStrategyEarnings").Where("f_uid = ? and `f_strategyID` = ? and `f_createTime`>= ? and `f_createTime`<= ?", uid, sid, startTime, endTime).Find(&userStrategyEarnings)
	if err != nil {
		return nil, err
	}
	return userStrategyEarnings, nil
}

func GetStrategyTotalInvests(engine *xorm.Engine, sid string) (float64, error) {
	var userStrategy types.UserStrategy
	total, err := engine.Table("`userStrategy`").Where("`f_strategyID`=?", sid).Sum(userStrategy, "`f_actualInvest`")
	if err != nil {
		return 0, err
	}
	return total, nil
}

func GetStrategyUserCount(engine *xorm.Engine) (int64, error) {
	var userStrategy types.UserStrategy
	total, err := engine.Table("`userStrategy`").Distinct("f_uid").Count(userStrategy)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func GetAllStrategy(engine *xorm.Engine) ([]*types.Strategy, error) {
	var Strategy []*types.Strategy
	err := engine.Table("strategys").Where("`f_isValid`=?", true).Find(&Strategy)
	if err != nil {
		return nil, err
	}
	return Strategy, nil
}

func GetUserStrategyLatestEarnings(engine *xorm.Engine, uid string, sid string) (*types.UserStrategyEarnings, error) {
	var userStrategyEarnings types.UserStrategyEarnings
	has, err := engine.Table("`userStrategyEarnings`").Where("f_uid=? and `f_strategyID`=?", uid, sid).Desc("f_id").Limit(1).Get(&userStrategyEarnings)
	if err != nil {
		return nil, err
	}
	if has {
		return &userStrategyEarnings, nil
	}
	return nil, nil
}

func GetUserIncome(engine *xorm.Engine) (float64, error) {
	var userStrategyEarnings types.UserStrategyEarnings
	total, err := engine.Table("`userStrategyEarnings`").Sum(userStrategyEarnings, "`f_totalBenefit`")
	if err != nil {
		return 0, err
	}
	return total, nil
}

func timeFmt(timeCycle int) (string, string) {
	// 1:0~6个月  2:6~12个月 3:12~36个月 4:36个月以上
	var startTime string
	var endTime string
	timeNow := time.Now()
	sixMonthsAgo := timeNow.AddDate(0, -6, 0).Format("2006-01-02")
	twelveMonthsAgo := timeNow.AddDate(-1, 0, 0).Format("2006-01-02")
	thirtySixMonthsAgo := timeNow.AddDate(-3, 0, 0).Format("2006-01-02")
	if timeCycle == 1 {
		startTime = sixMonthsAgo
		endTime = timeNow.Format("2006-01-02")
	} else if timeCycle == 2 {
		startTime = twelveMonthsAgo
		endTime = sixMonthsAgo
	} else if timeCycle == 3 {
		startTime = thirtySixMonthsAgo
		endTime = twelveMonthsAgo
	} else if timeCycle == 4 {
		startTime = "2006-01-02"
		endTime = thirtySixMonthsAgo
	} else {
		startTime = "2006-01-02"
		endTime = timeNow.Format("2006-01-02")
	}
	return startTime, endTime
}

func ExpectedYieldFmt(ExpectedYield int) (string, string) {
	// '预期收益率' -1全部 1:0~50%  2:50%~100% 3:100%~300%
	var startExpected string
	var endExpected string
	if ExpectedYield == 1 {
		startExpected = "0"
		endExpected = "50"
	} else if ExpectedYield == 2 {
		startExpected = "50"
		endExpected = "100"
	} else if ExpectedYield == 3 {
		startExpected = "100"
		endExpected = "300"
	} else {
		startExpected = "0"
		endExpected = "300"
	}
	return startExpected, endExpected
}

func WithdrawalRateFmt(WithdrawalRate int) (string, string) {
	// '最大回撤率' -1全部 1:0~20%  2:20%~40% 3:40%~60%
	var startWithdrawalRate string
	var endWithdrawalRate string
	if WithdrawalRate == 1 {
		startWithdrawalRate = "0"
		endWithdrawalRate = "20"
	} else if WithdrawalRate == 2 {
		startWithdrawalRate = "20"
		endWithdrawalRate = "40"
	} else if WithdrawalRate == 3 {
		startWithdrawalRate = "40"
		endWithdrawalRate = "60"
	} else {
		startWithdrawalRate = "0"
		endWithdrawalRate = "60"
	}
	return startWithdrawalRate, endWithdrawalRate
}

func GetScreenStrategy(engine *xorm.Engine, payload *types.StrategyInput, CollectStragety []int) ([]types.Strategy, error) {
	var strategy []types.Strategy
	sortMap := map[int]string{
		1: "f_totalYield",
		2: "f_totalRevenue",
		3: "f_participateNum",
		4: "f_maxDrawDown",
	}
	sessionSql := engine.Table("strategys").Where("`f_isValid`=?", true)
	// 我的收藏
	if payload.Strategy == 1 {
		sessionSql = sessionSql.In("`f_strategyID`", CollectStragety)
	}
	// 币种
	if payload.Currency != "" && payload.Currency != "-1" {
		sessionSql = sessionSql.Where("`f_coinName` = ?", payload.Currency)
	}
	// 产品来源
	if payload.StrategySource != 0 && payload.StrategySource != -1 {
		sessionSql = sessionSql.Where("`f_source` = ?", payload.StrategySource)
	}
	// 产品类别
	if payload.ProductCategory != 0 && payload.ProductCategory != -1 {
		sessionSql = sessionSql.Where("`f_type` = ?", payload.ProductCategory)
	}
	// 时间
	if payload.RunTime != 0 && payload.RunTime != -1 {
		startTime, endTime := timeFmt(payload.RunTime)
		sessionSql = sessionSql.Where("?<=`f_createTime`<?", startTime, endTime)
	}
	// 预期收益率
	if payload.ExpectedYield != 0 && payload.ExpectedYield != -1 {
		startExpected, endExpected := ExpectedYieldFmt(payload.ExpectedYield)
		sessionSql = sessionSql.Where("?<=`f_expectedBefenit`<?", startExpected, endExpected)
	}
	// 最大回撤率
	if payload.MaxWithdrawalRate != 0 && payload.MaxWithdrawalRate != -1 {
		startExpected, endExpected := WithdrawalRateFmt(payload.ExpectedYield)
		sessionSql = sessionSql.Where("?<=`f_maxDrawDown`<?", startExpected, endExpected)
	}
	// 排序
	if payload.ComprehensiveSorting != 0 && payload.ComprehensiveSorting != -1 {
		filed := sortMap[payload.ComprehensiveSorting]
		sessionSql = sessionSql.Desc("`?`", filed)
	} else {
		sessionSql = sessionSql.Desc("`f_participateNum`").Desc("`f_createTime`").Desc("`f_recommendRate`")
	}
	err := sessionSql.Limit(payload.PageSize, (payload.PageIndex-1)*payload.PageIndex).Find(&strategy)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return strategy, nil
}

func GetSearchScreenStrategy(engine *xorm.Engine, payload *types.StrategyInput, CollectStragety []int) ([]types.Strategy, error) {
	var strategy []types.Strategy
	sortMap := map[int]string{
		1: "f_totalYield",
		2: "f_totalRevenue",
		3: "f_participateNum",
		4: "f_maxDrawDown",
	}
	sessionSql := engine.Table("strategys").Where("`f_isValid`=?", true)
	// 我的收藏
	if payload.Strategy == 1 {
		sessionSql = sessionSql.In("`f_strategyID`", CollectStragety)
	}
	sessionSql.Where(fmt.Sprintf("`f_strategyName` like '%s' ", "%"+payload.Keywords+"%"))
	// 排序
	if payload.ComprehensiveSorting != 0 && payload.ComprehensiveSorting != -1 {
		filed := sortMap[payload.ComprehensiveSorting]
		sessionSql = sessionSql.Desc("`?`", filed)
	} else {
		sessionSql = sessionSql.Desc("`f_participateNum`").Desc("`f_createTime`").Desc("`f_recommendRate`")
	}
	err := sessionSql.Limit(payload.PageSize, (payload.PageIndex-1)*payload.PageSize).Find(&strategy)
	if err != nil {
		return nil, err
	}
	return strategy, nil
}

func TransactionRecords(engine *xorm.Engine, pageSizeInt int, pageIndexInt int, id string) ([]types.TransactionRecords, error) {
	var transactionRecords []types.TransactionRecords
	err := engine.Table("`transactionRecords`").Where("`f_strategyID`=?", id).
		Limit(pageSizeInt, (pageIndexInt-1)*pageSizeInt).Find(&transactionRecords)
	if err != nil {
		return nil, err
	}
	return transactionRecords, nil
}
