package services

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"database/sql"
	"encoding/base64"
	"fmt"
	"github.com/LinkinStars/go-scaffold/contrib/cryptor"
	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/ethereum/BGService/types"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"log"
	"strconv"
	"strings"
	"time"
)

type UserTxRecordService struct {
}

func NewUserTxRecordService() *UserTxRecordService {
	return &UserTxRecordService{}
}

func (c *UserTxRecordService) Name() string {
	return "UserTxRecordService"
}

const (
	DbDsn                       = "postgres://postgres:12345@127.0.0.1:5432/postgres?sslmode=disable"
	AesKey                      = "cure-d111y=1ziukr07k*!r$q=zcgto%" //AES密钥
	baseFutureTestnetBinanceUrl = "https://testnet.binancefuture.com"
)

type RUserStrategy struct {
	Strategyid   string
	Uid          string
	joinTime     time.Time
	actualInvest float64
	apiId        string
}

func RQueryUserStrategy(db *sql.DB) []RUserStrategy {
	UserStrategySql := `SELECT "f_uid", "f_joinTime", "f_strategyID", "f_actualInvest", "f_apiId" FROM "userStrategy" WHERE "f_isValid"='t'`
	rows, err := db.Query(UserStrategySql)
	if err != nil {
		log.Fatal("Failed to execute query: ", err)
	}

	var StrategyidList []RUserStrategy

	for rows.Next() {
		var UserStrategynew RUserStrategy

		err = rows.Scan(&UserStrategynew.Uid, &UserStrategynew.joinTime, &UserStrategynew.Strategyid, &UserStrategynew.actualInvest, &UserStrategynew.apiId)
		StrategyidList = append(StrategyidList, UserStrategynew)
	}
	return StrategyidList
}

func strategy(db *sql.DB, Strategyid string) (string, string) {
	// 策略数据查询
	StrategySql := `SELECT "f_coinName", "f_strategyName" FROM "strategys" WHERE "f_strategyID" = $1`
	fmt.Println(StrategySql, Strategyid)
	Strategyrows, err := db.Query(StrategySql, Strategyid)
	if err != nil {
		log.Fatal("Failed to execute query: ", err)
	}
	var coinName string
	var strategyName string
	for Strategyrows.Next() {
		err = Strategyrows.Scan(&coinName, &strategyName)
	}
	return coinName, strategyName
}

func userBindInfo(db *sql.DB, apiId string) (string, string) {
	Sql := fmt.Sprintf(`SELECT "f_apiKey", "f_apiSecret" FROM "userBindInfos" WHERE f_id = %s`, apiId)
	fmt.Println(Sql)
	rows, err := db.Query(Sql)
	if err != nil {
		log.Fatal("Failed to execute query: ", err)
	}
	var apiKey string
	var apiSecret string
	for rows.Next() {
		err = rows.Scan(&apiKey, &apiSecret)
	}
	api := fmt.Sprintf("apikey:%s, apisecret:%s", apiKey, apiSecret)
	fmt.Println(api)
	return apiKey, apiSecret
}

// 去码
func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// AES 解密
func AesDecrypt(cryted string, key string) string {
	// cryted需要解密的字符串，key加密密钥
	//使用RawURLEncoding 不要使用StdEncoding
	//不要使用StdEncoding  放在url参数中回导致错误
	crytedByte, _ := base64.RawURLEncoding.DecodeString(cryted)
	k := []byte(key)

	// 分组秘钥
	block, err := aes.NewCipher(k)
	if err != nil {
		panic(fmt.Sprintf("key 长度必须 16/24/32长度: %s", err.Error()))
	}
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 加密模式
	blockMode := cipher.NewCBCDecrypter(block, k[:blockSize])
	// 创建数组
	orig := make([]byte, len(crytedByte))
	// 解密
	blockMode.CryptBlocks(orig, crytedByte)
	// 去补全码
	orig = PKCS7UnPadding(orig)
	return string(orig)
}

// U本位合约--交易历史
func GetBinanceUMUserTxHistory(symbol string, limit int, apiKey string, secretKey string) ([]*futures.AccountTrade, error) {
	futuresClient := binance.NewFuturesClient(apiKey, secretKey) // USDT-M Futures
	futuresClient.SetApiEndpoint(baseFutureTestnetBinanceUrl)

	listAccountTrades, err := futuresClient.NewListAccountTradeService().Symbol(symbol).Limit(limit).Do(context.Background())

	if err != nil {
		logrus.Info(err)
		return nil, err
	}

	return listAccountTrades, nil
}

func insertRecords(db *sql.DB, orderId string, Uid string, address string, Strategyid string, side string, behavior string, t string) error {
	stmt, err := db.Prepare(`INSERT INTO "transactionRecords"("f_orderId", "f_uid", "f_address", "f_strategyID", "f_action", "f_behavior", "f_time") 
									VALUES($1, $2, $3, $4, $5, $6, $7)`)
	if err != nil {
		panic(err)
	}
	fmt.Println(orderId)
	res, err := stmt.Exec(orderId, Uid, address, Strategyid, side, behavior, t)
	if err != nil {
		panic(err)
		return err
	}

	fmt.Printf("res = %d", res)
	return nil
}

func (c *UserTxRecordService) Run() error {
	//Create DB pool
	db, err := sql.Open("postgres", DbDsn)
	if err != nil {
		log.Fatal("Failed to open a DB connection: ", err)
	}
	defer db.Close()
	// 用户策略查询
	StrategyidList := RQueryUserStrategy(db)
	for _, value := range StrategyidList {
		// 策略信息查询
		// coinName币种、strategyName交易对
		// 一期只做U本位合约
		_, strategyName := strategy(db, value.Strategyid)
		if value.apiId == "" {
			continue
		}
		// 用户api查询,查到的key是加密的需要解密
		apiKey, apiSecret := userBindInfo(db, value.apiId)
		if apiKey == "" && apiSecret == "" {
			continue
		}
		// 解密
		apiKey = cryptor.AesSimpleDecrypt(apiKey, types.AesKey)
		apiSecret = cryptor.AesSimpleDecrypt(apiSecret, types.AesKey)

		// 交易历史
		symbol := strings.Replace(strategyName, "/", "", 1)
		history, err := GetBinanceUMUserTxHistory(symbol, 1000, apiKey, apiSecret)
		if err != nil {
			return err
		}
		fmt.Println(history)
		for _, historyvalue := range history {
			timestamp := historyvalue.Time
			t := time.Unix(timestamp/1000, 0)
			side := historyvalue.Side                 // 动作
			orderId := historyvalue.OrderID           // 订单编号
			positionSide := historyvalue.PositionSide // 持仓
			price := historyvalue.Price               // 成交价
			quoteQty := historyvalue.QuoteQuantity    // 成交额
			// 行为
			behavior := fmt.Sprintf("在【%s】, 【%s】以均价【%s】 【%s】成交【%s】", "binance", positionSide, price, side, quoteQty)
			timeNow := time.Now()
			yesterday := timeNow.AddDate(0, 0, -1)
			if yesterday.Unix() < t.Unix() && t.Unix() < timeNow.Unix() {
				fmt.Println(behavior)
				str := strconv.FormatInt(timestamp, 10)
				strorderId := strconv.FormatInt(orderId, 10)
				err := insertRecords(db, strorderId+str, value.Uid, "binance", value.Strategyid, string(side), behavior, t.String())
				if err != nil {
					continue
				}
			}
		}
	}
	return nil
}
