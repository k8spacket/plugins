package tcp_connection_db

import (
	"github.com/k8spacket/plugins/idb"
	nodegraph_log "github.com/k8spacket/plugins/nodegraph/log"
	"github.com/k8spacket/plugins/nodegraph/metrics/nodegraph/model"
)

var db, _ = idb.StartDB[model.ConnectionItem]("tcp_connections")

func Read(key string) model.ConnectionItem {
	result, err := db.Read(key)
	if err != nil {
		// can happen, silent
		return model.ConnectionItem{}
	}
	return result
}

func ReadAll() []model.ConnectionItem {
	result, err := db.ReadAll()
	if err != nil {
		nodegraph_log.LOGGER.Printf("[db:tcp_connections:ReadAll] Error: %+v", err)
		return []model.ConnectionItem{}
	}
	return result
}

func Set(key string, value *model.ConnectionItem) {
	err := db.Upsert(key, *value)
	if err != nil {
		nodegraph_log.LOGGER.Printf("[db:tcp_connections:Upsert] Error: %+v", err)
	}
}
