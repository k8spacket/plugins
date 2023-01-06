package connections

import (
	"encoding/json"
	"fmt"
	"github.com/k8spacket/k8s-api"
	"github.com/k8spacket/plugins/tls-parser/metrics/model"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func TLSParserConnectionsHandler(w http.ResponseWriter, req *http.Request) {
	idParam := strings.TrimPrefix(req.URL.Path, "/api/data/")
	var id, _ = strconv.Atoi(idParam)
	if id > 0 {
		buildResponse(w, req, fmt.Sprintf("http://%%s:%s/tlsparser/connections/%s?%s", os.Getenv("K8S_PACKET_TCP_LISTENER_PORT"), id, req.URL.Query().Encode()))
	} else {
		buildResponse(w, req, fmt.Sprintf("http://%%s:%s/tlsparser/connections?%s", os.Getenv("K8S_PACKET_TCP_LISTENER_PORT"), req.URL.Query().Encode()))
	}

}

func buildResponse(w http.ResponseWriter, r *http.Request, url string) {
	var k8spacketIps = k8s.GetPodIPsByLabel("name", os.Getenv("K8S_PACKET_NAME_LABEL_VALUE"))

	var in []metrics.TLSConnection
	var tlsConnectionItems []metrics.TLSConnection

	for _, ip := range k8spacketIps {
		resp, err := http.Get(fmt.Sprintf(url, ip, os.Getenv("K8S_PACKET_TCP_LISTENER_PORT"), r.URL.Query().Encode()))

		if err != nil {
			fmt.Print(err.Error())
			os.Exit(1)
		}

		responseData, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		err = json.Unmarshal(responseData, &in)
		if err != nil {
			panic(err)
		}

		tlsConnectionItems = append(tlsConnectionItems, in...)
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(tlsConnectionItems)
	if err != nil {
		panic(err)
	}
}
