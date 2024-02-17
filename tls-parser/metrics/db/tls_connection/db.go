package tls_connection_db

import (
	"github.com/k8spacket/plugins/idb"
	tls_parser_log "github.com/k8spacket/plugins/tls-parser/log"
	"github.com/k8spacket/plugins/tls-parser/metrics/model"
)

var db, _ = idb.StartDB[model.TLSConnection]("tls_connections")

func ReadAll() []model.TLSConnection {
	result, err := db.ReadAll()
	if err != nil {
		tls_parser_log.LOGGER.Printf("[db:tls_connections:ReadAll] Error: %+v", err)
		return []model.TLSConnection{}
	}
	return result
}

func Upsert(key string, value *model.TLSConnection) {
	err := db.Upsert(key, *value)
	if err != nil {
		tls_parser_log.LOGGER.Printf("[db:tls_connections:Upsert] Error: %+v", err)
	}
}
