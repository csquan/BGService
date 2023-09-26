package services

import (
	"database/sql"
	"fmt"
	"github.com/ethereum/BGService/types"
	utils "github.com/ethereum/BGService/util"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"log"
	"strconv"
	"time"
	//_ "github.com/bmizerany/pq"
)

// 体验金总收益

type ActivityBenefitService struct {
}

func NewActivityBenefitService() *ActivityBenefitService {
	return &ActivityBenefitService{}
}
func (c *ActivityBenefitService) Name() string {
	return "UserBenefit"
}

const (
	activityDbDSN     = "postgres://postgres:12345@127.0.0.1:5432/postgres?sslmode=disable"
	activityApiKey    = "Xq2vyva4DUxw1EqywIHHZa8RDFIitXraDexa1LVONe3reuPNUEFuDYDs7JYjMY86"
	activityApiSecret = "reLDM7CYMHVPlw6FodmQvYpU9zRdndQ5NUlRFswKT6leKzcKl2BeP3tycqEaLBRZ"
	Asset             = "USDT"
)

var (
	//activity_future_testnet_binance_url = "https://testnet.binancefuture.com"
	activity_future_binance_url = "https://api.binance.com/api"
)

func queryActivityStrategyEarnings(db *sql.DB, strategyID string, createTime string) float64 {
	Sql := fmt.Sprintf(`SELECT "f_totalBenefit" FROM "platformExperienceEarnings" WHERE  "f_strategyID"='%s' and "f_createTime"='%s'`, strategyID, createTime)
	rows, err := db.Query(Sql)
	if err != nil {
		log.Fatal("Failed to execute query: ", err)
	}
	var totalBenefit string
	for rows.Next() {
		err = rows.Scan(&totalBenefit)
	}
	totalBenefitFloat, err := strconv.ParseFloat(totalBenefit, 64)
	if err != nil {
		logrus.Error(err)
	}
	return totalBenefitFloat
}

func insertActivityEarning(db *sql.DB, dayBenefit float64, totalBenefit float64, strategyID string) {
	insertSQL := `
		INSERT INTO "platformExperienceEarnings" ("f_strategyID", "f_dayBenefit", "f_totalBenefit")
		VALUES ($1, $2, $3)
	`
	// 要插入的数据
	data := []interface{}{strategyID, dayBenefit, totalBenefit}
	// 执行插入操作
	result, err := db.Exec(insertSQL, data...)
	if err != nil {
		log.Fatal(err)
	}

	// 获取受影响的行数
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("rowsAffected = %d", rowsAffected)
}

func (c *ActivityBenefitService) Run() error {
	// Create DB pool
	db, err := sql.Open("postgres", activityDbDSN)
	if err != nil {
		log.Fatal("Failed to open a DB connection: ", err)
	}
	defer db.Close()
	userData, err := utils.GetBinanceUMUserData(types.ApiKeySystem, types.ApiSecretSystem)
	if err != nil {
		logrus.Info(err)
	}
	var MarginBalance string // 当天余额
	for _, asset := range userData.Assets {
		if asset.Asset == Asset {
			fmt.Println("MarginBalance", asset.MarginBalance)
			MarginBalance = asset.MarginBalance
		}
	}
	MarginBalanceFloat, err := strconv.ParseFloat(MarginBalance, 64)
	if err != nil {
		logrus.Error(err)
	}
	// 获取当前时间
	now := time.Now()
	// 计算昨天的时间
	yesterday := now.AddDate(0, 0, -1).Format("2006-01-02")
	//todo:这里得本金是从平台体验金信息表中取

	// 累计收入
	totalBenefit := MarginBalanceFloat - 2000000
	// 查库中的累计收入
	totalBenefitFloat := queryActivityStrategyEarnings(db, "1", yesterday)
	// 今日收益
	totalDay := totalBenefit - totalBenefitFloat
	insertActivityEarning(db, totalDay, totalBenefit, "1")
	return nil
}
