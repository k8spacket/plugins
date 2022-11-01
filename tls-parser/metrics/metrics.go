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
	Client       struct {
		Domain       string
		TlsVersions  []string
		CipherSuites []string
	}
	Server struct {
		TlsVersion  string
		CipherSuite string
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
			tlsRecord.Client.Domain,
			tlsRecord.Server.TlsVersion,
			tlsRecord.Server.CipherSuite).Add(1)
		var j, _ = json.MarshalIndent(tlsRecord, "", "  ")
		tls_parser_log.LOGGER.Println("TLS Record:", j)
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
			tlsRecord.Client.Domain = record.ResolvedClientFields.ServerName
			tlsRecord.Client.TlsVersions = record.ResolvedClientFields.SupportedVersions
			tlsRecord.Client.CipherSuites = record.ResolvedClientFields.Ciphers
		} else if payload[5] == model.ServerHelloTLS {
			var record = tls_api.ParseTLSPayload(payload).(model.ServerHelloTLSRecord)
			tlsRecord.Server.TlsVersion = record.ResolvedServerFields.SupportedVersion
			tlsRecord.Server.CipherSuite = record.ResolvedServerFields.Cipher
		}
		tlsRecordMap[tcpPacketPayload.StreamId] = tlsRecord
	}
}
