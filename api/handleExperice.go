package api

import (
	"fmt"
	"github.com/ethereum/api-in/db"
	"github.com/ethereum/api-in/types"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (a *ApiService) checkoutQualification(c *gin.Context) {
	uid := c.Query("uid")
	fmt.Println(uid)
	// 查询三个条件-查询数据库
	user := db.GetUser(a.dbEngine, uid)

	if user.IsBindGoogle == false {
		logrus.Info("条件不满足，google未绑定", user.IsBindGoogle)
	}
	if user.IsApiBind == false {
		logrus.Info("条件不满足，API未绑定", user.IsApiBind)
	}
	if user.InviteNumber < 1 {
		logrus.Info("条件不满足，当前未邀请人", user.InviteNumber)
	}

	res := types.HttpRes{}
	res.Code = 1
	res.Message = "当前符合领取体验金资格"
	res.Data = ""
	c.SecureJSON(http.StatusOK, res)
	return
}

func (a *ApiService) getExperienceFund(c *gin.Context) {
	uid := c.Query("uid")
	fmt.Println(uid)

	res := types.HttpRes{}

	// 查询三个条件-查询数据库
	user := db.GetUser(a.dbEngine, uid)

	//这里首先也要检查资格
	if user.IsBindGoogle == false {
		logrus.Info("条件不满足，google未绑定", user.IsBindGoogle)
		res.Code = -1
		res.Message = "领取体验金失败"
		c.SecureJSON(http.StatusOK, res)
	}
	if user.IsApiBind == false {
		logrus.Info("条件不满足，API未绑定", user.IsApiBind)
		res.Code = -1
		res.Message = "领取体验金失败"
		c.SecureJSON(http.StatusOK, res)
	}
	if user.InviteNumber < 1 {
		logrus.Info("条件不满足，当前未邀请人", user.InviteNumber)
		res.Code = -1
		res.Message = "领取体验金失败"
		c.SecureJSON(http.StatusOK, res)
	}
	logrus.Info("条件符合。可以领取体验金", uid)

	session := a.dbEngine.NewSession()
	err := session.Begin()
	if err != nil {
		return
	}
	//首先用户体验金增加一条记录
	userExperience := types.UserExperience{}
	_, err = session.Insert(userExperience)
	if err != nil {
		err := session.Rollback()
		if err != nil {
			return
		}
		logrus.Fatal(err)
	}
	//平台体验金资金池减少相应的数额
	TotalRevenueInfo := db.GetTotalRevenue(a.dbEngine)

	//总金额减少
	TotalRevenueInfo.TotalSum = TotalRevenueInfo.TotalSum - TotalRevenueInfo.PerSum
	//接收人数+1
	TotalRevenueInfo.ReceivePersons = TotalRevenueInfo.ReceivePersons + 1
	//更新
	_, err = session.Update(TotalRevenueInfo)
	if err != nil {
		err := session.Rollback()
		if err != nil {
			return
		}
		logrus.Fatal(err)
	}

	err = session.Commit()
	if err != nil {
		logrus.Fatal(err)
	}

	res.Code = 0
	res.Message = "领取体验金成功"
	c.SecureJSON(http.StatusOK, res)
	return
}

func (a *ApiService) getUserExperience(c *gin.Context) {
	uid := c.Query("uid")
	fmt.Println(uid)

	res := types.HttpRes{}

	userExperience := db.GetUserExperience(a.dbEngine, uid)

	res.Code = 0
	res.Message = "获取用户体验金信息成功"
	res.Data = userExperience

	c.SecureJSON(http.StatusOK, res)
	return
}
func (a *ApiService) getPlatformExperience(c *gin.Context) {
	res := types.HttpRes{}

	platformExperience := db.GetTotalRevenue(a.dbEngine)

	res.Code = 0
	res.Message = "获取平台体验金信息成功"
	res.Data = platformExperience

	c.SecureJSON(http.StatusOK, res)
	return
}
