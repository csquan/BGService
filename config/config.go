package config

import (
	"fmt"
	"github.com/jialanli/windward"
	"github.com/spf13/viper"
	"os"
)

type stdout struct {
	Enable bool `mapstructure:"enable"`
	Level  int  `mapstructure:"level"`
}

type file struct {
	Enable bool   `mapstructure:"enable"`
	Path   string `mapstructure:"path"`
	Level  int    `mapstructure:"level"`
}

type kafka struct {
	Enable  bool     `mapstructure:"enable"`
	Level   int      `mapstructure:"level"`
	Brokers []string `mapstructure:"kafka_servers"`
	Topic   string   `mapstructure:"topic"`
}

type Config struct {
	Db struct {
		Name     string `yaml:"name"`
		Password string `yaml:"password"`
		Ip       string `yaml:"ip"`
		Port     int32  `yaml:"port"`
		Database string `yaml:"database"`
	}
	Server struct {
		Port string `yaml:"port"`
	}
	Log struct {
		Stdout stdout `mapstructure:"stdout"`
		File   file   `mapstructure:"file"`
		Kafka  kafka  `mapstructure:"kafka"`
	}
}

var Conf Config

func Readconfig(filename string) (*Config, error) {
	path, err := os.Getwd()
	if err != nil {
		fmt.Errorf("++++++err ++++++++++: %v", err)
		return nil, fmt.Errorf("err : %v", err)
	}
	//加载配置文件
	file := path + "/" + filename
	w := windward.GetWindward()
	w.InitConf([]string{file}) //初始化自定义的配置文件

	//获取数据库连接名密码等数据
	var config Config //定义结构体【注意：这里需要有两层结构，因为w.ReadConfig读取的是data以及data中的数据】

	err = w.ReadConfig(file, &config)
	if err != nil {
		fmt.Sprintln("初始化配置文件失败")
		return nil, err
	}
	return &config, nil
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
