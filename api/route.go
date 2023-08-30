package api

import (
	"fmt"
	"github.com/ethereum/BGService/config"
	"github.com/ethereum/BGService/db"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"github.com/sirupsen/logrus"
	"net/http"
)

type ApiService struct {
	dbEngine    *xorm.Engine
	config      *config.Config
	RedisEngine db.CustomizedRedis
}

func NewApiService(dbEngine *xorm.Engine, RedisEngine db.CustomizedRedis, cfg *config.Config) *ApiService {
	return &ApiService{
		dbEngine:    dbEngine,
		config:      cfg,
		RedisEngine: RedisEngine,
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
	v1 := r.Group("/api/auth")
	{
		////下单
		//v1.POST("/order", a.order)
		////增加一条记录到users中
		//v1.POST("/enroll", a.enroll)
		v1.GET("/email", a.email)
		v1.POST("/register", a.register)
		v1.GET("/generateSecret", a.generateSecret)
		v1.GET("/verifyCode", a.verifyCode)
	}

	v6 := r.Group("/api/experienceActivity")
	{
		//检查领取体验金资格
		v6.GET("/checkoutQualification", a.checkoutQualification)
		//领取体验金
		v6.GET("/getExperienceFund", a.getExperienceFund)

		//获得用户的体验金收益率
		v6.GET("/getUserExperience", a.getUserExperience)
		//获取平台的体验金收益率
		v6.GET("/getPlatformExperience", a.getPlatformExperience)
	}
	logrus.Info("BGService un at " + a.config.Server.Port)

	err := r.Run(fmt.Sprintf(":%s", a.config.Server.Port))
	if err != nil {
		logrus.Fatalf("start http server err:%v", err)
	}
}
