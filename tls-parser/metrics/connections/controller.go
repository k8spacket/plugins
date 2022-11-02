package connections

import (
	"encoding/json"
	"fmt"
	"github.com/k8spacket/plugins/tls-parser/metrics/model"
	"net/http"
	"sync"
)

var (
	tlsConnectionItems      = make(map[string]metrics.TLSConnection)
	tlsConnectionItemsMutex = sync.RWMutex{}
)

func TLSConnectionHandler(w http.ResponseWriter, _ *http.Request) {
	tlsConnectionItemsMutex.RLock()
	values := make([]metrics.TLSConnection, 0, len(tlsConnectionItems))
	for _, v := range tlsConnectionItems {
		values = append(values, v)
	}
	tlsConnectionItemsMutex.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(values)
	if err != nil {
		panic(err)
	}
}

func AddTLSConnection(tlsConnection metrics.TLSConnection) {
	tlsConnectionItemsMutex.Lock()
	var key = fmt.Sprintf("%s-%s", tlsConnection.Src, tlsConnection.Dst)
	tlsConnectionItems[key] = tlsConnection
	tlsConnectionItemsMutex.Unlock()
}
