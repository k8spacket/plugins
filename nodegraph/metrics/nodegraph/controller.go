package nodegraph

import (
	"encoding/json"
	nodegraph_log "github.com/k8spacket/plugins/nodegraph/log"
	tcp_connection_db "github.com/k8spacket/plugins/nodegraph/metrics/db/tcp_connection"
	"github.com/k8spacket/plugins/nodegraph/metrics/nodegraph/model"
	"net/http"
	"net/url"
	"regexp"
)

func ConnectionHandler(w http.ResponseWriter, r *http.Request) {
	connectionItemsMutex.RLock()
	var response = filterConnections(r.URL.Query())
	connectionItemsMutex.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		nodegraph_log.LOGGER.Printf("[api] Cannot prepare connections response: %+v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func filterConnections(query url.Values) map[string]model.ConnectionItem {
	var namespace = query["namespace"]
	var patternNs = ""
	if len(namespace) > 0 {
		patternNs = namespace[0]
	}

	var exclude = query["exclude"]
	var patternEx = ""
	if len(exclude) > 0 {
		patternEx = exclude[0]
	}

	var include = query["include"]
	var patternIn = ""
	if len(include) > 0 {
		patternIn = include[0]
	}

	var filteredConnectionItems = make(map[string]model.ConnectionItem)

	for _, conn := range tcp_connection_db.ReadAll() {
		var matchSrc, _ = regexp.Match(patternNs, []byte(conn.SrcNamespace))
		var matchDst, _ = regexp.Match(patternNs, []byte(conn.DstNamespace))

		var excludeSrc, _ = regexp.Match(patternEx, []byte(conn.SrcName+conn.Src))
		var excludeDst, _ = regexp.Match(patternEx, []byte(conn.DstName+conn.Dst))

		var includeSrc, _ = regexp.Match(patternIn, []byte(conn.SrcName+conn.Src))
		var includeDst, _ = regexp.Match(patternIn, []byte(conn.DstName+conn.Dst))

		if (patternNs == "" || matchSrc || matchDst) && (patternEx == "" || (!excludeSrc && !excludeDst)) && (patternIn == "" || (includeSrc || includeDst)) {
			filteredConnectionItems[conn.Src+"-"+conn.Dst] = conn
		}
	}
	return filteredConnectionItems
}
