package main

import (
	"github.com/k8spacket/plugin-api"
	"github.com/k8spacket/plugins/nodegraph/log"
	"github.com/k8spacket/plugins/nodegraph/metrics"
	"github.com/k8spacket/plugins/nodegraph/metrics/nodegraph"
)

type stream plugin_api.ReassembledStream

func (s stream) InitPlugin(manager plugin_api.PluginManager) {
	nodegraph_log.BuildLogger()

	manager.RegisterPlugin(s)
	manager.RegisterHttpHandler("/nodegraph/connections", nodegraph.ConnectionHandler)
	manager.RegisterHttpHandler("/nodegraph/api/health", nodegraph.Health)
	manager.RegisterHttpHandler("/nodegraph/api/graph/fields", nodegraph.NodeGraphFieldsHandler)
	manager.RegisterHttpHandler("/nodegraph/api/graph/data", nodegraph.NodeGraphDataHandler)
}

func (s stream) DistributeReassembledStream(reassembledStream plugin_api.ReassembledStream) {
	metrics.StoreNodegraphMetric(reassembledStream)
}

func (s stream) DistributeTCPPacketPayload(_ plugin_api.TCPPacketPayload) {
	//silent
}

func init() {}

// exported
var StreamPlugin stream
