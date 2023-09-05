package util

import (
	"context"
	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/sirupsen/logrus"
)

var (
	secretKey = "bd03129b1d27f3818a5ffd363424f9bc6ed655848d063ebfecf220f3037c03da"
	apiKey    = "da7bab67305b2037c103c1c97d7f192c11401606cf3947769340e3a1e4e7e9c6"

	base_future_testnet_binance_url = "https://testnet.binancefuture.com"
)

// U本位合约--得到账户余额
func GetBinanceUserData() (*futures.Account, error) {
	futuresClient := binance.NewFuturesClient(apiKey, secretKey) // USDT-M Futures
	futuresClient.SetApiEndpoint(base_future_testnet_binance_url)

	ret, err := futuresClient.NewGetAccountService().Do(context.Background())

	if err != nil {
		logrus.Info(err)
		return nil, err
	}

	return ret, nil
}
