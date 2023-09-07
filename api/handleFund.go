package api

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/ethereum/BGService/db"
	"github.com/ethereum/BGService/types"
	"github.com/ethereum/BGService/util"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/fbsobreira/gotron-sdk/pkg/client"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/api"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"net/http"
	"strconv"
	"time"
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

// 金额以sun为单位 1trx=1000,000 sun
func (a *ApiService) fundOut(c *gin.Context) {
	var fundOutParam types.FundOutParam

	err := c.BindJSON(&fundOutParam)
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}

	var conn *client.GrpcClient
	conn = client.NewGrpcClient(a.config.Endpoint.Trx)
	err = conn.Start(grpc.WithInsecure())
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}
	//查询uid对应的地址
	fromAddr, err := db.GetUserAddr(a.dbEngine, fundOutParam.Uid)
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}

	amount, err := strconv.ParseInt(fundOutParam.Amount, 10, 64)
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}

	var tx *api.TransactionExtention
	for {
		tx, err = conn.Transfer(fromAddr.Addr, fundOutParam.ToAddr, amount)
		if err != nil {
			time.Sleep(100 * time.Millisecond)
		} else {
			break
		}
	}

	//开始签名
	rawData, err := proto.Marshal(tx.Transaction.GetRawData())

	h256h := sha256.New()
	h256h.Write(rawData)
	hash := h256h.Sum(nil)

	// btcec.PrivKeyFromBytes only returns a secret key and public key

	//下面取出对应私钥签名，todo：移动到单独的私钥服务器
	pri, err := db.GetUserKey(a.dbEngine, fromAddr.Addr)
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}

	privateKeyBytes, _ := hex.DecodeString(pri.PrivateKey)
	sk, _ := btcec.PrivKeyFromBytes(privateKeyBytes)

	signature, err := crypto.Sign(hash, sk.ToECDSA())

	tx.Transaction.Signature = append(tx.Transaction.Signature, signature)

	for {
		_, err = conn.Broadcast(tx.Transaction)
		if err != nil {
			time.Sleep(100 * time.Millisecond)
		} else {
			break
		}
	}

	res := util.ResponseMsg(0, "success to send tx", "hash："+hex.EncodeToString(tx.Txid))
	c.SecureJSON(http.StatusOK, res)
	return
}

// todo:吧钱转出来
//func (a *ApiService) fundOut(c *gin.Context) {
//	uid := "47055457103956"                        //用户uid
//	toAddr := "TFSoDRmsSP289NjDp3mzAc2Rgi2ZGheiqD" //我的测试地址
//	logrus.Info(uid)
//	var conn *client.GrpcClient
//	conn = client.NewGrpcClient(a.config.Endpoint.Trx)
//	err := conn.Start(grpc.WithInsecure())
//	if err != nil {
//		res := util.ResponseMsg(-1, "fail", err)
//		c.SecureJSON(http.StatusOK, res)
//		return
//	}
//
//	//conn.SetAPIKey(apiKey)??
//	fromAddr := "TWK9oxSqfVc5J7GCCFj3MMYsh9w9Vce3tt" //用户地址
//	valueInt := int64(1000)
//
//	ks, acct, err := store.UnlockedKeystore(fromAddr, "")
//	if err != nil {
//		res := util.ResponseMsg(-1, "fail", err)
//		c.SecureJSON(http.StatusOK, res)
//		return
//	}
//
//	tx, err := conn.Transfer(fromAddr, toAddr, valueInt)
//	if err != nil {
//		res := util.ResponseMsg(-1, "fail", err.Error())
//		c.SecureJSON(http.StatusOK, res)
//		return
//	}
//	//mnemoric := "crystal gate zoo sock renew puppy process one cricket beach barely perfect praise side frost fat paddle age occur carbon clip claw yard yellow"
//	//ks1 := store.FromAccountName("csquan1")
//	ctrlr := transaction.NewController(conn, ks, acct, tx.Transaction)
//
//	err = ctrlr.ExecuteTransaction()
//	if err != nil {
//		res := util.ResponseMsg(-1, "fail", err)
//		c.SecureJSON(http.StatusOK, res)
//		return
//	}
//
//	res := util.ResponseMsg(0, "success to send tx", err)
//	c.SecureJSON(http.StatusOK, res)
//	return
//}
