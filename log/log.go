package log

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"

	"github.com/ethereum/api-in/config"

	"github.com/sirupsen/logrus"
)

var (
	enableDefaultFieldMap = false
	defaultFieldMap       = make(map[string]string)
)

func callerPrettyfier(f *runtime.Frame) (string, string) {
	fileName := fmt.Sprintf("%s:%d", f.File, f.Line)
	funcName := f.Function
	list := strings.Split(funcName, "/")
	if len(list) > 0 {
		funcName = list[len(list)-1]
	}
	return funcName, fileName
}

// for stdout
func callerFormatter(f *runtime.Frame) string {
	funcName, fileName := callerPrettyfier(f)
	return " @" + funcName + " " + fileName
}

func init() {
	logrus.SetReportCaller(true)
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetOutput(os.Stdout)
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: false,
		CallerPrettyfier: callerPrettyfier,
	})
}

// AddField add log fields
func AddField(key, value string) {
	if len(key) == 0 {
		return
	}
	if len(value) == 0 {
		return
	}
	enableDefaultFieldMap = true
	defaultFieldMap[key] = value
}

// DisableDefaultConsole 取消默认的控制台输出
func DisableDefaultConsole() {
	logrus.SetOutput(ioutil.Discard)
}

func getHookLevel(level int) []logrus.Level {
	if level < 0 || level > 5 {
		level = 5
	}
	return logrus.AllLevels[:level+1]
}

func Init(name string, config *config.Config) error {
	if config.Log.Stdout.Enable {
		AddConsoleOut(config.Log.Stdout.Level)
	}

	/*if config.Log.File.Enable {
		err := AddFileOut(config.Log.File.Path, config.Log.File.Level, 5)
		if err != nil {
			return err
		}
	}*/

	if config.Log.Kafka.Enable {
		err := AddKafkaHook(config.Log.Kafka.Topic, config.Log.Kafka.Brokers, config.Log.Kafka.Level)
		if err != nil {
			return err
		}

	}

	AddField("app", name)
	AddField("env_name", "prod")

	return nil
}
