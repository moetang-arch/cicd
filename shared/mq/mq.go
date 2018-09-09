package mq

import (
	"fmt"

	"github.com/nsqio/go-nsq"
)

type providerType int

const (
	TYPE_NSQ providerType = 1
)

type MqSender interface {
	Close() error
	Publish(topic string, data []byte) error
}

type nsqSenderImpl struct {
	producer *nsq.Producer
}

func (this *nsqSenderImpl) Publish(topic string, data []byte) error {
	return this.producer.Publish(topic, data)
}

func (this *nsqSenderImpl) Close() error {
	this.producer.Stop()
	return nil
}

type MqReceiver interface {
	Close() error
}

type nsqReceiverImpl struct {
	consumer *nsq.Consumer
}

func (this *nsqReceiverImpl) Close() error {
	this.consumer.Stop()
	return nil
}

type nsqHandler struct {
	ReceiveHandler
}

func (this nsqHandler) HandleMessage(message *nsq.Message) error {
	return this.ReceiveHandler(message.Body)
}

type ReceiveHandler func(data []byte) error

func NewMqSender(provider providerType, addr string) (MqSender, error) {
	switch provider {
	case TYPE_NSQ:
		config := nsq.NewConfig()
		producer, err := nsq.NewProducer(addr, config)
		if err != nil {
			return nil, err
		}

		impl := new(nsqSenderImpl)
		impl.producer = producer
		return impl, nil
	}
	panic(fmt.Sprint("unsupported provider:", provider))
}

func NewMqReceiver(provider providerType, addr string, topic, queue string, handler ReceiveHandler) (MqReceiver, error) {
	switch provider {
	case TYPE_NSQ:
		config := nsq.NewConfig()
		consumer, err := nsq.NewConsumer(topic, queue, config)
		if err != nil {
			return nil, err
		}
		consumer.AddHandler(nsqHandler{handler})
		err = consumer.ConnectToNSQD(addr)
		if err != nil {
			return nil, err
		}
		return &nsqReceiverImpl{consumer: consumer}, nil
	}
	panic(fmt.Sprint("unsupported provider:", provider))
}
