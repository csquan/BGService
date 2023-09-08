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

func NewServiceScheduler(conf *config.Config, engine *xorm.Engine) (t *ServiceScheduler, err error) {
	t = &ServiceScheduler{
		conf:     conf,
		engine:   engine,
		services: make([]types.IAsyncService, 0),
	}

	return
}

func (t *ServiceScheduler) Start() {
	userBenefitService := NewUserBenefitService(t.conf, t.engine)

	t.services = []types.IAsyncService{
		userBenefitService,
	}

	for _, s := range t.services {
		s.Run()
	}
}
