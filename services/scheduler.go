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
	UserBenefit := NewUserBenefitService()
	UserTxRecordService := NewUserTxRecordService()
	activityBenefitService := NewActivityBenefitService()

	t.services = []types.IAsyncService{
		UserBenefit,
		UserTxRecordService,
		activityBenefitService,
	}

	for _, s := range t.services {
		s.Run()
	}
}
