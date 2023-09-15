package services

import (
	"github.com/sirupsen/logrus"
)

type UserBonusService struct {
}

func NewUserBonusService() *UserBonusService {
	return &UserBonusService{}
}

func (c *UserBonusService) Run() error {
	logrus.Info("UserBonusService run.........")

	logrus.Info("UserBonusService done.........")
	return nil
}

func (c *UserBonusService) Name() string {
	return "BonusStats"
}
