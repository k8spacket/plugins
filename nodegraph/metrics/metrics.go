package metrics

import (
	"github.com/k8spacket/plugin-api"
	"github.com/k8spacket/plugins/nodegraph/log"
	"github.com/k8spacket/plugins/nodegraph/metrics/nodegraph"
	"github.com/k8spacket/plugins/nodegraph/metrics/prometheus"
	"os"
	"strconv"
)

func StoreNodegraphMetric(stream plugin_api.ReassembledStream) {
	hideSrcPort, _ := strconv.ParseBool(os.Getenv("K8S_PACKET_HIDE_SRC_PORT"))
	var srcPortMetrics = stream.SrcPort
	if hideSrcPort {
		srcPortMetrics = "dynamic"
	}

	prometheus.K8sPacketBytesSentMetric.WithLabelValues(stream.SrcNamespace, stream.Src, stream.SrcName, srcPortMetrics, stream.Dst, stream.DstName, stream.DstPort, strconv.FormatBool(stream.Closed)).Observe(stream.BytesSent)
	prometheus.K8sPacketBytesReceivedMetric.WithLabelValues(stream.SrcNamespace, stream.Src, stream.SrcName, srcPortMetrics, stream.Dst, stream.DstName, stream.DstPort, strconv.FormatBool(stream.Closed)).Observe(stream.BytesReceived)
	prometheus.K8sPacketDurationSecondsMetric.WithLabelValues(stream.SrcNamespace, stream.Src, stream.SrcName, srcPortMetrics, stream.Dst, stream.DstName, stream.DstPort, strconv.FormatBool(stream.Closed)).Observe(stream.Duration)

	nodegraph.UpdateNodeGraph(stream.Src, stream.SrcName, stream.SrcNamespace, stream.Dst, stream.DstName, stream.DstNamespace, stream.Closed, stream.BytesSent, stream.BytesReceived, stream.Duration)

	nodegraph_log.LOGGER.Printf("Connection: src=%v srcName=%v srcPort=%v srcNS=%v dst=%v dstName=%v dstPort=%v dstNS=%v closed=%v bytesSent=%v bytesReceived=%v duration=%v",
		stream.Src, stream.SrcName, stream.SrcPort, stream.SrcNamespace, stream.Dst, stream.DstName, stream.DstPort, stream.DstNamespace, stream.Closed, stream.BytesSent, stream.BytesReceived, stream.Duration)
}
