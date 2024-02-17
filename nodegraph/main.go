package main

import (
	"github.com/k8spacket/plugin-api/v2"
	"github.com/k8spacket/plugins/nodegraph/log"
	"github.com/k8spacket/plugins/nodegraph/metrics"
	"github.com/k8spacket/plugins/nodegraph/metrics/nodegraph"
)

type event plugin_api.TCPEvent

func (e event) InitPlugin(manager plugin_api.PluginManager) {
	nodegraph_log.BuildLogger()

	manager.RegisterTCPPlugin(e)
	manager.RegisterHttpHandler("/nodegraph/connections", nodegraph.ConnectionHandler)
	manager.RegisterHttpHandler("/nodegraph/api/health", nodegraph.Health)
	manager.RegisterHttpHandler("/nodegraph/api/graph/fields", nodegraph.NodeGraphFieldsHandler)
	manager.RegisterHttpHandler("/nodegraph/api/graph/data", nodegraph.NodeGraphDataHandler)
}

func (e event) DistributeTCPEvent(event plugin_api.TCPEvent) {
	metrics.StoreNodegraphMetric(event)
}

func init() {}

// exported
var TCPConsumerPlugin event
