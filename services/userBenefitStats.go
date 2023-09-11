package services

import (
	"github.com/ethereum/BGService/config"
	"github.com/go-xorm/xorm"
	"github.com/sirupsen/logrus"
)

type UserBenefitService struct {
	engine *xorm.Engine
	config *config.Config
}

func NewUserBenefitService(c *config.Config, engine *xorm.Engine) *UserBenefitService {
	return &UserBenefitService{
		engine: engine,
		config: c,
	}
}

// 这里取用户每日收益存储进db
func (c *UserBenefitService) Run() error {
	logrus.Info("UserBenefitService run.........")
	// 首先查询参与策略产品得用户UID，得到对应得APIKEY，根据策略属性
	// 现货：GET /sapi/v1/accountSnapshot 取每日资产快照 u本位合约  币本位合约--币

	logrus.Info("UserBenefitService done.........")
	return nil
}

func (c *UserBenefitService) Name() string {
	return "BurnStats"
}
