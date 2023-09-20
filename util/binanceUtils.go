package util

import (
	"context"
	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/delivery"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/ethereum/BGService/types"
	"github.com/sirupsen/logrus"
)

var (
	secretKey                       = "bd03129b1d27f3818a5ffd363424f9bc6ed655848d063ebfecf220f3037c03da"
	apiKey                          = "da7bab67305b2037c103c1c97d7f192c11401606cf3947769340e3a1e4e7e9c6"
	base_future_testnet_binance_url = "https://testnet.binancefuture.com"
	base_future_binance_url         = "https://fapi.binance.com"
)

// U本位合约--得到账户余额--todo:查询用户对应得真实APIKEY APISECRET
func GetBinanceUMUserData() (*futures.Account, error) {
	//binanceClient := futures.NewClient(apiKey, secretKey) // USDT-M Futures
	//binanceClient.SetApiEndpoint(base_future_testnet_binance_url)
	binanceClient := futures.NewClient(types.ApiKey, types.ApiSecret) // USDT-M Futures
	binanceClient.SetApiEndpoint(base_future_binance_url)

	ret, err := binanceClient.NewGetAccountService().Do(context.Background())

	if err != nil {
		logrus.Info(err)
		return nil, err
	}
	return ret, nil
}

// 币本位合约--得到账户余额
func GetBinanceCMUserData() (*delivery.Account, error) {
	delivery.UseTestnet = true
	binanceClient := delivery.NewClient(apiKey, secretKey) // USDT-M Futures
	binanceClient.SetApiEndpoint(base_future_testnet_binance_url)

	ret, err := binanceClient.NewGetAccountService().Do(context.Background())

	if err != nil {
		logrus.Info(err)
		return nil, err
	}

	return ret, nil
}

// 现货--得到账户余额
func GetBinanceSpotUserData() ([]*binance.CoinInfo, error) {
	//binance.UseTestnet = true
	client := binance.NewClient(types.ApiKey, types.ApiSecret)
	client.SetApiEndpoint(types.Base_binance_url)

	ret, err := client.NewGetAllCoinsInfoService().Do(context.Background())

	if err != nil {
		logrus.Info(err)
		return nil, err
	}

	return ret, nil
}

// U本位合约--得到交易记录
func GetBinanceUMUserTxHistory(symbol string, limit int) ([]*futures.AccountTrade, error) {
	binanceClient := futures.NewClient(types.ApiKey, types.ApiSecret) // USDT-M Futures
	binanceClient.SetApiEndpoint(base_future_binance_url)

	listAccountTrades, err := binanceClient.NewListAccountTradeService().Symbol(symbol).Limit(limit).Do(context.Background())

	if err != nil {
		logrus.Info(err)
		return nil, err
	}

	return listAccountTrades, nil
}
