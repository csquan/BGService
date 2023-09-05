package api

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"

	// 导入session存储引擎
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

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// todo 登录校验
		//session := sessions.Default(c)
		//Uid := session.Get("Uid")
		//if Uid == nil {
		//	c.JSON(http.StatusUnauthorized, gin.H{
		//		"error": "Unauthorized",
		//	})
		//	c.Abort()
		//	return
		//}
		//// 用户已登录，将用户 ID 传递给后续的处理函数
		//c.Set("Uid", Uid)
		c.Set("Uid", "24670980929080")
		c.Set("invitationCode", "VCZ34Z71")
		c.Next()
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
	store := cookie.NewStore([]byte("secret123456"))
	//session中间件生效，参数mysession，是浏览器端cookie的名字
	r.Use(sessions.Sessions("mysession", store))

	//验证token--先不验证token
	//r.Use(auth.MustExtractUser())
	v1 := r.Group("/api/auth")
	{
		v1.GET("/email", a.email)
		v1.POST("/register", a.register)
		v1.POST("/login", a.login)
		v1.POST("/logout", authMiddleware(), a.logout)
		v1.POST("/forgotPassword", authMiddleware(), a.forgotPassword)
		v1.POST("/resetPassword", authMiddleware(), a.resetPassword)
		v1.GET("/generateSecret", authMiddleware(), a.generateSecret)
		v1.POST("/verifyCode", a.verifyCode)
	}
	v2 := r.Group("/api/user")
	{
		v2.GET("info", authMiddleware(), a.info)
		v2.GET("myApi", authMiddleware(), a.myApi)
		v2.POST("bindingApi", authMiddleware(), a.bindingApi)
		v2.GET("unbindingApi", authMiddleware(), a.unbindingApi)
		v2.GET("invite", authMiddleware(), a.invite)
		v2.GET("inviteRanking", authMiddleware(), a.inviteRanking)
		v2.GET("/getStrategy", a.getStrategy)
		v2.POST("/unbindingGoogle", a.unbindingGoogle)
	}

	v3 := r.Group("/api/info")
	{
		v3.GET("/list", a.list)
		v3.GET("/hotSpotList", a.hotSpotList)
		v3.GET("/details", a.details)
	}

	v4 := r.Group("/api/market")
	{
		//添加/移除自选
		v4.POST("/addConcern", a.addConcern)
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

	v7 := r.Group("/api/wallet")
	{
		v7.GET("/getTradeAccountDetail", a.getTradeAccountDetail)
	}

	v8 := r.Group("/api/product")
	{
		v8.GET("/overview", a.overview)
	}

	logrus.Info("BGService un at " + a.config.Server.Port)

	err := r.Run(fmt.Sprintf(":%s", a.config.Server.Port))
	if err != nil {
		logrus.Fatalf("start http server err:%v", err)
	}
}
