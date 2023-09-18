package util

import (
	"math/big"
	"strings"
)

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

func RemoveElement(str string, elem string) string {
	return strings.Replace(str, elem, "", 1)
}
