package connections

import (
	"encoding/json"
	tls_parser_log "github.com/k8spacket/plugins/tls-parser/log"
	tls_connection_db "github.com/k8spacket/plugins/tls-parser/metrics/db/tls_connection"
	tls_detail_db "github.com/k8spacket/plugins/tls-parser/metrics/db/tls_detail"
	"github.com/k8spacket/plugins/tls-parser/metrics/model"
	"net/http"
	"reflect"
	"strings"
)

func TLSConnectionHandler(w http.ResponseWriter, req *http.Request) {
	id := strings.TrimPrefix(req.URL.Path, "/tlsparser/connections/")
	if len(id) > 0 {
		w.Header().Set("Content-Type", "application/json")
		var details = tls_detail_db.Read(id)
		if !reflect.DeepEqual(details, model.TLSDetails{}) {
			err := json.NewEncoder(w).Encode(details)
			if err != nil {
				tls_parser_log.LOGGER.Printf("[api] Cannot prepare connection details response: %+v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Not Found 404"))
		}
	} else {
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(tls_connection_db.ReadAll())
		if err != nil {
			tls_parser_log.LOGGER.Printf("[api] Cannot prepare connections response: %+v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
