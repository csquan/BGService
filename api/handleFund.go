package api

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
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
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"net/http"
	"strconv"
	"time"
)

// 简单版本充值：去链上查询这个地址，获取余额和db中最新的一条比对 正规做法：需要爬快 kafka传消息-待迭代
func (a *ApiService) haveFundIn(c *gin.Context) {
	var fundInParam *types.FundInParam

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

	session := a.dbEngine.NewSession()
	err = session.Begin()
	if err != nil {
		return
	}

	//首先插入或修改用户充值记录
	fundInAmount, err := util.ModifyUserFundIn(session, a.dbEngine, fundInParam, userAddr)
	if err != nil {
		err := session.Rollback()
		if err != nil {
			logrus.Fatal(err)
		}
	}
	//下面更新用户资产表
	userAsset, err := db.GetUserAsset(a.dbEngine, fundInParam.Uid)
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}
	fundInDec, err := decimal.NewFromString(fundInAmount)
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}

	if userAsset == nil {
		userAsset = &types.UserAsset{
			Uid:       fundInParam.Uid,
			Network:   fundInParam.Network,
			CoinName:  "usdt",
			Available: fundInAmount,
			Total:     fundInAmount,
		}
	} else {
		dec1, err := decimal.NewFromString(userAsset.Total)
		if err != nil {
			res := util.ResponseMsg(-1, "fail", err)
			c.SecureJSON(http.StatusOK, res)
			return
		}
		userAsset.Total = decimal.Sum(dec1, fundInDec).String()
	}

	_, err = session.Table("userAsset").Insert(userAsset)
	if err != nil {
		err := session.Rollback()
		if err != nil {
			logrus.Fatal(err)
		}
	}

	err = session.Commit()
	if err != nil {
		logrus.Fatal(err)
	}

	res = util.ResponseMsg(0, "success", nil)
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
	//查询uid对应的地址--todo：ADD平台奖励账户地址  应该从平台账户地址打出钱
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

	//下面取出对应私钥签名
	pri, err := db.GetUserKey(a.dbEngine, fromAddr.Addr)
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}
	//先解密
	priDecrypt := util.AesDecrypt(pri.PrivateKey, types.AesKey)

	privateKeyBytes, _ := hex.DecodeString(priDecrypt)

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
	//todo:这里将记录插入提币记录表tx.EnergyUsed 是手续费么？？

	res := util.ResponseMsg(0, "success to send tx", "hash："+hex.EncodeToString(tx.Txid))
	c.SecureJSON(http.StatusOK, res)
	return
}

// 得到用户体验金-从用户体验表中取出即可
func (a *ApiService) getUserExperience(c *gin.Context) {
	//uid, _ := c.Get("Uid")
	//uidFormatted := fmt.Sprintf("%s", uid)

	res := util.ResponseMsg(0, "getUserPlatformFundIn success", nil)
	c.SecureJSON(http.StatusOK, res)
	return
}

// 得到用户佣金--从用户分佣记录表中取出即可
func (a *ApiService) getUserShare(c *gin.Context) {
	//uid, _ := c.Get("Uid")
	//uidFormatted := fmt.Sprintf("%s", uid)

	res := util.ResponseMsg(0, "getUserShare success", nil)
	c.SecureJSON(http.StatusOK, res)
	return
}

// 得到充值记录--转入
func (a *ApiService) getUserPlatformFundIn(c *gin.Context) {
	uid, _ := c.Get("Uid")
	uidFormatted := fmt.Sprintf("%s", uid)
	fundIns, err := db.GetUserAllFundIn(a.dbEngine, uidFormatted)
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}
	var recordOutputs []types.RecordOutput

	for _, fundIn := range *fundIns {
		var recordOutput types.RecordOutput

		recordOutput.Time = fundIn.CreateTime.String()
		recordOutput.Coin = fundIn.Coin
		recordOutput.Type = "Fund IN"
		recordOutput.Amount = fundIn.FundInAmount
		recordOutput.Addr = fundIn.Addr
		recordOutput.Status = "Arrived"

		recordOutputs = append(recordOutputs, recordOutput)
	}

	res := util.ResponseMsg(0, "getUserPlatformFundIn success", recordOutputs)
	c.SecureJSON(http.StatusOK, res)
	return
}

// 得到充值记录--转出
func (a *ApiService) getUserPlatformFundOut(c *gin.Context) {
	uid, _ := c.Get("Uid")
	uidFormatted := fmt.Sprintf("%s", uid)

	//先根据UID查询对应的用户地址

	userAddr, err := db.GetUserAddr(a.dbEngine, uidFormatted)
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}

	fundOuts, err := db.GetUserAllFundOut(a.dbEngine, userAddr.Addr)
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}
	var recordOutputAndGases []types.RecordOutputAndGas

	for _, fundout := range *fundOuts {
		var recordOutputGas types.RecordOutputAndGas

		recordOutputGas.Time = fundout.CreateTime.String()
		recordOutputGas.Coin = fundout.CoinName
		recordOutputGas.Type = "Fund OUT"
		recordOutputGas.Amount = fundout.Amount
		recordOutputGas.Addr = fundout.ToAddr
		recordOutputGas.Status = "Arrived"
		recordOutputGas.Gas = fundout.Gas

		recordOutputAndGases = append(recordOutputAndGases, recordOutputGas)
	}

	res := util.ResponseMsg(0, "getUserPlatformFundOut success", recordOutputAndGases)
	c.SecureJSON(http.StatusOK, res)
	return
}

// 得到分佣记录
func (a *ApiService) getUserPlatformShare(c *gin.Context) {
	uid, _ := c.Get("Uid")
	uidFormatted := fmt.Sprintf("%s", uid)

	userShares, err := db.GetUserAllShare(a.dbEngine, uidFormatted)
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}

	var recordOutputs []types.RecordOutput

	for _, userShare := range *userShares {
		var recordOutput types.RecordOutput

		recordOutput.Time = userShare.CreateTime.String()
		recordOutput.Coin = userShare.CoinName
		recordOutput.Type = "SHARE"
		userShare.Amount = userShare.Amount
		recordOutput.Status = "Arrived"

		recordOutputs = append(recordOutputs, recordOutput)
	}
	res := util.ResponseMsg(0, "getUserPlatformFundOut success", recordOutputs)
	c.SecureJSON(http.StatusOK, res)
	return
}

// 得到体验金记录
func (a *ApiService) getUserPlatformExperience(c *gin.Context) {
	uid, _ := c.Get("Uid")
	uidFormatted := fmt.Sprintf("%s", uid)

	userExperiences, err := db.GetUserAllExperience(a.dbEngine, uidFormatted)
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}

	var expRecordOutputs []types.ExpRecordOutput

	for _, userExperience := range *userExperiences {
		var expRecordOutput types.ExpRecordOutput

		expRecordOutput.Time = userExperience.CreateTime.String()
		expRecordOutput.Coin = userExperience.CoinName

		switch userExperience.Type {
		case "1":
			expRecordOutput.Type = "quart product exp"
		}
		//expRecordOutput.Amount = userExperience.ReceiverSum
		expRecordOutput.Status = "not used"
		expRecordOutput.Valid = userExperience.ValidStartTime + "-" + userExperience.ValidEndTime

		expRecordOutputs = append(expRecordOutputs, expRecordOutput)
	}

	res := util.ResponseMsg(0, "getUserPlatformFundOut success", expRecordOutputs)
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
