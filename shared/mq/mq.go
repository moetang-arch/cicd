package mq

import "github.com/nsqio/go-nsq"

type providerType int

const (
	TYPE_NSQ providerType = 1
)

type Mq interface {
	Close() error
	Publish(topic string, data []byte) error
}

type nsqImpl struct {
	producer *nsq.Producer
}

func (this *nsqImpl) Publish(topic string, data []byte) error {
	return this.producer.Publish(topic, data)
}

func (this *nsqImpl) Close() error {
	this.producer.Stop()
	return nil
}

func NewMqSender(provider providerType, addr string) (Mq, error) {
	switch provider {
	case TYPE_NSQ:
		config := nsq.NewConfig()
		producer, err := nsq.NewProducer(addr, config)
		if err != nil {
			return nil, err
		}

		impl := new(nsqImpl)
		impl.producer = producer
		return impl, nil
	}
	panic("")
}
