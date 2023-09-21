package util

import (
	"context"
	"fmt"
	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/delivery"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/ethereum/BGService/types"
	"github.com/sirupsen/logrus"
)

var (
	base_future_testnet_binance_url = "https://testnet.binancefuture.com"
	base_future_binance_url         = "https://fapi.binance.com"
)

// U本位合约--得到账户余额
func GetBinanceUMUserData(apiKey, apiSecret string) (*futures.Account, error) {
	binanceClient := futures.NewClient(apiKey, apiSecret) // USDT-M Futures
	binanceClient.SetApiEndpoint(base_future_binance_url)

	ret, err := binanceClient.NewGetAccountService().Do(context.Background())

	if err != nil {
		logrus.Info(err)
		return nil, err
	}
	return ret, nil
}

// 币本位合约--得到账户余额
func GetBinanceCMUserData(apiKey, apiSecret string) (*delivery.Account, error) {
	binanceClient := delivery.NewClient(apiKey, apiSecret) // USDT-M Futures
	binanceClient.SetApiEndpoint(base_future_testnet_binance_url)

	ret, err := binanceClient.NewGetAccountService().Do(context.Background())

	if err != nil {
		logrus.Info(err)
		return nil, err
	}

	return ret, nil
}

// 现货--得到账户余额
func GetBinanceSpotUserData(apiKey, apiSecret string) ([]*binance.CoinInfo, error) {
	client := binance.NewClient(apiKey, apiSecret)
	client.SetApiEndpoint(types.Base_binance_url)

	ret, err := client.NewGetAllCoinsInfoService().Do(context.Background())

	if err != nil {
		logrus.Info(err)
		return nil, err
	}

	return ret, nil
}

// U本位合约--得到交易记录
func GetBinanceUMUserTxHistory(apiKey, apiSecret, symbol string, limit int) ([]*futures.AccountTrade, error) {
	binanceClient := futures.NewClient(apiKey, apiSecret) // USDT-M Futures
	binanceClient.SetApiEndpoint(base_future_binance_url)

	listAccountTrades, err := binanceClient.NewListAccountTradeService().Symbol(symbol).Limit(limit).Do(context.Background())

	if err != nil {
		logrus.Info(err)
		return nil, err
	}

	return listAccountTrades, nil
}

// 1D 5分钟--288条 1W 1h--168条 1M 6h
func GetBinanceKlinesHistory(interval string, startTime int64, endTime int64, KlineType int, symbol string) ([]*binance.Kline, error) {
	//首先参数检验
	err := CheckGetKlineParam(interval, startTime, endTime, KlineType)
	if err != nil {
		logrus.Info(err)
		return nil, err
	}
	var klines []*binance.Kline
	client := binance.NewClient(types.ApiKey, types.ApiSecret)
	switch KlineType {
	case 1:
		//1D 5分钟--288
		klines, err = client.NewKlinesService().Symbol(symbol).
			Interval(interval).StartTime(startTime).EndTime(endTime).Do(context.Background())
	case 2:
		//1W 1h--168条
		klines, err = client.NewKlinesService().Symbol(symbol).
			Interval(interval).StartTime(startTime).EndTime(endTime).Do(context.Background())
	case 3:
		//1M 6h--120条左右
		klines, err = client.NewKlinesService().Symbol(symbol).
			Interval(interval).StartTime(startTime).EndTime(endTime).Do(context.Background())
	}
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return klines, nil
}
