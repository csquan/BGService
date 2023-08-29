package web

import (
	"encoding/hex"
	"github.com/ethereum/api-in/pkg/model"
	"github.com/ethereum/api-in/pkg/util"
	"github.com/ethereum/api-in/pkg/util/ecies"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	jsoniter "github.com/json-iterator/go"
	"net/http"
	"time"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func BadRes(c *gin.Context, err util.Err) {
	er := HttpMsg{
		Code:    err.Code(),
		Message: err.Msg(),
	}
	c.AbortWithStatusJSON(http.StatusOK, er)
}

// HttpMsg error msg
type HttpMsg struct {
	Code    int    `json:"code" example:"502"`
	Message string `json:"message" example:"可耻滴失败鸟"`
}

// HttpData success data
type HttpData[T any] struct {
	Code int `json:"code" example:"0"`
	Data T   `json:"data"`
}

type HttpRes[T any] struct {
	Code    int    `json:"code" example:"502"`
	Message string `json:"message" example:"可耻滴失败鸟"`
	Data    T      `json:"data"`
}

// GoodResp success response
func GoodResp[T any](c *gin.Context, data T) {
	ret := HttpData[T]{
		Code: 0,
		Data: data,
	}
	c.JSON(http.StatusOK, ret)
}

func GetTokenUser(c *gin.Context) *model.TokenUser {
	if value, exists := c.Get("currentUser"); !exists {
		return nil
	} else {
		if tu, ok := value.(model.TokenUser); ok {
			return &tu
		} else {
			return nil
		}
	}
}

type MyInnerCli struct {
	Cli *resty.Client
	Pub *ecies.PublicKey
}

func NewCli(baseUrl, pubStr string) (*MyInnerCli, util.Err) {
	cli := resty.New()
	cli.JSONMarshal = json.Marshal
	cli.JSONUnmarshal = json.Unmarshal
	//cli.SetRetryCount(3).
	//	SetRetryWaitTime(5 * time.Second).
	//	SetRetryMaxWaitTime(30 * time.Second).
	//	SetRetryAfter(func(client *resty.Client, resp *resty.Response) (time.Duration, error) {
	//		return 0, nil
	//	})
	cli.SetBaseURL(baseUrl)
	pubKey, err := ecies.PublicFromString(pubStr)
	if err != nil {
		return nil, err
	}
	return &MyInnerCli{Cli: cli, Pub: pubKey}, nil
}

func (cli *MyInnerCli) Post(url string, req map[string]interface{}, res interface{}) util.Err {
	nowStr := time.Now().UTC().Format(http.TimeFormat)
	ct, err := ecies.Encrypt(cli.Pub, []byte(nowStr))
	if err != nil {
		return err
	}
	req["verified"] = hex.EncodeToString(ct)

	resp, er := cli.Cli.R().SetBody(req).SetResult(&res).Post(url)
	if er != nil {
		return util.ErrResty
	}
	if resp.StatusCode() != http.StatusOK {
		return util.ErrResty
	}
	return nil
}
