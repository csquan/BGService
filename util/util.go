package util

import (
	"math/big"
)

type CoinInfo struct {
	Name         string `yaml:"name"`
	Symbol       string `yaml:"symbol"`
	ContractAddr string `yaml:"contract_addr"`
	Decimal      string `yaml:"decimal"`
	Total_Supply string `yaml:"total_supply"`
}

type HolderInfo struct {
	Addr          string `yaml:"addr"`
	Balance       string `yaml:"balance"`
	Contract_addr string `yaml:"contract_addr"`
}

type HistoryInfo struct {
	Symbol  string `yaml:"symbol"`
	Time    string `yaml:"time"`
	Balance string `yaml:"balance"`
	Op      string `yaml:"op"`
	Amount  string `yaml:"amount"`
	Params  string `yaml:"params"`
}

type BlockRange struct {
	BeginBlock *big.Int
	EndBlock   *big.Int
}
