package flightsql

import (
	"context"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/apache/arrow/go/v10/arrow/array"
	"github.com/apache/arrow/go/v10/arrow/flight"
	"github.com/apache/arrow/go/v10/arrow/flight/flightsql"
	"github.com/apache/arrow/go/v10/arrow/memory"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"google.golang.org/grpc/metadata"
)

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

	if !reader.Next() {
		http.Error(w, "table not found", http.StatusNotFound)
		return
	}
	rec := reader.Record()
	rec.Retain()
	defer rec.Release()
	reader.Next()
	if err := reader.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	indices := rec.Schema().FieldIndices("table_schema")
	if len(indices) == 0 {
		http.Error(w, "table_schema field not found", http.StatusInternalServerError)
		return
	}
	col := rec.Column(indices[0])
	serializedSchema := array.NewStringData(col.Data()).Value(0)
	schema, err := flight.DeserializeSchema([]byte(serializedSchema), memory.DefaultAllocator)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var resp backend.DataResponse
	resp.Frames = append(resp.Frames, newFrame(schema, ""))
	if err := writeDataResponse(w, resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func newDataResponse(reader *flight.Reader) backend.DataResponse {
	var resp backend.DataResponse
	frame := newFrame(reader.Schema(), "")
READER:
	for reader.Next() {
		record := reader.Record()
		for i, col := range record.Columns() {
			if err := copyData(frame.Fields[i], col); err != nil {
				resp.Error = err
				break READER
			}
		}
		if err := reader.Err(); err != nil && !errors.Is(err, io.EOF) {
			resp.Error = err
			break
		}
	}
	resp.Frames = append(resp.Frames, frame)
	return resp
}

func writeDataResponse(w io.Writer, resp backend.DataResponse) error {
	json, err := resp.MarshalJSON()
	if err != nil {
		return err
	}
	_, err = w.Write(json)
	return err
}
