package main

import (
	"github.com/k8spacket/plugin-api/v2"
	"github.com/k8spacket/plugins/tls-parser/log"
	"github.com/k8spacket/plugins/tls-parser/metrics"
	"github.com/k8spacket/plugins/tls-parser/metrics/connections"
)

type event plugin_api.TLSEvent

func (e event) InitPlugin(manager plugin_api.PluginManager) {
	tls_parser_log.BuildLogger()

	manager.RegisterTLSPlugin(e)
	manager.RegisterHttpHandler("/tlsparser/connections/", connections.TLSConnectionHandler)
	manager.RegisterHttpHandler("/tlsparser/api/data", connections.TLSParserConnectionsHandler)
	manager.RegisterHttpHandler("/tlsparser/api/data/", connections.TLSParserConnectionDetailsHandler)
}

func (e event) DistributeTLSEvent(event plugin_api.TLSEvent) {
	metrics.StoreTLSMetrics(event)
}

func init() {}

// exported
var TLSConsumerPlugin event
