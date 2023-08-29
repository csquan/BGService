package router

import (
	"github.com/gin-gonic/gin"
	"user/pkg/ctrl"
	"user/pkg/web/auth"
)

type PubRouter struct {
	pubCtrl ctrl.PubCtrl
}

func NewPubRouter(pubCtrl ctrl.PubCtrl) PubRouter {
	return PubRouter{pubCtrl: pubCtrl}
}

func (pr *PubRouter) Route(rg *gin.RouterGroup) {
	r := rg.Group("pub")
	r.Use(auth.SilentExtractUser())
	r.POST("/myEmail", pr.pubCtrl.MyEmail)
	r.POST("/myMobile", pr.pubCtrl.MyMobile)
	r.POST("/verifyEmail", pr.pubCtrl.VerifyEmail)
	r.POST("/verifyMobile", pr.pubCtrl.VerifyMobile)
	r.POST("/signIn", pr.pubCtrl.SignIn)
	r.GET("/refresh", pr.pubCtrl.Refresh)
	r.POST("/forget-e", pr.pubCtrl.ForgetE)
	r.POST("/forget-m", pr.pubCtrl.ForgetM)
	r.POST("/forget-send-email-by-email", pr.pubCtrl.ForgetSendEmailByEmail)
	r.POST("/forget-send-sms-by-email", pr.pubCtrl.ForgetSendSMSByEmail)
	r.POST("/forget-send-email-by-mobile", pr.pubCtrl.ForgetSendEmailByMobile)
	r.POST("/forget-send-sms-by-mobile", pr.pubCtrl.ForgetSendSMSByMobile)
	r.POST("/preset-password", pr.pubCtrl.PresetPassword)
	r.POST("/kyc-ok", pr.pubCtrl.KYC)
	r.POST("/kyc-query", pr.pubCtrl.KycQuery)
	r.POST("/validate-email", pr.pubCtrl.ValidateEmail)
	r.POST("/validate-mobile", pr.pubCtrl.ValidateMobile)
	r.POST("/signOut", pr.pubCtrl.SignOut)
	r.GET("/pubKey", pr.pubCtrl.PubKey)
	r.POST("/verifyGa", pr.pubCtrl.VerifyGa)
	r.POST("/resEmail", pr.pubCtrl.ResEmail)
	r.POST("/resMobile", pr.pubCtrl.ResMobile)
	r.POST("/revEmail", pr.pubCtrl.RevEmail)
	r.POST("/revMobile", pr.pubCtrl.RevMobile)
	r.POST("/signUp", pr.pubCtrl.SignUp)

	r.POST("/i-q-user-by-addr", pr.pubCtrl.QueryUserByAddr)
	r.POST("/kyc-set-user-status", pr.pubCtrl.KycSetUserStatus)
	r.POST("/kyc-user-list", pr.pubCtrl.KycUserList)
	r.POST("/i-q-user-multi", pr.pubCtrl.InnerUserQueryMulti)
}
