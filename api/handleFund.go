package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/ethereum/BGService/db"
	"github.com/ethereum/BGService/types"
	"github.com/ethereum/BGService/util"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/fbsobreira/gotron-sdk/pkg/client"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/api"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"net/http"
	"strconv"
	"time"
)

const base_tron_url = "https://api.trongrid.io"

// 简单版本充值：去链上查询这个地址，获取交易记录 正规做法：需要爬快 kafka传消息
func (a *ApiService) haveFundIn(c *gin.Context) {
	var fundInParam types.FundInParam
	var UserFundIns []types.UserFundIn

	err := c.BindJSON(&fundInParam)
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}

	userAddr, err := db.GetUserAddr(a.dbEngine, fundInParam.Uid)
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}

	res := types.HttpRes{}

	//这里简单逻辑修改下：直接取最新的区块，然后取余额，与数据库中上次修改的余额相减，得到本次充值
	url := base_tron_url + "/wallet/getaccount"

	accountParam := types.AccountParam{
		Address: userAddr.Addr,
		Visible: true,
	}

	bodyStr, err := json.Marshal(accountParam)
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}

	str1, err := util.Post(url, bodyStr)
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}
	balance := gjson.Get(str1, "balance")

	url = base_tron_url + "/v1/accounts/" + userAddr.Addr + "/transactions"

	//取出用户最近的充值记录表
	userFundIn, err := db.GetUserFundIn(a.dbEngine, fundInParam.Uid, fundInParam.Network)
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}
	var DBBlockHeight int64

	if userFundIn == nil {
		//insert
	} else {
		if userFundIn.IsCollect == true { //发生过归集 本次充值金额为 归集后链上余额 + 目前的链上余额

		} else { //未发生归集 本次充值金额为 本次充值后链上余额-上次充值后链上余额

		}
	}

	//这里用事务更新UserFundIns进db 这里的金额作为本次充值的金额
	session := a.dbEngine.NewSession()
	err = session.Begin()
	if err != nil {
		return
	}

	for _, fundIn := range UserFundIns {
		//首先插入用户充值记录
		_, err = session.Table("platformExperience").Insert(fundIn)
		if err != nil {
			err := session.Rollback()
			if err != nil {
				return
			}
			logrus.Fatal(err)
		}
		//下面应该更新用户资产表
		userAsset, err := db.GetUserAsset(a.dbEngine, fundInParam.Uid)
		if err != nil {

		}
		dec1, err := decimal.NewFromString(userAsset.Total)
		if err != nil {

		}
		dec2, err := decimal.NewFromString(fundIn.Amount)
		if err != nil {

		}

		dec := decimal.Sum(dec1, dec2)
		userAsset.Total = dec.String()

		_, err = session.Table("userAssets").Update(userAsset)
		if err != nil {
			err := session.Rollback()
			if err != nil {
				return
			}
			logrus.Fatal(err)
		}
	}

	err = session.Commit()
	if err != nil {
		logrus.Fatal(err)
	}

	res = util.ResponseMsg(0, "success", array)
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
