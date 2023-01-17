package flightsql

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/apache/arrow/go/v10/arrow"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"google.golang.org/grpc/metadata"
)

// QueryData executes batches of ad-hoc queries and returns a batch of results.
func (d *FlightSQLDatasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	response := backend.NewQueryDataResponse()
	for _, qreq := range req.Queries {
		var q queryRequest
		if err := json.Unmarshal(qreq.JSON, &q); err != nil {
			response.Responses[qreq.RefID] = backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("unmarshal query request: %s", err))
			continue
		}
		response.Responses[qreq.RefID] = d.query(ctx, q.Text)
	}
	return response, nil
}

// queryRequest is an inbound query request as part of a batch of queries sent
// to (*FlightSQLDatasource).QueryData.
type queryRequest struct {
	RefID                string `json:"refId"`
	Text                 string `json:"queryText"`
	IntervalMilliseconds int    `json:"intervalMs"`
	MaxDataPoints        int    `json:"maxDataPoints"`
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

	// We've implemented our own conversions from Arrow to Data Frame, because
	// the Arrow dependency bundled with Grafana SDK is ancient. If we were to
	// use their functions, we'd end up writing the same amount of conversion
	// code to adapt the APIs.
	var resp backend.DataResponse
	for reader.Next() {
		record := reader.Record()
		schema := record.Schema()

		// Detect whether or not we should convert this Frame into the wide
		// format for time-series. For now this is pretty simplistic. We should
		// consider improving this in the future by reading column metadata.
		//
		// https://grafana.com/docs/grafana/latest/developers/plugins/data-frames/#wide-format
		// https://grafana.com/docs/grafana/latest/developers/plugins/data-frames/#long-format
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
			copyData(frame.Fields[i], col)
		}

		if hasTimeField && hasStringField {
			// Convert the long format we received into wide, because we're
			// pretty sure this is time-series data. This will produce a table
			// that contains columns for each distinct label pair for a matching
			// `time` value.
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
