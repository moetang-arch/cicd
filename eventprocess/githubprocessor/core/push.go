package core

import (
	"github.com/moetang-arch/cicd/shared"
	"github.com/moetang-arch/cicd/shared/model"
	"github.com/moetang-arch/cicd/shared/mq"
)

type PushService struct {
	sender mq.MqSender
}

func NewPushService() *PushService {
	ps := new(PushService)
	sender, err := mq.NewMqSender(mq.TYPE_NSQ, "127.0.0.1:4150")
	if err != nil {
		panic("init mq error. " + err.Error())
	}
	ps.sender = sender
	return ps
}

func (this *PushService) SendPushEvent(pushEvent *PushEvent) error {
	// 1. convert to generic model
	pe := new(model.PushEvent)
	pe.PushSource = "github"
	pe.SourceId = 0 //FIXME need obtain from db
	pe.Branch = pushEvent.Ref
	pe.CommitVersion = pushEvent.AfterCommit
	pe.CloneSource = pushEvent.Repository.CloneUrl

	// 2. send to event bus for processing
	data, err := pe.ToBytes()
	if err != nil {
		return err
	}
	err = this.sender.Publish(shared.PUSH_EVENT_TOPIC, data)
	if err != nil {
		return err
	}
	return nil
}
