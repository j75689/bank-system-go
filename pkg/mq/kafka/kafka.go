package kafka

import (
	"bank-system-go/pkg/mq"
	"context"
	"errors"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
)

type KafkaOption struct {
	Brokers        []string
	ConsumerGroup  string
	OffsetsInitial int64
	LoggerAdapter  watermill.LoggerAdapter
}

var _ mq.MQ = (*KafkaMQ)(nil)

func NewKafkaMQ(option KafkaOption) (*KafkaMQ, error) {
	publisher, err := kafka.NewPublisher(
		kafka.PublisherConfig{
			Brokers:   option.Brokers,
			Marshaler: kafka.DefaultMarshaler{},
		}, option.LoggerAdapter,
	)
	if err != nil {
		return nil, err
	}

	saramaSubscriberConfig := kafka.DefaultSaramaSubscriberConfig()
	saramaSubscriberConfig.Consumer.Offsets.Initial = option.OffsetsInitial
	subscriber, err := kafka.NewSubscriber(
		kafka.SubscriberConfig{
			Brokers:       option.Brokers,
			Unmarshaler:   kafka.DefaultMarshaler{},
			ConsumerGroup: option.ConsumerGroup,
		}, option.LoggerAdapter,
	)
	if err != nil {
		return nil, err
	}

	return &KafkaMQ{
		publisher:  publisher,
		subscriber: subscriber,
	}, nil
}

type KafkaMQ struct {
	publisher  message.Publisher
	subscriber message.Subscriber
}

func (mq *KafkaMQ) Publish(topic, key string, data []byte) error {
	if len(key) == 0 {
		key = watermill.NewUUID()
	}
	return mq.publisher.Publish(topic, message.NewMessage(key, message.Payload(data)))
}

func (mq *KafkaMQ) Subscribe(ctx context.Context, topic string, process func(key string, data []byte) (bool, error), errCallBack ...func(error)) error {
	if process == nil {
		return errors.New("process is nil function")
	}

	message, err := mq.subscriber.Subscribe(ctx, topic)
	if err != nil {
		return err
	}

	for m := range message {
		isAck, err := process(m.UUID, m.Payload)
		if err != nil {
			for _, cb := range errCallBack {
				cb(err)
			}
		}
		if isAck {
			m.Ack()
		} else {
			m.Nack()
		}
	}

	return nil
}

func (mq *KafkaMQ) Close() error {
	if err := mq.publisher.Close(); err != nil {
		return err
	}
	return mq.subscriber.Close()
}