package msg

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"time"
	"user/pkg/log"
)

var logger = log.C.Logger().With().Str("base", "kafka").Logger()

func NewKafkaProducer() *kafka.Producer {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "kafka:9092"})
	if err != nil {
		logger.Err(err).Msg("Failed to create kafka producer")
		panic(err)
	}
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				// The message delivery report, indicating success or
				// permanent failure after retries have been exhausted.
				// Application level retries won't help since the client
				// is already configured to do that.
				m := ev
				if m.TopicPartition.Error != nil {
					logger.Error().AnErr("Delivery failed: ", m.TopicPartition.Error).Send()
				} else {
					logger.Printf("Delivered message to topic %s [%d] at offset %v\n",
						*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
				}
			case kafka.Error:
				// Generic client instance-level errors, such as
				// broker connection failures, authentication issues, etc.
				//
				// These errors should generally be considered informational
				// as the underlying client will automatically try to
				// recover from any errors encountered, the application
				// does not need to take action on them.
				logger.Err(ev).Send()
			default:
				logger.Info().Interface("Ignored event: ", ev).Send()
			}
		}
	}()
	logger.Info().Interface("Created Producer: ", p).Send()
	return p
}

func Send(p *kafka.Producer, topic, key string, value []byte) {
	err := p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            []byte(key),
		Value:          value,
		Headers:        []kafka.Header{},
	}, nil)

	if err != nil {
		if err.(kafka.Error).Code() == kafka.ErrQueueFull {
			// Producer queue is full, wait 1s for messages
			// to be delivered then try again.
			time.Sleep(time.Second)
			Send(p, topic, key, value)
		}
		logger.Error().AnErr("Failed to produce message: ", err).Send()
	}
	// Flush and close the producer and the events channel
	//for p.Flush(10000) > 0 {
	//	fmt.Print("Still waiting to flush outstanding messages\n", err)
	//}
}
