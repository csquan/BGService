package router

import (
	"github.com/gin-gonic/gin"
	"user/pkg/ctrl"
	"user/pkg/web/auth"
)

type UserRouter struct {
	userCtrl ctrl.UserCtrl
}

func NewUserRouter(userCtrl ctrl.UserCtrl) UserRouter {
	return UserRouter{userCtrl: userCtrl}
}

func (ur *UserRouter) Route(rg *gin.RouterGroup) {
	r := rg.Group("user")
	r.Use(auth.MustExtractUser())
	r.POST("/genGa", ur.userCtrl.GenGa)
	r.POST("/bindGa", ur.userCtrl.BindGa)
	r.POST("/verifyGa", ur.userCtrl.VerifyGa)
	r.POST("/resEmail", ur.userCtrl.ResEmail)
	r.POST("/resMobile", ur.userCtrl.ResMobile)
	r.POST("/revEmail", ur.userCtrl.RevEmail)
	r.POST("/revMobile", ur.userCtrl.RevMobile)
	r.POST("/myEmail", ur.userCtrl.MyEmail)
	r.POST("/myMobile", ur.userCtrl.MyMobile)
	r.POST("/myNewEmail", ur.userCtrl.MyNewEmail)
	r.POST("/myNewMobile", ur.userCtrl.MyNewMobile)
	r.POST("/changeEmail", ur.userCtrl.ChangeEmail)
	r.POST("/changeMobile", ur.userCtrl.ChangeMobile)
	r.POST("/bindEmailBy", ur.userCtrl.BindEmailBy)
	r.POST("/bindMobileBy", ur.userCtrl.BindMobileBy)
	r.POST("/unbindGa", ur.userCtrl.UnbindGa)
	r.POST("/do-reset-password", ur.userCtrl.DoResetPassword)
	r.POST("/change-nick", ur.userCtrl.ChangeNick)

	r.POST("/query-free-nick-by-id", ur.userCtrl.QueryFreeNickByID)
	r.POST("/add-firm-user-by-id", ur.userCtrl.AddFirmUserByID)
	r.POST("/del-firm-user-by-id", ur.userCtrl.DelFirmUserByID)
	r.GET("/list-firm-user", ur.userCtrl.ListFirmUser)
	r.POST("/validate-n-sign", ur.userCtrl.ValidateAndSign)
}
