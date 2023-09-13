package services

import (
	"github.com/ethereum/BGService/config"
	"github.com/ethereum/BGService/types"
	"github.com/go-xorm/xorm"
)

type ServiceScheduler struct {
	conf     *config.Config
	engine   *xorm.Engine
	services []types.IAsyncService
}

func NewServiceScheduler() (t *ServiceScheduler, err error) {
	t = &ServiceScheduler{
		services: make([]types.IAsyncService, 0),
	}

	return
}

func (t *ServiceScheduler) Start() {
	UserTxRecordService := NewUserTxRecordService()
	UserBonusService := NewUserBonusService()
	UserBenefit := NewUserBenefitService()

	t.services = []types.IAsyncService{
		UserTxRecordService,
		UserBonusService,
		UserBenefit,
	}

	for _, s := range t.services {
		s.Run()
	}
}
