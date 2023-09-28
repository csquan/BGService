package services

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/LinkinStars/go-scaffold/contrib/cryptor"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/ethereum/BGService/types"
	utils "github.com/ethereum/BGService/util"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
	//_ "github.com/bmizerany/pq"
)

type UserBenefitService struct {
}

func NewUserBenefitService() *UserBenefitService {
	return &UserBenefitService{}
}
func (c *UserBenefitService) Name() string {
	return "UserBenefit"
}

const (
	DB_DSN = "postgres://postgres:12345@127.0.0.1:5432/postgres?sslmode=disable"
)

var (
	base_future_testnet_binance_url = "https://testnet.binancefuture.com"
	//base_future_testnet_binance_url = "https://api.binance.com/api"
)

type UserStrategy struct {
	Strategyid   string
	Uid          string
	joinTime     time.Time
	actualInvest float64
	apiId        string
}

func queryUserStrategy(db *sql.DB) []UserStrategy {
	UserStrategySql := `SELECT "f_uid", "f_joinTime", "f_strategyID", "f_actualInvest", "f_apiId" FROM "userStrategy" WHERE "f_isValid"='t'`
	rows, err := db.Query(UserStrategySql)
	if err != nil {
		logrus.Error("Failed to execute query: ", err)
	}

	var StrategyidList []UserStrategy

	for rows.Next() {
		var UserStrategynew UserStrategy

		err = rows.Scan(&UserStrategynew.Uid, &UserStrategynew.joinTime, &UserStrategynew.Strategyid, &UserStrategynew.actualInvest, &UserStrategynew.apiId)
		StrategyidList = append(StrategyidList, UserStrategynew)
	}
	return StrategyidList
}

func queryUserStrategyEarnings(db *sql.DB, uid string, f_strategyID string, yesterday string) float64 {
	Sql := fmt.Sprintf(`SELECT "f_totalBenefit" FROM "userStrategyEarnings" WHERE f_uid = '%s' and "f_strategyID"='%s' and "f_createTime"='%s'`, uid, f_strategyID, yesterday)
	rows, err := db.Query(Sql)
	if err != nil {
		logrus.Error("Failed to execute query: ", err)
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

func insertEarning(db *sql.DB, dayBenefit float64, totalBenefit float64, uid string, strategyID string) {
	insertSQL := `
		INSERT INTO "userStrategyEarnings" ("f_strategyID", "f_dayBenefit", "f_totalBenefit", "f_uid")
		VALUES ($1, $2, $3,$4)
	`
	// 要插入的数据
	data := []interface{}{strategyID, dayBenefit, totalBenefit, uid}
	// 执行插入操作
	result, err := db.Exec(insertSQL, data...)
	if err != nil {
		logrus.Error(err)
		return
	}

	// 获取受影响的行数
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logrus.Error(err)
		return
	}
	fmt.Printf("rowsAffected = %d", rowsAffected)
}

func (c *UserBenefitService) Run() error {
	logrus.Info("***************************开始每日任务：用户每日投资收益统计***************************")
	// Create DB pool
	db, err := sql.Open("postgres", DB_DSN)
	if err != nil {
		logrus.Error("Failed to open a DB connection: ", err)
		return err
	}
	defer db.Close()
	StrategyidList := queryUserStrategy(db)
	for _, value := range StrategyidList {
		Sql := fmt.Sprintf(`SELECT "f_apiKey", "f_apiSecret" FROM "userBindInfos" WHERE f_id = %s`, value.apiId)
		fmt.Println(Sql)
		rows, err := db.Query(Sql)
		if err != nil {
			logrus.Error("Failed to execute query: ", err)
			return err
		}
		StrategySql := `SELECT "f_coinName" FROM "strategys" WHERE "f_strategyID" = $1`
		fmt.Println(StrategySql, value.Strategyid)
		Strategyrows, err := db.Query(StrategySql, value.Strategyid)
		if err != nil {
			logrus.Error("Failed to execute query: ", err)
			return err
		}
		var apiKey string
		var apiSecret string
		var Asset string
		for rows.Next() {
			err = rows.Scan(&apiKey, &apiSecret)
		}
		api := fmt.Sprintf("apikey:%s, apisecret:%s", apiKey, apiSecret)
		logrus.Info(api)
		if apiKey != "" && apiSecret != "" {
			for Strategyrows.Next() {
				err = Strategyrows.Scan(&Asset)
			}
			logrus.Info("Asset:", Asset)
			// 解密
			apiKey = cryptor.AesSimpleDecrypt(apiKey, types.AesKey)
			apiSecret = cryptor.AesSimpleDecrypt(apiSecret, types.AesKey)

			err1 := errors.New("init error")
			var userData *futures.Account
			for {
				userData, err1 = utils.GetBinanceUMUserData(apiKey, apiSecret)
				if err1 != nil {
					logrus.Info(err1)
				} else {
					logrus.Info("成功请求到币安UM接口")
					break
				}
			}

			if userData == nil {
				logrus.Info("userData is null")
				return nil
			}

			umSum := float64(0)

			for _, asset := range userData.Assets {
				MarginBalanceFloat, err := strconv.ParseFloat(asset.MarginBalance, 64)
				if err != nil {
					logrus.Error(err)
					return err
				}
				if MarginBalanceFloat > 0 {
					price := float64(0)
					if asset.Asset != "USDT" {
						symbols := make([]string, 1)
						symbols[0] = asset.Asset + "USDT"

						prices, err := utils.GetBinancePrice(types.ApiKeySystem, types.ApiSecretSystem, symbols)
						if err != nil {
							logrus.Error(err)
							return err
						}
						if err != nil {
							logrus.Error(prices)
							return err
						}
						price, err = strconv.ParseFloat(prices[0].Price, 64)
						if err != nil {
							logrus.Error(err)
							return err
						}
					} else {
						price = 1
					}
					assetSum := MarginBalanceFloat * float64(price)
					umSum = assetSum + float64(umSum)
				}
			}
			// 累计收入
			totalBenefit := umSum - value.actualInvest
			// 获取当前时间
			now := time.Now()
			// 计算昨天的时间
			yesterday := now.AddDate(0, 0, -1).Format("2006-01-02")
			// 查库中的累计收入
			totalBenefitFloat := queryUserStrategyEarnings(db, value.Uid, value.Strategyid, yesterday)
			// 今日收益
			totalDay := totalBenefit - totalBenefitFloat
			insertEarning(db, totalDay, totalBenefit, value.Uid, value.Strategyid)
		}

	}
	return nil
}
