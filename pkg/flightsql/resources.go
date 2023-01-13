package flightsql

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/apache/arrow/go/v10/arrow/flight/flightsql"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"google.golang.org/grpc/metadata"
)

// getSQLInfo calls the `CommandGetSqlInfo` command and returns it as JSON.
//
// It is invoked as a Resource Handler via the CallResource gRPC method.
func (d *FlightSQLDatasource) getSQLInfo(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()
	ctx = metadata.AppendToOutgoingContext(ctx, mdBucketName, d.database)
	info, err := d.client.GetSqlInfo(ctx, []flightsql.SqlInfo{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	reader, err := d.client.DoGet(ctx, info.Endpoint[0].Ticket)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer reader.Release()

	if err := writeDataResponse(w, newDataResponse(reader)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// getTables calls the `CommandGetTables` command and returns it as JSON.
//
// It is invoked as a Resource Handler via the CallResource gRPC method.
func (d *FlightSQLDatasource) getTables(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()
	ctx = metadata.AppendToOutgoingContext(ctx, mdBucketName, d.database)
	info, err := d.client.GetTables(ctx, &flightsql.GetTablesOpts{
		TableTypes: []string{"BASE TABLE"},
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	reader, err := d.client.DoGet(ctx, info.Endpoint[0].Ticket)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer reader.Release()

	if err := writeDataResponse(w, newDataResponse(reader)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// getColumns calls `CommandGetTables` with `include_schemas` for a single
// table. The Arrow Schema of the table is used to produce a JSON response that
// describes the table's structure and has no data rows.
//
// It is invoked as a Resource Handler via the CallResource gRPC method.
func (d *FlightSQLDatasource) getColumns(w http.ResponseWriter, r *http.Request) {
	tableName := r.URL.Query().Get("table")
	if tableName == "" {
		http.Error(w, `query parameter "table" is required`, http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()
	ctx = metadata.AppendToOutgoingContext(ctx, mdBucketName, d.database)
	info, err := d.client.GetTables(ctx, &flightsql.GetTablesOpts{
		TableNameFilterPattern: &tableName,
		IncludeSchema:          true,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	reader, err := d.client.DoGet(ctx, info.Endpoint[0].Ticket)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer reader.Release()

	schemas, err := readTableSchemas(reader)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var resp backend.DataResponse
	resp.Frames = append(resp.Frames, newFrame(schemas[tableName], ""))
	if err := writeDataResponse(w, resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// writeDataResponse writes a Grafana DataResponse to an io.Writer.
func writeDataResponse(w io.Writer, resp backend.DataResponse) error {
	json, err := resp.MarshalJSON()
	if err != nil {
		return err
	}
	_, err = w.Write(json)
	return err
}
