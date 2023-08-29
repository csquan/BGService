package api

import (
	"github.com/ethereum/api-in/db"
	"github.com/ethereum/api-in/types"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"net/http"
)

//func (a *ApiService) init(c *gin.Context) {
//	buf := make([]byte, 2048)
//	n, _ := c.Request.Body.Read(buf)
//	data1 := string(buf[0:n])
//	res := types.HttpRes{}
//
//	isValid := gjson.Valid(data1)
//	if isValid == false {
//		logrus.Error("Not valid json")
//		res.Code = http.StatusBadRequest
//		res.Message = "Not valid json"
//		c.SecureJSON(http.StatusBadRequest, res)
//		return
//	}
//	name := gjson.Get(data1, "name")
//	apiKey := gjson.Get(data1, "apiKey")
//	apiSecret := gjson.Get(data1, "apiSecret")
//
//	mechanismData := types.Mechanism{
//		Name:      name.String(),
//		ApiKey:    apiKey.String(),
//		ApiSecret: apiSecret.String(),
//	}
//
//	err := a.db.CommitWithSession(a.db, func(s *xorm.Session) error {
//		if err := a.db.InsertMechanism(s, &mechanismData); err != nil {
//			logrus.Errorf("insert  InsertMechanism task error:%v tasks:[%v]", err, mechanismData)
//			return err
//		}
//		return nil
//	})
//	if err != nil {
//		logrus.Error(err)
//	}
//
//	res.Code = 0
//	res.Message = err.Error()
//	res.Data = ""
//
//	c.SecureJSON(http.StatusOK, res)
//}

func (a *ApiService) order(c *gin.Context) {
	buf := make([]byte, 1024)
	n, _ := c.Request.Body.Read(buf)
	data1 := string(buf[0:n])
	res := types.HttpRes{}

	isValid := gjson.Valid(data1)
	if isValid == false {
		logrus.Error("Not valid json")
		res.Code = http.StatusBadRequest
		res.Message = "Not valid json"
		c.SecureJSON(http.StatusBadRequest, res)
		return
	}
	//下面将信息存入db
	res.Code = 0
	res.Message = "success"
	res.Data = "null"

	c.SecureJSON(http.StatusOK, res)
}

func (a *ApiService) enroll(c *gin.Context) {
	buf := make([]byte, 1024)
	n, _ := c.Request.Body.Read(buf)
	data1 := string(buf[0:n])
	res := types.HttpRes{}

	isValid := gjson.Valid(data1)
	if isValid == false {
		logrus.Error("Not valid json")
		res.Code = http.StatusBadRequest
		res.Message = "Not valid json"
		c.SecureJSON(http.StatusBadRequest, res)
		return
	}
	uid := gjson.Get(data1, "uid")
	password := gjson.Get(data1, "password")

	user := types.Users{
		Uid:      uid.String(),
		Password: password.String(),
	}

	db.InsertUser(a.dbEngine, &user)
	//下面将信息存入db
	res.Code = 0
	res.Message = "success"
	res.Data = "null"

	c.SecureJSON(http.StatusOK, res)
}
