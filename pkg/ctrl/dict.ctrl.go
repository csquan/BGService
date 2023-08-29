package ctrl

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"user/pkg/db"
	"user/pkg/log"
	"user/pkg/model"
	"user/pkg/util"
	"user/pkg/web"
)

type DictCtrl struct {
	DB  *gorm.DB
	Log zerolog.Logger
}

func NewDictCtrl() DictCtrl {
	logger := log.C.Logger().With().Str("ctrl", "dict").Logger()
	return DictCtrl{DB: db.DB, Log: logger}
}

// GetDictByGrp godoc
//
//	@Summary		从字典取分组信息
//	@Description	直接把字典分组放在url上，取该组字典信息，属于公共接口，不设权限
//	@Tags			dict
//	@Accept			json
//	@Produce		json
//	@Param			grp	path		int	true	"group ID"
//	@Success		200	{object}	web.HttpData[[]model.DictView]
//	@Failure		502	{object}	web.HttpMsg
//	@Router			/dict/{grp} [get]
func (dc *DictCtrl) GetDictByGrp(c *gin.Context) {
	grp := c.Param("grp")

	var dict model.Dict
	var firmTypes []model.DictView
	err := dc.DB.Model(&dict).Find(&firmTypes, "\"group\"=?", grp).Error
	if err != nil {
		web.BadRes(c, util.ErrDB)
		dc.Log.Err(err).Send()
		return
	}
	web.GoodResp(c, firmTypes)
}
