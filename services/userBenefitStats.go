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

func (c *UserBenefitService) Run() error {
	logrus.Info("UserBenefitService run.........")

	logrus.Info("UserBenefitService done.........")
	return nil
}

func (c *UserBenefitService) Name() string {
	return "BurnStats"
}
