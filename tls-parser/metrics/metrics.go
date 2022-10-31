package metrics

import (
	"bytes"
	"crypto/tls"
	"encoding/binary"
	"github.com/k8spacket/plugin-api"
	tls_parser_log "github.com/k8spacket/plugins/tls-parser/log"
	"github.com/k8spacket/plugins/tls-parser/metrics/prometheus"
	"github.com/k8spacket/tls-api"
	"github.com/k8spacket/tls-api/model"
)

type TLSRecord struct {
	srcNamespace string
	src          string
	srcName      string
	dst          string
	dstName      string
	dstPort      string
	domain       string
	tlsVersion   string
	cipherSuite  string
}

var tlsRecordMap = make(map[uint32]TLSRecord)

func StoreStreamMetrics(reassembledStream plugin_api.ReassembledStream) {
	tlsRecord, ok := tlsRecordMap[reassembledStream.StreamId]
	if ok {
		tlsRecord.srcNamespace = reassembledStream.SrcNamespace
		tlsRecord.src = reassembledStream.Src
		tlsRecord.srcName = reassembledStream.SrcName
		tlsRecord.dst = reassembledStream.Dst
		tlsRecord.dstName = reassembledStream.DstName
		tlsRecord.dstPort = reassembledStream.DstPort
		prometheus.K8sPacketTLSRecordMetric.WithLabelValues(
			tlsRecord.srcNamespace,
			tlsRecord.src,
			tlsRecord.srcName,
			tlsRecord.dst,
			tlsRecord.dstName,
			tlsRecord.dstPort,
			tlsRecord.domain,
			tlsRecord.tlsVersion,
			tlsRecord.cipherSuite).Add(1)
		tls_parser_log.LOGGER.Printf("TLS Record: src=%v srcName=%v srcNS=%v dst=%v dstName=%v dstPort=%v domain=%v tlsVersion=%v cipherSuite=%v",
			tlsRecord.src, tlsRecord.srcName, tlsRecord.srcNamespace, tlsRecord.dst, tlsRecord.dstName, tlsRecord.dstPort, tlsRecord.domain, tlsRecord.tlsVersion, tlsRecord.cipherSuite)
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
		if payload[5] == model.ClientHelloTLSRecord {
			var record = tls_api.ParseTLSPayload(payload).(tls_api.ClientHelloTLSRecord)
			tlsRecord.domain = getServerName(record)
		} else if payload[5] == model.ServerHelloTLSRecord {
			var record = tls_api.ParseTLSPayload(payload).(tls_api.ServerHelloTLSRecord)
			tlsRecord.tlsVersion = model.GetTLSVersion(record.HandshakeProtocol.TLSVersion)
			extension := record.Extensions.Extensions[model.TLSVersionExt]
			if extension.Value != nil {
				tlsRecord.tlsVersion = model.GetTLSVersion(binary.BigEndian.Uint16(extension.Value))
			}
			tlsRecord.cipherSuite = tls.CipherSuiteName(record.CipherSuite.Value)
		}
		tlsRecordMap[tcpPacketPayload.StreamId] = tlsRecord
	}
}

func getServerName(record tls_api.ClientHelloTLSRecord) string {
	extension := record.Extensions.Extensions[model.ServerNameExt]

	var serverNameExtension model.ServerNameExtension

	reader := bytes.NewReader(extension.Value)
	binary.Read(reader, binary.BigEndian, &serverNameExtension.ListLength)
	binary.Read(reader, binary.BigEndian, &serverNameExtension.Type)
	binary.Read(reader, binary.BigEndian, &serverNameExtension.Length)
	serverNameValue := make([]byte, serverNameExtension.Length)
	binary.Read(reader, binary.BigEndian, &serverNameValue)
	serverNameExtension.Value = serverNameValue

	return string(serverNameExtension.Value)
}
