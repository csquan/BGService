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
	UserBonusService := NewUserBonusService()
	UserBenefit := NewUserBenefitService()
	UserTxRecordService := NewUserTxRecordService()

	t.services = []types.IAsyncService{
		UserBonusService,
		UserBenefit,
		UserTxRecordService,
	}

	for _, s := range t.services {
		s.Run()
	}
}
