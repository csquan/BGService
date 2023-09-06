package api

import (
	"github.com/ethereum/BGService/types"
	"github.com/ethereum/BGService/util"
	"github.com/fbsobreira/gotron-sdk/pkg/client"
	"github.com/fbsobreira/gotron-sdk/pkg/client/transaction"
	"github.com/fbsobreira/gotron-sdk/pkg/store"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net/http"
)

// 去链上查询这个地址，获取交易记录
func (a *ApiService) fundIn(c *gin.Context) {
	symbols := c.Query("symbols")

	res := types.HttpRes{}

	url := base_binance_url + "api/v3/ticker/price?symbols=" + symbols

	data, err := util.Get(url)
	if err != nil {
		logrus.Info("获取币价失败", err)

		res.Code = 0
		res.Message = "成功获取价格"
		res.Data = err

		c.SecureJSON(http.StatusOK, res)
		return
	}
	res.Code = 0
	res.Message = "成功获取价格"
	res.Data = data

	c.SecureJSON(http.StatusOK, res)
	return
}

// 插入发送交易任务表
func (a *ApiService) fundOut(c *gin.Context) {
	uid := "47055457103956"                        //用户uid
	toAddr := "TFSoDRmsSP289NjDp3mzAc2Rgi2ZGheiqD" //我的测试地址
	logrus.Info(uid)
	var conn *client.GrpcClient
	conn = client.NewGrpcClient(a.config.Endpoint.Trx)
	err := conn.Start(grpc.WithInsecure())
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}

	//conn.SetAPIKey(apiKey)??
	fromAddr := "TWK9oxSqfVc5J7GCCFj3MMYsh9w9Vce3tt" //用户地址
	valueInt := int64(1000)

	ks, acct, err := store.UnlockedKeystore(fromAddr, "")
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}

	tx, err := conn.Transfer(fromAddr, toAddr, valueInt)
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err.Error())
		c.SecureJSON(http.StatusOK, res)
		return
	}
	//mnemoric := "crystal gate zoo sock renew puppy process one cricket beach barely perfect praise side frost fat paddle age occur carbon clip claw yard yellow"
	//ks := store.FromAccountName("csquan1")
	ctrlr := transaction.NewController(conn, ks, acct, tx.Transaction)

	err = ctrlr.ExecuteTransaction()
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}

	res := util.ResponseMsg(0, "success to send tx", err)
	c.SecureJSON(http.StatusOK, res)
	return
}
