package metrics

import (
	"encoding/json"
	"fmt"
	"github.com/grantae/certinfo"
	"github.com/k8spacket/plugin-api"
	tls_parser_log "github.com/k8spacket/plugins/tls-parser/log"
	"github.com/k8spacket/plugins/tls-parser/metrics/db"
	"github.com/k8spacket/plugins/tls-parser/metrics/model"
	"github.com/k8spacket/plugins/tls-parser/metrics/prometheus"
	"github.com/k8spacket/tls-api"
	"github.com/k8spacket/tls-api/model"
	"hash/fnv"
	"reflect"
	"strings"
)

var tlsConnectionMap = make(map[uint32]metrics.TLSConnection)
var tlsDetailsMap = make(map[uint32]metrics.TLSDetails)

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
			tlsConnection.Domain,
			tlsConnection.UsedTLSVersion,
			tlsConnection.UsedCipherSuite).Add(1)
		tlsDetails, _ := tlsDetailsMap[reassembledStream.StreamId]
		storeInDatabase(tlsConnection, tlsDetails)
		var j, _ = json.Marshal(tlsConnection)
		tls_parser_log.LOGGER.Println("TLS Record:", string(j))
		delete(tlsConnectionMap, reassembledStream.StreamId)
		delete(tlsDetailsMap, reassembledStream.StreamId)
	}
}

func CollectTCPPacketPayload(streamId uint32, payload []byte) {
	tlsConnection, ok := tlsConnectionMap[streamId]
	if !ok {
		tlsConnection = metrics.TLSConnection{}
		tlsConnection.StreamId = streamId
	}
	tlsDetails, ok := tlsDetailsMap[streamId]
	if !ok {
		tlsDetails = metrics.TLSDetails{}
		tlsDetails.StreamId = streamId
	}
	var tlsWrapper = tls_api.ParseTLSPayload(payload)
	if !reflect.DeepEqual(tlsWrapper.ClientHelloTLSRecord, model.ClientHelloTLSRecord{}) {
		var record = tlsWrapper.ClientHelloTLSRecord
		tlsConnection.Domain = record.ResolvedClientFields.ServerName
		tlsDetails.Domain = record.ResolvedClientFields.ServerName
		tlsDetails.ClientTLSVersions = record.ResolvedClientFields.SupportedVersions
		tlsDetails.ClientCipherSuites = record.ResolvedClientFields.Ciphers
	}
	if !reflect.DeepEqual(tlsWrapper.ServerHelloTLSRecord, model.ServerHelloTLSRecord{}) {
		var record = tlsWrapper.ServerHelloTLSRecord
		tlsConnection.UsedTLSVersion = record.ResolvedServerFields.SupportedVersion
		tlsConnection.UsedCipherSuite = record.ResolvedServerFields.Cipher
		tlsDetails.UsedTLSVersion = record.ResolvedServerFields.SupportedVersion
		tlsDetails.UsedCipherSuite = record.ResolvedServerFields.Cipher
		if tlsDetails.UsedTLSVersion == model.GetTLSVersion(0x0304) {
			tlsDetails.ServerChain = "ENCRYPTED"
		}
	}
	if !reflect.DeepEqual(tlsWrapper.CertificateTLSRecord, model.CertificateTLSRecord{}) {
		var record = tlsWrapper.CertificateTLSRecord
		tlsDetails.ServerChain = ""
		for _, cert := range record.Certificates {
			var certString, _ = certinfo.CertificateText(&cert)
			certString = strings.Replace(certString, "\n\n", "\n", -1)
			tlsDetails.ServerChain += certString
		}
	}
	tlsConnectionMap[streamId] = tlsConnection
	tlsDetailsMap[streamId] = tlsDetails
}

func storeInDatabase(tlsConnection metrics.TLSConnection, tlsDetails metrics.TLSDetails) {
	var id = hash(fmt.Sprintf("%s-%s", tlsConnection.Src, tlsConnection.Dst))
	tlsConnection.Id = id
	db.Insert(int(id), tlsConnection)
	tlsDetails.Id = id
	db.Insert(int(id), tlsDetails)
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
