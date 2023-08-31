package api

import (
	"fmt"
	"github.com/ethereum/BGService/db"
	"github.com/ethereum/BGService/util"
	"github.com/gin-gonic/gin"
	"net/http"
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
