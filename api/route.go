package api

import (
	"fmt"
	"github.com/ethereum/api-in/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"github.com/sirupsen/logrus"
	"net/http"
)

type ApiService struct {
	dbEngine *xorm.Engine
	config   *config.Config
}

func NewApiService(dbEngine *xorm.Engine, cfg *config.Config) *ApiService {
	return &ApiService{
		dbEngine: dbEngine,
		config:   cfg,
	}
}

func (a *ApiService) Run() {
	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowCredentials = true
	corsConfig.AllowOrigins = []string{"*"}
	r.Use(func(ctx *gin.Context) {
		method := ctx.Request.Method
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Headers", "*")
		// ctx.Header("Access-Control-Allow-Headers", "Content-Type,addr,GoogleAuth,AccessToken,X-CSRF-Token,Authorization,Token,token,auth,x-token")
		ctx.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		ctx.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		ctx.Header("Access-Control-Allow-Credentials", "true")
		if method == "OPTIONS" {
			ctx.AbortWithStatus(http.StatusNoContent)
			return
		}
		ctx.Next()
	})

	//验证token--先不验证token
	//r.Use(auth.MustExtractUser())

	//r.POST("/init", a.init)
	//下单
	r.POST("/order", a.order)
	//增加一条记录到users中
	r.POST("/enroll", a.enroll)

	r.GET("/query", a.query)

	logrus.Info("BGService un at " + a.config.Server.Port)

	err := r.Run(fmt.Sprintf(":%s", a.config.Server.Port))
	if err != nil {
		logrus.Fatalf("start http server err:%v", err)
	}
}
