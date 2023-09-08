package services

import (
	"github.com/ethereum/BGService/config"
	"github.com/go-xorm/xorm"
	"github.com/sirupsen/logrus"
)

type UserTxRecordService struct {
	engine *xorm.Engine
	config *config.Config
}

func NewUserTxRecordService(c *config.Config, engine *xorm.Engine) *UserTxRecordService {
	return &UserTxRecordService{
		engine: engine,
		config: c,
	}
}

func (c *UserTxRecordService) Run() error {
	logrus.Info("UserTxRecordService run.........")

	logrus.Info("UserTxRecordService done.........")
	return nil
}

func (c *UserTxRecordService) Name() string {
	return "UserTxRecordService"
}
