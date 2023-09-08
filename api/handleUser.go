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

func (a *ApiService) invite(c *gin.Context) {
	invitationCode, _ := c.Get("invitationCode")
	// 邀请码
	InviteCode := fmt.Sprintf("%s", invitationCode)
	// 邀请数量
	total := c.Query("total")
	totalInt, err := strconv.Atoi(total)
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}
	err, inviteUser := db.QueryInviteNumLimit(a.dbEngine, InviteCode, totalInt)
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}
	var inviteUserList []interface{}
	if len(inviteUser) > 0 {
		for _, value := range inviteUser {
			inviteUserInfo := make(map[string]interface{})
			inviteUserInfo["username"] = value.UserName
			inviteUserInfo["register"] = true
			inviteUserInfo["isBindGoogle"] = value.IsBindGoogle
			inviteUserInfo["createtime"] = value.CreateTime
			inviteUserList = append(inviteUserList, inviteUserInfo)
		}
	}
	body := make(map[string]interface{})
	body["total"] = len(inviteUser)
	body["List"] = inviteUserList
	res := util.ResponseMsg(1, "success", body)
	c.SecureJSON(http.StatusOK, res)
	return
}

func (a *ApiService) inviteRanking(c *gin.Context) {
	uid, _ := c.Get("Uid")
	// 根据uid查询用户信息
	uidFormatted := fmt.Sprintf("%s", uid)
	// 数量
	total := c.Query("total")
	totalInt, err := strconv.Atoi(total)
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}
	err, inviteUserNum := db.QueryClaimRewardNumber(a.dbEngine)
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}
	var inviteUserList []interface{}
	var myPlaced int
	var myCommission int
	if len(inviteUserNum) > 0 {
		for i := 0; i < len(inviteUserNum); i++ {
			value := inviteUserNum[i]
			// 邀请到人的情况能查到排名
			if value.Uid == uidFormatted {
				myPlaced = i + 1
				myCommission = value.ClaimRewardNumber * 10
			}
			inviteUserInfo := make(map[string]interface{})
			inviteUserInfo["placed"] = i + 1
			inviteUserInfo["username"] = value.UserName
			inviteUserInfo["commission"] = value.ClaimRewardNumber * 10
			inviteUserList = append(inviteUserList, inviteUserInfo)
		}
	}
	// 没邀请到人的情况排名在邀请人的最后一名，佣金为0
	if myPlaced == 0 && myCommission == 0 {
		myPlaced = len(inviteUserNum) + 1
		myCommission = 0
	}
	var inviteUserListRes []interface{}
	if len(inviteUserList) < totalInt {
		inviteUserListRes = inviteUserList
	} else {
		inviteUserListRes = inviteUserList[:totalInt]
	}
	body := make(map[string]interface{})
	body["total"] = len(inviteUserListRes)
	body["ranking"] = inviteUserListRes
	body["myPlaced"] = myPlaced
	body["myCommission"] = myCommission
	res := util.ResponseMsg(1, "success", body)
	c.SecureJSON(http.StatusOK, res)
	return
}

func (a *ApiService) getStrategy(c *gin.Context) {
	uid := c.Query("uid")

	userStrategys, err := db.GetUserStrategys(a.dbEngine, uid)
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}

	res := util.ResponseMsg(1, "success", userStrategys)
	c.SecureJSON(http.StatusOK, res)
	return
}

// 这个要验证下动态码

func (a *ApiService) unbindingGoogle(c *gin.Context) {

	var userCode types.UserCodeInfos

	err := c.BindJSON(&userCode)
	if err != nil {
		logrus.Info("not valid json", err)

		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}
	uid := userCode.Uid
	code := userCode.Code //验证动态码

	_, secret := db.QuerySecret(a.dbEngine, uid)

	codeint, err := strconv.ParseInt(code, 10, 64)

	if err != nil {
		logrus.Info("not valid code", code)

		res := util.ResponseMsg(-1, "fail", "google code is not pass,so can not unbinding google")
		c.SecureJSON(http.StatusOK, res)
		return
	}

	isTrue := VerifyCode(secret.Secret, int32(codeint))

	if !isTrue {
		res := util.ResponseMsg(-1, "fail", "google code is not pass,so can not unbinding google")
		c.SecureJSON(http.StatusOK, res)
	}
	logrus.Info("code pass verify,next unbind google")

	//下面才可以解绑--将db更新即可
	user, err := db.GetUser(a.dbEngine, uid)
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}

	if user == nil {
		res := util.ResponseMsg(-1, "fail", "apiKey is not exist")
		c.SecureJSON(http.StatusOK, res)
		return
	}

	user.IsBindGoogle = false

	err = db.UpdateUser(a.dbEngine, uid, user)

	if err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}
	res := util.ResponseMsg(0, "unbinding success", nil)
	c.SecureJSON(http.StatusOK, res)
	return
}

func (a *ApiService) userRevenueRanking(c *gin.Context) {
	total := c.Query("total")
	err, Revenue := db.UserRevenue(a.dbEngine, total)
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}
	// 总收益排行
	var UserRevenueList []interface{}
	for i := 0; i < len(Revenue); i++ {
		UserRevenue := make(map[string]interface{})
		err, user := db.QuerySecret(a.dbEngine, Revenue[i].Id)
		if err != nil {
			res := util.ResponseMsg(-1, "fail", err)
			c.SecureJSON(http.StatusOK, res)
			return
		}
		UserRevenue["placed"] = i + 1
		UserRevenue["username"] = user.UserName
		UserRevenue["revenue"] = Revenue[i].TotalBenefit
		UserRevenueList = append(UserRevenueList, UserRevenue)
	}
	// 总收益率排行
	// 总收益
	err, AllRevenue := db.UserAllRevenue(a.dbEngine)
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}
	// 总投资
	err, AllInvest := db.UserAllInvest(a.dbEngine)
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err)
		c.SecureJSON(http.StatusOK, res)
		return
	}

	body := make(map[string]interface{})
	body["revenueAmount"] = UserRevenueList
	res := util.ResponseMsg(0, "success", nil)
	c.SecureJSON(http.StatusOK, res)
	return
}
