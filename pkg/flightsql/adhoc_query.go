package flightsql

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/apache/arrow/go/v10/arrow"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/grafana-plugin-sdk-go/data/sqlutil"
	"google.golang.org/grpc/metadata"
)

// QueryData executes batches of ad-hoc queries and returns a batch of results.
func (d *FlightSQLDatasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	response := backend.NewQueryDataResponse()
	for _, qreq := range req.Queries {
		// The shape of the payload we're parsing from backend.DataQuery.JSON
		// doesn't have the same fields as what sqlutil.Query expects to parse
		// when using sqlutil.GetQuery. Parse it ourselves and create a
		// sqlutil.Query value.
		var q queryRequest
		if err := json.Unmarshal(qreq.JSON, &q); err != nil {
			response.Responses[qreq.RefID] = backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("unmarshal query request: %s", err))
			continue
		}
		query := &sqlutil.Query{
			RawSQL:        q.Text,
			RefID:         q.RefID,
			MaxDataPoints: q.MaxDataPoints,
			Interval:      time.Duration(q.IntervalMilliseconds) * time.Millisecond,
			TimeRange:     qreq.TimeRange,
			// There are other fields here that are worth looking into. Things
			// seem to work just fine, but Format and FillMissing seem useful.
		}

		// Process macros and execute the query.
		sql, err := sqlutil.Interpolate(query, macros)
		if err != nil {
			response.Responses[qreq.RefID] = backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("macro interpolation: %s", err))
			continue
		}
		response.Responses[qreq.RefID] = d.query(ctx, sql)
	}
	return response, nil
}

// queryRequest is an inbound query request as part of a batch of queries sent
// to (*FlightSQLDatasource).QueryData.
type queryRequest struct {
	RefID                string `json:"refId"`
	Text                 string `json:"queryText"`
	IntervalMilliseconds int    `json:"intervalMs"`
	MaxDataPoints        int64  `json:"maxDataPoints"`
}

// query executes a SQL statement by issuing a `CommandStatementQuery` command to Flight SQL.
func (d *FlightSQLDatasource) query(ctx context.Context, sql string) backend.DataResponse {
	ctx = metadata.AppendToOutgoingContext(ctx, mdBucketName, d.database)

	info, err := d.client.Execute(ctx, sql)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusInternal, fmt.Sprintf("flightsql: %s", err))
	}
	if len(info.Endpoint) != 1 {
		return backend.ErrDataResponse(backend.StatusInternal, fmt.Sprintf("unsupported endpoint count in response: %d", len(info.Endpoint)))
	}
	reader, err := d.client.DoGet(ctx, info.Endpoint[0].Ticket)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusInternal, fmt.Sprintf("flightsql: %s", err))
	}
	defer reader.Release()

	return newQueryDataResponse(reader, sql)
}

type flightReader interface {
	Next() bool
	Record() arrow.Record
	Err() error
}

// newQueryDataResponse builds a DataResponse from a stream of Arrow data. The
// DataResponse is divided into multiple Frames.
//
// If the data contains a `time` column of type `timestamp` and at least one
// column of type `string`, the resulting table will be converted to wide
// format.
//
// Data Frame Formatting Reference:
// - https://grafana.com/docs/grafana/latest/developers/plugins/data-frames/#wide-format
// - https://grafana.com/docs/grafana/latest/developers/plugins/data-frames/#long-format
func newQueryDataResponse(reader flightReader, sql string) backend.DataResponse {
	var resp backend.DataResponse

READER:
	for reader.Next() {
		record := reader.Record()
		schema := record.Schema()

		var (
			hasTimeField   bool
			hasStringField bool
		)
		for _, f := range schema.Fields() {
			if f.Name == "time" && f.Type.ID() == arrow.TIMESTAMP {
				hasTimeField = true
			} else if f.Type.ID() == arrow.STRING {
				hasStringField = true
			}
		}

		frame := newFrame(schema, sql)
		for i, col := range record.Columns() {
			if err := copyData(frame.Fields[i], col); err != nil {
				resp.Error = err
				break READER
			}
		}
		if hasTimeField && hasStringField {
			var err error
			frame, err = data.LongToWide(frame, nil)
			if err != nil {
				resp.Error = err
				break
			}
		}
		resp.Frames = append(resp.Frames, frame)

		if err := reader.Err(); err != nil && !errors.Is(err, io.EOF) {
			resp.Error = err
			break
		}
	}
	return resp
}
