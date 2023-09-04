package api

import (
	"github.com/ethereum/BGService/db"
	"github.com/ethereum/BGService/util"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func handleParam(c *gin.Context) (int, int, string, error) {
	pageSize, ok := c.GetQuery("pageSize")
	if !ok {
		logrus.Error("pageSize not exist.")
		return 0, 0, "", errors.New("pageSize not exist.")
	}
	pageIndex, ok := c.GetQuery("pageIndex")
	if !ok {
		logrus.Error("pageIndex not exist.")
		return 0, 0, "", errors.New("pageIndex not exist.")
	}
	Type, ok := c.GetQuery("type")
	if !ok {
		logrus.Error("type not exist.")

		return 0, 0, "", errors.New("type not exist.")
	}
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		logrus.Error(err)
		return 0, 0, "", err
	}
	pageIndexInt, err := strconv.Atoi(pageIndex)
	if err != nil {
		logrus.Error(err)
		return 0, 0, "", err
	}
	return pageSizeInt, pageIndexInt, Type, nil
}

func (a *ApiService) list(c *gin.Context) {
	pageSizeInt, pageIndexInt, Type, err := handleParam(c)
	if err != nil {
		logrus.Error(err.Error())
		res := util.ResponseMsg(-1, "fail", err.Error())
		c.SecureJSON(http.StatusOK, res)
		return
	}
	msg, err := db.GetMsg(a.dbEngine, pageSizeInt, pageIndexInt, Type)
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err.Error())
		c.SecureJSON(http.StatusOK, res)
		return
	}
	var msgInfoList []interface{}
	for _, value := range msg {
		msgInfo := make(map[string]interface{})
		msgInfo["id"] = value.ID
		msgInfo["title"] = value.Title
		msgInfo["content"] = value.Content
		msgInfo["image"] = value.Cover
		msgInfo["datetime"] = value.CreateTime
		msgInfoList = append(msgInfoList, msgInfo)
	}
	body := make(map[string]interface{})
	body["total"] = len(msg)
	body["list"] = msgInfoList
	res := util.ResponseMsg(1, "success", body)
	c.SecureJSON(http.StatusOK, res)
	return
}

func (a *ApiService) hotSpotList(c *gin.Context) {
	Total, ok := c.GetQuery("total")
	if !ok {
		logrus.Error("total not exist.")
		res := util.ResponseMsg(-1, "fail", "total not exist.")
		c.SecureJSON(http.StatusOK, res)
		return
	}
	Type, ok := c.GetQuery("type")
	if !ok {
		logrus.Error("type not exist.")
		res := util.ResponseMsg(-1, "fail", "type not exist.")
		c.SecureJSON(http.StatusOK, res)
		return
	}

	msg, err := db.GetMsg(a.dbEngine, pageSizeInt, pageIndexInt, Type)
	if err != nil {
		res := util.ResponseMsg(-1, "fail", err.Error())
		c.SecureJSON(http.StatusOK, res)
		return
	}
	var msgInfoList []interface{}
	for _, value := range msg {
		msgInfo := make(map[string]interface{})
		msgInfo["id"] = value.ID
		msgInfo["title"] = value.Title
		msgInfo["content"] = value.Content
		msgInfo["image"] = value.Cover
		msgInfo["datetime"] = value.CreateTime
		msgInfoList = append(msgInfoList, msgInfo)
	}
	body := make(map[string]interface{})
	body["total"] = len(msg)
	body["list"] = msgInfoList
	res := util.ResponseMsg(1, "success", body)
	c.SecureJSON(http.StatusOK, res)
	return
}
