package wireset

import (
	"bank-system-go/internal/config"
	"bank-system-go/pkg/logger"
	"bank-system-go/pkg/mq"
	"bank-system-go/pkg/mq/kafka"
	"strings"

	"github.com/pkg/errors"
)

func InitMQ(config config.Config, logger logger.Logger) (mq.MQ, error) {
	switch strings.ToLower(config.MQ.Driver) {
	case "kafka":
		return kafka.NewKafkaMQ(kafka.KafkaOption{
			Brokers:        config.MQ.KafkaOption.Brokers,
			ConsumerGroup:  config.MQ.KafkaOption.ConsumerGroup,
			OffsetsInitial: config.MQ.KafkaOption.OffsetsInitial,
			LoggerAdapter:  logger.WrapedWatermillLogger(),
		})
	}
	return nil, errors.New("no supported driver [" + config.MQ.Driver + "]")
}
