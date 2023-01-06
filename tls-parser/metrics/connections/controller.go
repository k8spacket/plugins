package connections

import (
	"encoding/json"
	"github.com/k8spacket/plugins/tls-parser/metrics/db"
	"github.com/k8spacket/plugins/tls-parser/metrics/model"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

func TLSConnectionHandler(w http.ResponseWriter, req *http.Request) {
	idParam := strings.TrimPrefix(req.URL.Path, "/tlsparser/connections/")
	var id, _ = strconv.Atoi(idParam)
	if id > 0 {
		w.Header().Set("Content-Type", "application/json")
		var tlsDetails = db.Read(id, metrics.TLSDetails{})
		if !reflect.DeepEqual(tlsDetails, metrics.TLSDetails{}) {
			_ = json.NewEncoder(w).Encode(db.Read(id, metrics.TLSDetails{}))
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Not Found 404"))
		}
	} else {
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(db.ReadAll(metrics.TLSConnection{}))
		if err != nil {
			panic(err)
		}
	}
}
