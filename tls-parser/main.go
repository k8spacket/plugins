package main

import (
	"github.com/k8spacket/plugin-api"
	"github.com/k8spacket/plugins/tls-parser/log"
	"github.com/k8spacket/plugins/tls-parser/metrics"
	"github.com/k8spacket/plugins/tls-parser/metrics/connections"
	"github.com/k8spacket/tls-api/model"
)

type stream plugin_api.ReassembledStream

func (s stream) InitPlugin(manager plugin_api.PluginManager) {
	tls_parser_log.BuildLogger()
	manager.RegisterPlugin(s)
	manager.RegisterHttpHandler("/tlsparser/connections/", connections.TLSConnectionHandler)
	manager.RegisterHttpHandler("/tlsparser/api/data/", connections.TLSParserConnectionsHandler)
}

func (s stream) DistributeReassembledStream(reassembledStream plugin_api.ReassembledStream) {
	metrics.StoreStreamMetrics(reassembledStream)
}

func (s stream) DistributeTCPPacketPayload(tcpPacketPayload plugin_api.TCPPacketPayload) {
	if len(tcpPacketPayload.Payload) > 5 && tcpPacketPayload.Payload[0] == model.TLSRecord {
		metrics.CollectTCPPacketPayload(tcpPacketPayload.StreamId, tcpPacketPayload.Payload)
	}
}

func init() {}

// exported
var StreamPlugin stream
