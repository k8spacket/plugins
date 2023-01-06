package db

import (
	"encoding/json"
	"github.com/HouzuoGuo/tiedot/db"
	"github.com/fatih/structs"
	"github.com/k8spacket/plugins/tls-parser/metrics/model"
)

const connectionColName = "TLSConnections"
const detailsColName = "TLSDetails"

var connectionsCol, detailsCol = buildDatabase()

func buildDatabase() (*db.Col, *db.Col) {

	dbDir := "./Database"

	database, err := db.OpenDB(dbDir)
	if err != nil {
		panic(err)
	}

	if database.ColExists(connectionColName) == false {
		if err := database.Create(connectionColName); err != nil {
			panic(err)
		}
		if err := database.Use(connectionColName).Index([]string{"id"}); err != nil {
			panic(err)
		}
	}
	if database.ColExists(detailsColName) == false {
		if err := database.Create(detailsColName); err != nil {
			panic(err)
		}
		if err := database.Use(detailsColName).Index([]string{"id"}); err != nil {
			panic(err)
		}
	}

	return database.Use(connectionColName), database.Use(detailsColName)
}

func Insert[T metrics.TLSDetails | metrics.TLSConnection](id int, document T) {
	var col = getCol(document)
	var doc, _ = col.Read(id)
	if len(doc) > 0 {
		_ = col.Update(id, structs.Map(document))
	} else {
		_ = col.InsertRecovery(id, structs.Map(document))
	}
}

func Read[T metrics.TLSDetails | metrics.TLSConnection](docId int, s T) T {
	var document, _ = getCol(s).Read(docId)
	var jsonBytes, _ = json.Marshal(document)
	_ = json.Unmarshal(jsonBytes, &s)
	return s
}

func ReadAll[T metrics.TLSDetails | metrics.TLSConnection](s T) []T {
	col := getCol(s)

	var query interface{}
	err := json.Unmarshal([]byte(`["all"]`), &query)
	if err != nil {
		return nil
	}

	queryResult := make(map[int]struct{})
	if err := db.EvalQuery(query, col, &queryResult); err != nil {
		panic(err)
	}

	var result []T
	for id := range queryResult {
		result = append(result, Read(id, s))
	}
	return result
}

func getCol[T metrics.TLSDetails | metrics.TLSConnection](s T) *db.Col {
	switch any(s).(type) {
	case metrics.TLSConnection:
		return connectionsCol
	case metrics.TLSDetails:
		return detailsCol
	}
	return nil
}
