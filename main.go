package main

import (
	"flag"
	"fmt"
	"github.com/ethereum/BGService/api"
	"github.com/ethereum/BGService/config"
	"github.com/ethereum/BGService/db"
	"github.com/ethereum/BGService/log"
	"github.com/ethereum/BGService/services"
	"github.com/jasonlvhit/gocron"
	"github.com/sirupsen/logrus"
	"os"
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
	if config.Conf, err = config.LoadConfig("./configTest"); err != nil {
		logrus.Error("🚀 Could not load environment variables")
		panic(err)
	}

	flag.Parse()

	err = log.Init("BGService", &config.Conf)
	if err != nil {
		log.Fatal(err)
	}

	dbEngine := db.GetDBEngine(&config.Conf)
	RedisEngine := db.GetRedisEngine(&config.Conf)

	//setup scheduler
	scheduler, err := services.NewServiceScheduler()
	if err != nil {
		return
	}

	gocron.Every(1).Day().At(config.Conf.Schedule.Time).Do(scheduler.Start)

	apiService := api.NewApiService(dbEngine, RedisEngine, &config.Conf)
	go apiService.Run()

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
