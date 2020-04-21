package nats

import (
	broker "github.com/nats-io/nats.go"
)

type natsPubSub struct {
	nc *broker.Conn
}

func New(nc *broker.Conn) {
	return
}