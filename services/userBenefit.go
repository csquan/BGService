package services

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/adshao/go-binance/v2"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"log"
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
		log.Fatal("Failed to execute query: ", err)
	}

	var StrategyidList []UserStrategy

	for rows.Next() {
		var UserStrategynew UserStrategy

		err = rows.Scan(&UserStrategynew.Uid, &UserStrategynew.joinTime, &UserStrategynew.Strategyid, &UserStrategynew.actualInvest, &UserStrategynew.apiId)
		StrategyidList = append(StrategyidList, UserStrategynew)
	}
	return StrategyidList
}

func queryUserStrategyEarnings(db *sql.DB, uid string, stragetyID string) float64 {
	Sql := fmt.Sprintf(`SELECT "f_totalBenefit" FROM "userStrategyEarnings" WHERE f_uid = '%s' and "f_stragetyID"='%s'`, uid, stragetyID)
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

func updateEarning(db *sql.DB, dayBenefit float64, totalBenefit float64, uid string, stragetyID string) {
	stmt, err := db.Prepare(`update "userStrategyEarnings" set "f_dayBenefit"=$1, "f_totalBenefit"=$2 where "f_uid"=$3 and "f_stragetyID"=$4`)
	if err != nil {
		panic(err)
	}
	res, err := stmt.Exec(dayBenefit, totalBenefit, uid, stragetyID)
	if err != nil {
		panic(err)
	}

	fmt.Printf("res = %d", res)
}

func (c *UserBenefitService) Run() error {
	// Create DB pool
	db, err := sql.Open("postgres", DB_DSN)
	if err != nil {
		log.Fatal("Failed to open a DB connection: ", err)
	}
	defer db.Close()
	StrategyidList := queryUserStrategy(db)
	for _, value := range StrategyidList {
		Sql := fmt.Sprintf(`SELECT "f_apiKey", "f_apiSecret" FROM "userBindInfos" WHERE f_id = %s`, value.apiId)
		fmt.Println(Sql)
		rows, err := db.Query(Sql)
		if err != nil {
			log.Fatal("Failed to execute query: ", err)
		}
		StrategySql := `SELECT "f_coinName" FROM "strategys" WHERE "f_strategyID" = $1`
		fmt.Println(StrategySql, value.Strategyid)
		Strategyrows, err := db.Query(StrategySql, value.Strategyid)
		if err != nil {
			log.Fatal("Failed to execute query: ", err)
		}
		var apiKey string
		var apiSecret string
		var Asset string
		for rows.Next() {
			err = rows.Scan(&apiKey, &apiSecret)
		}
		api := fmt.Sprintf("apikey:%s, apisecret:%s", apiKey, apiSecret)
		fmt.Println(api)
		if apiKey != "" && apiSecret != "" {
			for Strategyrows.Next() {
				err = Strategyrows.Scan(&Asset)
			}
			fmt.Println("Asset:", Asset)
			futuresClient := binance.NewFuturesClient(apiKey, apiSecret) // USDT-M Futures
			futuresClient.SetApiEndpoint(base_future_testnet_binance_url)
			userData, err := futuresClient.NewGetAccountService().Do(context.Background())
			if err != nil {
				logrus.Info(err)
			}
			var MarginBalance string
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
			// 累计收入
			totalBenefit := MarginBalanceFloat - value.actualInvest
			// 查库中的累计收入
			totalBenefitFloat := queryUserStrategyEarnings(db, value.Uid, value.Strategyid)
			// 今日收益
			totalDay := totalBenefit - totalBenefitFloat
			updateEarning(db, totalDay, totalBenefit, value.Uid, value.Strategyid)
		}

	}
	return nil
}
