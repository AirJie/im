package tcp_conn

import (
	"log"
	"sync"
)

var manager sync.Map

func load(deviceId int64) *TcpConn {
	value, ok := manager.Load(deviceId)
	if !ok {
		log.Fatal("Failed to load", ok)
		return nil
	}
	return value.(*TcpConn)
}

func store(deviceId int64, link *TcpConn) {
	manager.Store(deviceId, link)
}

func delete(deviceId int64) {
	manager.Delete(deviceId)
}
