package core

import "fmt"

type PushService struct {
}

func NewPushService() *PushService {
	ps := new(PushService)
	return ps
}

func (this *PushService) SendPushEvent(pushEvent *PushEvent) error {
	//TODO 1. convert to generic model
	//TODO 2. send to event bus for processing
	fmt.Println(pushEvent)
	return nil
}
