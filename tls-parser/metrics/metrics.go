package metrics

import (
	"encoding/json"
	"github.com/k8spacket/plugin-api"
	tls_parser_log "github.com/k8spacket/plugins/tls-parser/log"
	"github.com/k8spacket/plugins/tls-parser/metrics/connections"
	"github.com/k8spacket/plugins/tls-parser/metrics/model"
	"github.com/k8spacket/plugins/tls-parser/metrics/prometheus"
	"github.com/k8spacket/tls-api"
	"github.com/k8spacket/tls-api/model"
)

var tlsConnectionMap = make(map[uint32]metrics.TLSConnection)

func StoreStreamMetrics(reassembledStream plugin_api.ReassembledStream) {
	tlsConnection, ok := tlsConnectionMap[reassembledStream.StreamId]
	if ok {
		tlsConnection.SrcNamespace = reassembledStream.SrcNamespace
		tlsConnection.Src = reassembledStream.Src
		tlsConnection.SrcName = reassembledStream.SrcName
		tlsConnection.Dst = reassembledStream.Dst
		tlsConnection.DstName = reassembledStream.DstName
		tlsConnection.DstPort = reassembledStream.DstPort
		prometheus.K8sPacketTLSRecordMetric.WithLabelValues(
			tlsConnection.SrcNamespace,
			tlsConnection.Src,
			tlsConnection.SrcName,
			tlsConnection.Dst,
			tlsConnection.DstName,
			tlsConnection.DstPort,
			tlsConnection.TLS.Domain,
			tlsConnection.TLS.UsedTLSVersion,
			tlsConnection.TLS.UsedCipherSuite).Add(1)
		connections.AddTLSConnection(tlsConnection)
		var j, _ = json.Marshal(tlsConnection)
		tls_parser_log.LOGGER.Println("TLS Record:", string(j))
		delete(tlsConnectionMap, reassembledStream.StreamId)
	}
}

func CollectTCPPacketPayload(tcpPacketPayload plugin_api.TCPPacketPayload) {
	payload := tcpPacketPayload.Payload
	tlsConnection, ok := tlsConnectionMap[tcpPacketPayload.StreamId]
	if !ok {
		tlsConnection = metrics.TLSConnection{}
	}
	if len(payload) > 5 && payload[0] == model.TLSRecord {
		if payload[5] == model.ClientHelloTLS {
			var record = tls_api.ParseTLSPayload(payload).(model.ClientHelloTLSRecord)
			tlsConnection.TLS.Domain = record.ResolvedClientFields.ServerName
			tlsConnection.TLS.ClientTLSVersions = record.ResolvedClientFields.SupportedVersions
			tlsConnection.TLS.ClientCipherSuites = record.ResolvedClientFields.Ciphers
		} else if payload[5] == model.ServerHelloTLS {
			var record = tls_api.ParseTLSPayload(payload).(model.ServerHelloTLSRecord)
			tlsConnection.TLS.UsedTLSVersion = record.ResolvedServerFields.SupportedVersion
			tlsConnection.TLS.UsedCipherSuite = record.ResolvedServerFields.Cipher
		}
		tlsConnectionMap[tcpPacketPayload.StreamId] = tlsConnection
	}
}
