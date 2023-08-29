package api

import (
	"github.com/ethereum/api-in/types"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
)

const ADDRLEN = 42

const Ok = 0

func checkAddr(addr string) error {
	if addr[:2] != "0x" {
		return errors.New("addr must start with 0x")
	}
	if len(addr) != ADDRLEN {
		return errors.New("addr len wrong ,must 40")
	}
	return nil
}

// 首先查询balance_erc20表，得到地址持有的代币合约地址，然后根据代币合约地址查erc20_info表
func (a *ApiService) query(c *gin.Context) {
	addr := c.Param("contractAddr")
	res := types.HttpRes{}

	err := checkAddr(addr)
	if err != nil {
		logrus.Error(err)
		res.Code = http.StatusBadRequest
		res.Message = err.Error()
		c.SecureJSON(http.StatusBadRequest, res)
		return
	}
	res.Code = Ok
	res.Message = "success"
	res.Data = ""
	c.SecureJSON(http.StatusOK, res)
}

// 首先查询balance_erc20表，得到地址持有的代币合约地址，然后根据代币合约地址查erc20_info表
func (a *ApiService) getCoinHistory(c *gin.Context) {
	addr := c.Param("contractAddr")
	res := types.HttpRes{}

	err := checkAddr(addr)
	if err != nil {
		logrus.Error(err)
		res.Code = http.StatusBadRequest
		res.Message = err.Error()
		c.SecureJSON(http.StatusBadRequest, res)
		return
	}
	res.Code = Ok
	res.Message = "success"
	res.Data = ""
	c.SecureJSON(http.StatusOK, res)
}
