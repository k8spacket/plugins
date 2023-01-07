package metrics

type TLSConnection struct {
	Id              uint32 `json:"id"`
	StreamId        uint32 `json:"streamId"`
	Src             string `json:"src"`
	SrcName         string `json:"srcName"`
	SrcNamespace    string `json:"srcNamespace"`
	Dst             string `json:"dst"`
	DstName         string `json:"dstName"`
	DstPort         string `json:"dstPort"`
	Domain          string `json:"domain"`
	UsedTLSVersion  string `json:"usedTLSVersion"`
	UsedCipherSuite string `json:"usedCipherSuite"`
}

type TLSDetails struct {
	Id                 uint32   `json:"id"`
	StreamId           uint32   `json:"streamId"`
	Domain             string   `json:"domain"`
	ClientTLSVersions  []string `json:"clientTLSVersions"`
	ClientCipherSuites []string `json:"clientCipherSuites"`
	UsedTLSVersion     string   `json:"usedTLSVersion"`
	UsedCipherSuite    string   `json:"usedCipherSuite"`
	ServerChain        string   `json:"serverChain"`
}
