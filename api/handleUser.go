package api

import (
	"fmt"
	"github.com/ethereum/BGService/db"
	"github.com/ethereum/BGService/types"
	"github.com/ethereum/BGService/util"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

func (a *ApiService) info(c *gin.Context) {
	uid, _ := c.Get("Uid")
	// 根据uid查询用户信息
	uidFormatted := fmt.Sprintf("%s", uid)
	err, user := db.QuerySecret(a.dbEngine, uidFormatted)
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}
	// 邀请人数查询
	err, inviteUser := db.QueryInviteNum(a.dbEngine, user.InvitationCode)
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}
	// api绑定信息
	userBindInfos, err := db.GetUserBindInfos(a.dbEngine, uidFormatted)
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}
	var bindNum bool
	if len(userBindInfos.Uid) < 1 {
		bindNum = false
	} else {
		bindNum = true
	}
	body := make(map[string]interface{})
	body["uid"] = user.Uid
	body["userName"] = user.UserName
	body["isBindGoogle"] = user.IsBindGoogle
	body["isIDVerify"] = user.IsIDVerify
	body["mobile"] = user.Mobile
	body["invitation"] = len(inviteUser)
	body["apiBinding"] = bindNum
	body["email"] = user.MailBox
	body["inviteCode"] = user.InvitationCode
	res := util.ResponseMsg(1, "success", body)
	c.SecureJSON(http.StatusOK, res)
	return
}

func (a *ApiService) myApi(c *gin.Context) {
	uid, _ := c.Get("Uid")
	// 根据uid查询用户信息
	uidFormatted := fmt.Sprintf("%s", uid)
	userBindInfos, err := db.GetAllUserBindInfos(a.dbEngine, uidFormatted)
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}
	body := make(map[string]interface{})
	var allOkCex []interface{}
	var allBinanceCex []interface{}
	for _, value := range userBindInfos {
		if value.Cex == "okex" {
			oneCex := make(map[string]interface{})
			oneCex["id"] = value.ID
			oneCex["cex"] = value.Cex
			oneCex["apiKey"] = value.ApiKey
			oneCex["secretKey"] = value.ApiSecret
			oneCex["passphrase"] = value.Passphrase
			oneCex["account"] = value.Account
			oneCex["alias"] = value.Alias
			oneCex["synchronizeTime"] = value.SynchronizeTime
			oneCex["permission"] = value.Permission
			allOkCex = append(allOkCex, oneCex)
		}
		if value.Cex == "binance" {
			oneCex := make(map[string]interface{})
			oneCex["id"] = value.ID
			oneCex["cex"] = value.Cex
			oneCex["apiKey"] = value.ApiKey
			oneCex["secretKey"] = value.ApiSecret
			oneCex["passphrase"] = value.Passphrase
			oneCex["account"] = value.Account
			oneCex["alias"] = value.Alias
			oneCex["synchronizeTime"] = value.SynchronizeTime
			oneCex["permission"] = value.Permission
			allBinanceCex = append(allBinanceCex, oneCex)
		}
	}
	body["okex"] = allOkCex
	body["binance"] = allBinanceCex
	res := util.ResponseMsg(1, "success", body)
	c.SecureJSON(http.StatusOK, res)
	return
}

func (a *ApiService) bindingApi(c *gin.Context) {
	uid, _ := c.Get("Uid")
	// 根据uid查询用户信息
	uidFormatted := fmt.Sprintf("%s", uid)
	var payload *types.UserBindInfoInput
	if err := c.ShouldBindJSON(&payload); err != nil {
		logrus.Error(err)
		res := util.ResponseMsg(-1, "fail", err.Error())
		c.SecureJSON(http.StatusOK, res)
		return
	}
	userBindInfos, err := db.GetApiKeyUserBindInfos(a.dbEngine, payload.ApiKey)
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}
	// 校验apikey是否已绑定
	if userBindInfos != nil {
		res := util.ResponseMsg(-1, "fail", "apiKey is Bound")
		c.SecureJSON(http.StatusOK, res)
		return
	}
	// todo 缺一个查询此apikey交易权限
	nowTime := time.Now()
	UserBindInfo := types.InsertUserBindInfo{
		Uid:             uidFormatted,
		Cex:             payload.Cex,
		ApiKey:          payload.ApiKey,
		ApiSecret:       payload.ApiSecret,
		Passphrase:      payload.Passphrase,
		Alias:           payload.Alias,
		Account:         payload.Account,
		SynchronizeTime: nowTime,
		Permission:      true,
	}
	if err := db.InsertUserBindInfo(a.dbEngine, &UserBindInfo); err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}
	body := make(map[string]interface{})
	res := util.ResponseMsg(1, "success", body)
	c.SecureJSON(http.StatusOK, res)
	return
}

func (a *ApiService) unbindingApi(c *gin.Context) {
	uid, _ := c.Get("Uid")
	// 根据uid查询用户信息
	uidFormatted := fmt.Sprintf("%s", uid)
	apiId := c.Query("id")
	id, err := strconv.Atoi(apiId)
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}

	userBindInfos, err := db.GetIdUserBindInfos(a.dbEngine, uidFormatted, apiId)
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}
	// 校验apikey是否存在
	if userBindInfos == nil {
		res := util.ResponseMsg(-1, "fail", "apiKey is not exist")
		c.SecureJSON(http.StatusOK, res)
		return
	}
	err = db.DeleteUserBindInfo(a.dbEngine, id)
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}
	body := make(map[string]interface{})
	res := util.ResponseMsg(1, "success", body)
	c.SecureJSON(http.StatusOK, res)
	return
}
