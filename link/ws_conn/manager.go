package ws_conn

import (
	"log"
	"sync"
)

var manager sync.Map

func load(deviceId int64) *Linker {
	value, ok := manager.Load(deviceId)
	if !ok {
		log.Fatal("Failed to load", ok)
		return nil
	}
	return value.(*Linker)
}

func store(deviceId int64, link *Linker) {
	manager.Store(deviceId, link)
}

func delete(deviceId int64) {
	manager.Delete(deviceId)
}
