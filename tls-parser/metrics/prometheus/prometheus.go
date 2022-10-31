package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	K8sPacketTLSRecordMetric = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "k8s_packet_tls_record",
			Help: "Kubernetes packet TLS Record",
		},
		[]string{"ns", "src", "src_name", "dst", "dst_name", "dst_port", "domain", "tls_version", "cipher_suite"},
	)
)

func init() {
	prometheus.MustRegister(K8sPacketTLSRecordMetric)
}
