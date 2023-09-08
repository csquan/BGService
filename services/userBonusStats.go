package services

import (
	"github.com/ethereum/BGService/config"
	"github.com/go-xorm/xorm"
	"github.com/sirupsen/logrus"
)

type UserBonusService struct {
	engine *xorm.Engine
	config *config.Config
}

func NewUserBonusService(c *config.Config, engine *xorm.Engine) *UserBonusService {
	return &UserBonusService{
		engine: engine,
		config: c,
	}
}

func (c *UserBonusService) Run() error {
	logrus.Info("UserBonusService run.........")

	logrus.Info("UserBonusService done.........")
	return nil
}

func (c *UserBonusService) Name() string {
	return "BonusStats"
}
