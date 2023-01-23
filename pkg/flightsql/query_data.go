package flightsql

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
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
// to [(*FlightSQLDatasource).QueryData].
type queryRequest struct {
	RefID                string `json:"refId"`
	Text                 string `json:"queryText"`
	IntervalMilliseconds int    `json:"intervalMs"`
	MaxDataPoints        int64  `json:"maxDataPoints"`
}

// query executes a SQL statement by issuing a `CommandStatementQuery` command to Flight SQL.
func (d *FlightSQLDatasource) query(ctx context.Context, sql string) backend.DataResponse {
	ctx = metadata.NewOutgoingContext(ctx, d.md)

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
