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
