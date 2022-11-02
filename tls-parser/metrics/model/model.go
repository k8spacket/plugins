package metrics

type TLSConnection struct {
	Src          string `json:"src"`
	SrcName      string `json:"srcName"`
	SrcNamespace string `json:"srcNamespace"`
	Dst          string `json:"dst"`
	DstName      string `json:"dstName"`
	DstPort      string `json:"dstPort"`
	TLS          TLS    `json:"tls"`
}

type TLS struct {
	Domain             string   `json:"domain"`
	ClientTLSVersions  []string `json:"clientTLSVersions"`
	ClientCipherSuites []string `json:"clientCipherSuites"`
	UsedTLSVersion     string   `json:"usedTLSVersion"`
	UsedCipherSuite    string   `json:"usedCipherSuite"`
}
