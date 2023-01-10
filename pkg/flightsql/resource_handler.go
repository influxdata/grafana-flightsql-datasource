package flightsql

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/apache/arrow/go/v10/arrow/flight/flightsql"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/resource/httpadapter"
	"google.golang.org/grpc/metadata"
)

func newResourceHandler(client *flightsql.Client, db string) backend.CallResourceHandler {
	mux := http.NewServeMux()
	// mux.HandleFunc("/get-tables", getTables)

	mux.HandleFunc("/get-tables", func(w http.ResponseWriter, req *http.Request) {
		getTables(w, req, client, db)
	})

	mux.HandleFunc("/get-columns", getColumns)

	return httpadapter.New(mux)
}

// type getTablesResponse struct {
// 	Tables []string `json:"tables"`
// }

func getTables(w http.ResponseWriter, r *http.Request, client *flightsql.Client, db string) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	ctx := context.Background()
	sql := "SELECT TABLE_NAME FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_TYPE='BASE TABLE'"

	tableResp := query(ctx, client, sql, db)

	jsonData, err := tableResp.MarshalJSON()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(jsonData)
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
	// SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = N'coindesk'
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

func query(ctx context.Context, client *flightsql.Client, sql string, db string) backend.DataResponse {
	ctx = metadata.AppendToOutgoingContext(ctx, mdBucketName, db)

	info, err := client.Execute(ctx, sql)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusInternal, fmt.Sprintf("flightsql: %s", err))
	}

	reader, err := client.DoGet(ctx, info.Endpoint[0].Ticket)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusInternal, fmt.Sprintf("flightsql: %s", err))
	}
	defer reader.Release()

	var resp backend.DataResponse
	frame := newFrame(reader.Schema(), sql)
	for reader.Next() {
		record := reader.Record()
		for i, col := range record.Columns() {
			copyData(frame.Fields[i], col)
		}

		if err := reader.Err(); err != nil && !errors.Is(err, io.EOF) {
			resp.Error = err
			break
		}
	}
	resp.Frames = append(resp.Frames, frame)
	return resp
}
