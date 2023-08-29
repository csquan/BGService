package log

import (
	"crypto/tls"
	"log"
	"time"

	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type kafkaLogrusHook struct {
	id           string
	defaultTopic string
	levels       []logrus.Level
	formatter    logrus.Formatter
	producer     sarama.AsyncProducer
}

func newKafkaLogrusHook(id string,
	levels []logrus.Level,
	formatter logrus.Formatter,
	brokers []string,
	defaultTopic string,
	tls *tls.Config) (*kafkaLogrusHook, error) {
	var err error
	var producer sarama.AsyncProducer
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.RequiredAcks = sarama.WaitForLocal            // Only wait for the leader to ack
	kafkaConfig.Producer.Compression = sarama.CompressionSnappy        // Compress messages
	kafkaConfig.Producer.Flush.Frequency = 500 * time.Millisecond      // Flush batches every 500ms
	kafkaConfig.Producer.Partitioner = sarama.NewRoundRobinPartitioner // one partition at one time

	// check here if provided *tls.Config is not nil and assign to the sarama config
	// NOTE: we automatically enabled the TLS config because sarama would error out if our
	//       config were non-nil but disabled. To avoid issue father down the stack, we enable.
	if tls != nil {
		kafkaConfig.Net.TLS.Enable = true
		kafkaConfig.Net.TLS.Config = tls
	}

	kafkaConfig.Net.DialTimeout = time.Second
	if producer, err = sarama.NewAsyncProducer(brokers, kafkaConfig); err != nil {
		return nil, err
	}

	go func() {
		for err := range producer.Errors() {
			log.Printf("Failed to send log entry to Kafka: %v\n", err)
		}
	}()

	hook := &kafkaLogrusHook{
		id,
		defaultTopic,
		levels,
		formatter,
		producer,
	}

	return hook, nil
}

func (hook *kafkaLogrusHook) ID() string {
	return hook.id
}

func (hook *kafkaLogrusHook) Levels() []logrus.Level {
	return hook.levels
}

func (hook *kafkaLogrusHook) Fire(entry *logrus.Entry) error {
	if enableDefaultFieldMap {
		for key, value := range defaultFieldMap {
			if _, ok := entry.Data[key]; !ok {
				entry.Data[key] = value
			}
		}
	}
	var partitionKey sarama.ByteEncoder
	var b []byte
	var err error

	t, _ := entry.Data["time"].(time.Time)
	if b, err = t.MarshalBinary(); err != nil {
		return err
	}
	partitionKey = sarama.ByteEncoder(b)

	if b, err = hook.formatter.Format(entry); err != nil {
		return err
	}
	value := sarama.ByteEncoder(b)

	topic := hook.defaultTopic
	if tsRaw, ok := entry.Data["topic"]; ok {
		ts, ok := tsRaw.(string)
		if !ok {
			return errors.New("incorrect topic filed type (should be string)")
		}
		topic = ts
	}
	hook.producer.Input() <- &sarama.ProducerMessage{
		Key:   partitionKey,
		Topic: topic,
		Value: value,
	}
	return nil
}
