package log

import (
	"github.com/sirupsen/logrus"
)

// AddKafkaHook ...
func AddKafkaHook(topic string, brokers []string, level int) error {
	hook, err := newKafkaLogrusHook("group id not set yet!",
		getHookLevel(level),
		&logrus.JSONFormatter{
			DisableTimestamp: false,
			CallerPrettyfier: callerPrettyfier,
		},
		brokers,
		topic,
		nil)
	if err != nil {
		return err
	}

	logrus.AddHook(hook)
	return nil
}
