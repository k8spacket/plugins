package metrics

import (
	"encoding/json"
	"github.com/k8spacket/plugin-api"
	tls_parser_log "github.com/k8spacket/plugins/tls-parser/log"
	"github.com/k8spacket/plugins/tls-parser/metrics/prometheus"
	"github.com/k8spacket/tls-api"
	"github.com/k8spacket/tls-api/model"
)

type TLSRecord struct {
	SrcNamespace string
	Src          string
	SrcName      string
	Dst          string
	DstName      string
	DstPort      string
	TLS          struct {
		Domain                string
		SupportedTLSVersions  []string
		SupportedCipherSuites []string
		UsedTLSVersion        string
		UsedCipherSuite       string
	}
}

var tlsRecordMap = make(map[uint32]TLSRecord)

func StoreStreamMetrics(reassembledStream plugin_api.ReassembledStream) {
	tlsRecord, ok := tlsRecordMap[reassembledStream.StreamId]
	if ok {
		tlsRecord.SrcNamespace = reassembledStream.SrcNamespace
		tlsRecord.Src = reassembledStream.Src
		tlsRecord.SrcName = reassembledStream.SrcName
		tlsRecord.Dst = reassembledStream.Dst
		tlsRecord.DstName = reassembledStream.DstName
		tlsRecord.DstPort = reassembledStream.DstPort
		prometheus.K8sPacketTLSRecordMetric.WithLabelValues(
			tlsRecord.SrcNamespace,
			tlsRecord.Src,
			tlsRecord.SrcName,
			tlsRecord.Dst,
			tlsRecord.DstName,
			tlsRecord.DstPort,
			tlsRecord.TLS.Domain,
			tlsRecord.TLS.UsedTLSVersion,
			tlsRecord.TLS.UsedCipherSuite).Add(1)
		var j, _ = json.Marshal(tlsRecord)
		tls_parser_log.LOGGER.Println("TLS Record:", string(j))
		delete(tlsRecordMap, reassembledStream.StreamId)
	}
}

func CollectTCPPacketPayload(tcpPacketPayload plugin_api.TCPPacketPayload) {
	payload := tcpPacketPayload.Payload
	tlsRecord, ok := tlsRecordMap[tcpPacketPayload.StreamId]
	if !ok {
		tlsRecord = TLSRecord{}
	}
	if len(payload) > 5 && payload[0] == model.TLSRecord {
		if payload[5] == model.ClientHelloTLS {
			var record = tls_api.ParseTLSPayload(payload).(model.ClientHelloTLSRecord)
			tlsRecord.TLS.Domain = record.ResolvedClientFields.ServerName
			tlsRecord.TLS.SupportedTLSVersions = record.ResolvedClientFields.SupportedVersions
			tlsRecord.TLS.SupportedCipherSuites = record.ResolvedClientFields.Ciphers
		} else if payload[5] == model.ServerHelloTLS {
			var record = tls_api.ParseTLSPayload(payload).(model.ServerHelloTLSRecord)
			tlsRecord.TLS.UsedTLSVersion = record.ResolvedServerFields.SupportedVersion
			tlsRecord.TLS.UsedCipherSuite = record.ResolvedServerFields.Cipher
		}
		tlsRecordMap[tcpPacketPayload.StreamId] = tlsRecord
	}
}
