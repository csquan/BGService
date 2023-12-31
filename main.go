package main

import (
	"flag"
	"fmt"
	"github.com/ethereum/BGService/api"
	"github.com/ethereum/BGService/config"
	"github.com/ethereum/BGService/db"
	"github.com/ethereum/BGService/log"
	"github.com/ethereum/BGService/services"
	"github.com/go-co-op/gocron"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

const CONTRACTLEN = 42

var (
	conffile string
	env      string
)

func init() {
	flag.StringVar(&conffile, "conf", "config.yaml", "conf file")
	flag.StringVar(&env, "env", "prod", "Deploy environment: [ prod | test ]. Default value: prod")
}

func main() {
	var err error
	if config.Conf, err = config.LoadConfig("./conf"); err != nil {
		logrus.Info("🚀 Could not load environment variables")
		return
	}

	flag.Parse()

	err = log.Init("BGService", &config.Conf)
	if err != nil {
		logrus.Info(err)
	}

	dbEngine := db.GetDBEngine(&config.Conf)
	RedisEngine := db.GetRedisEngine(&config.Conf)

	//setup scheduler
	scheduler, err := services.NewServiceScheduler()
	if err != nil {
		return
	}

	//activityBenefitService := services.NewActivityBenefitService()
	//activityBenefitService.Run()
	//UserBenefit := services.NewUserBenefitService()
	//UserBenefit.Run()

	timezone, _ := time.LoadLocation("Asia/Shanghai")
	s := gocron.NewScheduler(timezone)
	s.Every(1).Day().At(config.Conf.Schedule.Time).Do(scheduler.Start)

	apiService := api.NewApiService(dbEngine, RedisEngine, &config.Conf)
	go apiService.Run()

	s.StartBlocking()
	//listen kill signal
	closeCh := make(chan os.Signal, 1)

	for {
		select {
		case <-closeCh:
			fmt.Printf("receive os close sigal")
			return
		}
	}
}
