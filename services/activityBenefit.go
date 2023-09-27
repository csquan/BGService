package services

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/delivery"
	"github.com/adshao/go-binance/v2/futures"
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
	totalBenefitFloat := float64(0)
	if totalBenefit != "" {
		totalBenefitFloat, err = strconv.ParseFloat(totalBenefit, 64)
		if err != nil {
			logrus.Error(err)
		}
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
	logrus.Info("***************************开始每日任务：系统200万投资收益统计***************************")
	// Create DB pool
	db, err := sql.Open("postgres", activityDbDSN)
	if err != nil {
		log.Fatal("Failed to open a DB connection: ", err)
	}
	defer db.Close()

	err1 := errors.New("init error")
	var userData *futures.Account
	for {
		userData, err1 = utils.GetBinanceUMUserData(types.ApiKeySystem, types.ApiSecretSystem)
		if err1 != nil {
			logrus.Info(err1)
		} else {
			logrus.Info("成功请求到币安U本位接口")
			break
		}
	}
	umSum := float64(0)
	for _, asset := range userData.Assets {
		MarginBalanceFloat, err := strconv.ParseFloat(asset.MarginBalance, 64)
		if err != nil {
			logrus.Error(err)
		}
		if MarginBalanceFloat > 0 {
			price := float64(0)
			if asset.Asset != "USDT" {
				symbols := make([]string, 1)
				symbols[0] = asset.Asset + "USDT"

				prices, err := utils.GetBinancePrice(types.ApiKeySystem, types.ApiSecretSystem, symbols)
				if err != nil {
					logrus.Error(err)
				}
				if err != nil {
					logrus.Error(prices)
				}
				price, err = strconv.ParseFloat(prices[0].Price, 64)
				if err != nil {
					logrus.Error(err)
				}
			} else {
				price = 1
			}

			assetSum := MarginBalanceFloat * float64(price)
			umSum = MarginBalanceFloat + float64(umSum)

			logrus.Info("取出对应资产：", asset.Asset, "价格为：", price)
			logrus.Info("该资产价值：", assetSum)
			logrus.Info("经过计算得到U本位累加资产", umSum)
		}
	}
	logrus.Info("U本位余额为", umSum)

	spotSum := float64(0)
	err2 := errors.New("init error")
	var userData2 *binance.Account
	for {
		userData2, err2 = utils.GetBinanceSpotUserData(types.ApiKeySystem, types.ApiSecretSystem)
		if err2 != nil {
			logrus.Info(err2)
		} else {
			logrus.Info("成功请求到币安现货接口")
			break
		}
	}
	for _, balance := range userData2.Balances {
		MarginBalanceFloat, err := strconv.ParseFloat(balance.Locked, 64)
		if err != nil {
			logrus.Error(err)
		}
		if MarginBalanceFloat > 0 {
			price := float64(0)
			if balance.Asset != "USDT" {
				symbols := make([]string, 1)
				symbols[0] = balance.Asset + "USDT"

				prices, err := utils.GetBinancePrice(types.ApiKeySystem, types.ApiSecretSystem, symbols)
				if err != nil {
					logrus.Error(err)
				}
				if err != nil {
					logrus.Error(prices)
				}
				price, err = strconv.ParseFloat(prices[0].Price, 64)
				if err != nil {
					logrus.Error(err)
				}
			} else {
				price = 1
			}

			assetSum := MarginBalanceFloat * float64(price)
			spotSum = spotSum + assetSum

			logrus.Info("取出对应资产：", balance.Asset, "价格为：", price)
			logrus.Info("该资产价值：", assetSum)
			logrus.Info("经过计算得到现货累加资产", spotSum)
		}
	}
	logrus.Info("现货余额为", spotSum)

	cmSum := float64(0)
	err3 := errors.New("init error")
	var userData3 *delivery.Account
	for {
		userData3, err3 = utils.GetBinanceCMUserData(types.ApiKeySystem, types.ApiSecretSystem)
		if err3 != nil {
			logrus.Info(err3)
		} else {
			logrus.Info("成功请求到币安币本位接口")
			break
		}
	}
	for _, asset := range userData3.Assets {
		MarginBalanceFloat, err := strconv.ParseFloat(asset.MarginBalance, 64)
		if err != nil {
			logrus.Error(err)
		}
		if MarginBalanceFloat > 0 {
			price := float64(0)
			if asset.Asset != "USDT" {
				symbols := make([]string, 1)
				symbols[0] = asset.Asset + "USDT"

				prices, err := utils.GetBinancePrice(types.ApiKeySystem, types.ApiSecretSystem, symbols)
				if err != nil {
					logrus.Error(err)
				}
				if err != nil {
					logrus.Error(prices)
				}
				price, err = strconv.ParseFloat(prices[0].Price, 64)
				if err != nil {
					logrus.Error(err)
				}
			} else {
				price = 1
			}

			assetSum := MarginBalanceFloat * float64(price)
			cmSum = cmSum + assetSum

			logrus.Info("取出对应资产：", asset.Asset, "价格为：", price)
			logrus.Info("该资产价值：", assetSum)
			logrus.Info("经过计算得到币本位累加资产", cmSum)
		}
	}
	logrus.Info("币本位余额为", cmSum)

	MarginBalanceFloat := cmSum + cmSum + spotSum

	logrus.Info("账户总余额为", MarginBalanceFloat)
	// 获取当前时间
	now := time.Now()
	// 计算昨天的时间
	yesterday := now.AddDate(0, 0, -1).Format("2006-01-02")

	// 累计收入
	totalBenefit := MarginBalanceFloat - 2000000
	// 查库中的累计收入
	totalBenefitFloat := queryActivityStrategyEarnings(db, "1", yesterday)
	// 今日收益
	totalDay := totalBenefit - totalBenefitFloat
	logrus.Info("今日收益", totalDay)
	insertActivityEarning(db, totalDay, totalBenefit, "1")
	return nil
}
