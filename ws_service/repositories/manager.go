package repositories

import (
	"golang.org/x/net/websocket"
	"log"
	"sync"
)

var manager sync.Map



func load(deviceId int64) *Link {
	value, ok := manager.Load(deviceId)
	if !ok {
		log.Fatal("Failed to load", ok)
		return nil
	}
	return value.(*Link)
}

func store(deviceId int64, link *Link) {
	manager.Store(deviceId, link)
}

func delete(deviceId int64) {
	manager.Delete(deviceId)
}
