package wireset

import (
	"bank-system-go/internal/config"
	"bank-system-go/pkg/logger"
	"bank-system-go/pkg/mq"
	"bank-system-go/pkg/mq/kafka"
	"strings"

	"github.com/pkg/errors"
)

func InitMQ(config config.Config, log logger.Logger) (queue mq.MQ, err error) {

	switch strings.ToLower(config.MQ.Driver) {
	case "kafka":
		queue, err = kafka.NewKafkaMQ(kafka.KafkaOption{
			Brokers:        config.MQ.KafkaOption.Brokers,
			ConsumerGroup:  config.MQ.KafkaOption.ConsumerGroup,
			OffsetsInitial: config.MQ.KafkaOption.OffsetsInitial,
			LoggerAdapter:  logger.WrapWatermillLogger(log.Logger),
		})
	default:
		err = errors.New("no supported driver [" + config.MQ.Driver + "]")
	}
	if queue != nil {
		queue.SubscriberMiddleware(func(key string, data []byte) {
			log.Info().Str("request_id", key).Bytes("message", data).Send()
		})
	}
	return
}
