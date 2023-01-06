package flightsql

import (
	"encoding/json"
	"net/http"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/resource/httpadapter"
)

func newResourceHandler() backend.CallResourceHandler {
	mux := http.NewServeMux()
	mux.HandleFunc("/get-tables", getTables)
	mux.HandleFunc("/get-columns", getColumns)

	return httpadapter.New(mux)
}

type getTablesResponse struct {
	Tables []string `json:"tables"`
}

func getTables(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	// change this to a sql query
	// sql := "SELECT * FROM information_schema.tables"
	// map to response object that looks like this:

	tables := &getTablesResponse{
		Tables: []string{
			"coindesk",
			"airSensor",
		},
	}

	j, err := json.Marshal(tables)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(j)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

type getColumnsResponse struct {
	Columns []string `json:"columns"`
}

func getColumns(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	// change this to a sql query
	// sql := select * from information_schema.columns where table_name='tableName'
	// map to response object that looks like this:
	columns := &getColumnsResponse{
		Columns: []string{
			"time",
			"price",
		},
	}

	j, err := json.Marshal(columns)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(j)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
