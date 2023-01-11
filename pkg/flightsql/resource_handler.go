package flightsql

import (
	"context"
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

	mux.HandleFunc("/get-tables", func(w http.ResponseWriter, req *http.Request) {
		getTables(w, req, client, db)
	})

	mux.HandleFunc("/get-columns", func(w http.ResponseWriter, req *http.Request) {
		getColumns(w, req, client, db)
	})

	return httpadapter.New(mux)
}

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

func getColumns(w http.ResponseWriter, r *http.Request, client *flightsql.Client, db string) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	ctx := context.Background()
	tableName := r.URL.Query().Get("table")
	sql := fmt.Sprintf("select column_name from INFORMATION_SCHEMA.COLUMNS where table_name = '%s'", tableName)
	columnResp := query(ctx, client, sql, db)

	jsonData, err := columnResp.MarshalJSON()

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
