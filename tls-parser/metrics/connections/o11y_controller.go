package connections

import (
	"encoding/json"
	"fmt"
	"github.com/k8spacket/k8s-api"
	"github.com/k8spacket/plugins/tls-parser/metrics/model"
	"io"
	"net/http"
	"os"
	"reflect"
	"strings"
)

func TLSParserConnectionsHandler(w http.ResponseWriter, req *http.Request) {
	idParam := strings.TrimPrefix(req.URL.Path, "/tlsparser/api/data/")
	if len(strings.TrimSpace(idParam)) > 0 {
		resultFunc := func(destination, source metrics.TLSDetails) metrics.TLSDetails {
			if !reflect.DeepEqual(source, metrics.TLSDetails{}) {
				return source
			} else {
				return destination
			}
		}
		buildResponse(w, fmt.Sprintf("http://%%s:%s/tlsparser/connections/%s?%s", os.Getenv("K8S_PACKET_TCP_LISTENER_PORT"), idParam, req.URL.Query().Encode()), resultFunc)
	} else {
		resultFunc := func(destination, source []metrics.TLSConnection) []metrics.TLSConnection {
			return append(destination, source...)
		}
		buildResponse(w, fmt.Sprintf("http://%%s:%s/tlsparser/connections/?%s", os.Getenv("K8S_PACKET_TCP_LISTENER_PORT"), req.URL.Query().Encode()), resultFunc)
	}
}

func buildResponse[T metrics.TLSDetails | []metrics.TLSConnection](w http.ResponseWriter, url string, resultFunc func(d T, s T) T) {
	var k8spacketIps = k8s.GetPodIPsByLabel("name", os.Getenv("K8S_PACKET_NAME_LABEL_VALUE"))

	var in T
	var out T

	for _, ip := range k8spacketIps {
		resp, err := http.Get(fmt.Sprintf(url, ip))

		if err != nil {
			continue
		}

		responseData, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		err = json.Unmarshal(responseData, &in)
		if err != nil {
			continue
		}

		out = resultFunc(out, in)
	}

	w.Header().Set("Content-Type", "application/json")
	if !reflect.ValueOf(out).IsZero() {
		_ = json.NewEncoder(w).Encode(out)
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found 404"))
	}
}
