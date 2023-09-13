package services

import (
	"github.com/sirupsen/logrus"
)

type UserTxRecordService struct {
}

func NewUserTxRecordService() *UserTxRecordService {
	return &UserTxRecordService{}
}

func (c *UserTxRecordService) Run() error {
	logrus.Info("UserTxRecordService run.........")

	logrus.Info("UserTxRecordService done.........")
	return nil
}

func (c *UserTxRecordService) Name() string {
	return "UserTxRecordService"
}
