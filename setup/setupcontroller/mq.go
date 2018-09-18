package main

import (
	"io"

	"github.com/moetang-arch/cicd/shared/mq"
	"github.com/moetang-arch/cicd/shared"
)

func startNsq(nsqAddr string, fn func(data []byte) error) (io.Closer, error) {
	r, err := mq.NewMqReceiver(mq.TYPE_NSQ, nsqAddr, shared.PUSH_EVENT_TOPIC, "setup_master", fn)
	return r, err
}
