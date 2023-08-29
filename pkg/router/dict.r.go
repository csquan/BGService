package router

import (
	"github.com/gin-gonic/gin"
	"user/pkg/ctrl"
)

type DictRouter struct {
	dictCtrl ctrl.DictCtrl
}

func NewDictRoute(dictCtrl ctrl.DictCtrl) DictRouter {
	return DictRouter{dictCtrl}
}

func (dr *DictRouter) Route(rg *gin.RouterGroup) {
	router := rg.Group("dict")
	router.GET("/:grp", dr.dictCtrl.GetDictByGrp)
}
