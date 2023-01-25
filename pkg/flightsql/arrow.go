package flightsql

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"runtime/debug"
	"time"

	"github.com/apache/arrow/go/v10/arrow"
	"github.com/apache/arrow/go/v10/arrow/array"
	"github.com/apache/arrow/go/v10/arrow/scalar"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/grafana-plugin-sdk-go/data/sqlutil"
)

// TODO(brett): Make this configurable. This is an arbitrary value right
// now. Grafana used to have a 1M row rowLimit established in open-source. I'll
// let users hit that for now until we decide how to proceed.
const rowLimit = 1_000_000

type recordReader interface {
	Next() bool
	Schema() *arrow.Schema
	Record() arrow.Record
	Err() error
}

// newQueryDataResponse builds a [backend.DataResponse] from a stream of
// [arrow.Record]s.
//
// The backend.DataResponse contains a single [data.Frame].
func newQueryDataResponse(reader recordReader, query *sqlutil.Query) backend.DataResponse {
	var resp backend.DataResponse
	frame, err := frameForRecords(reader)
	if err != nil {
		resp.Error = err
	}
	frame.Meta.ExecutedQueryString = query.RawSQL
	frame.Meta.DataTopic = data.DataTopic(query.RawSQL)

	switch query.Format {
	case sqlutil.FormatOptionTimeSeries:
		if _, idx := frame.FieldByName("time"); idx == -1 {
			resp.Error = fmt.Errorf("no time column found")
			return resp
		}
		var err error
		frame, err = data.LongToWide(frame, nil)
		if err != nil {
			resp.Error = err
			return resp
		}
	case sqlutil.FormatOptionTable:
		// No changes to the output. Send it as is.
	case sqlutil.FormatOptionLogs:
		// TODO(brett): We need to find out what this actually is and if its
		// worth supporting. Pass through as "table" for now.
	default:
		resp.Error = fmt.Errorf("unsupported format")
	}

	resp.Frames = data.Frames{frame}
	return resp
}

// frameForRecords creates a [data.Frame] from a stream of [arrow.Record]s.
func frameForRecords(reader recordReader) (*data.Frame, error) {
	var (
		frame = newFrame(reader.Schema())
		rows  int64
	)
	for reader.Next() {
		record := reader.Record()
		for i, col := range record.Columns() {
			if err := copyData(frame.Fields[i], col); err != nil {
				return frame, err
			}
		}

		rows += record.NumRows()
		if rows > rowLimit {
			frame.AppendNotices(data.Notice{
				Severity: data.NoticeSeverityWarning,
				Text:     fmt.Sprintf("Results have been limited to %v because the SQL row limit was reached", rowLimit),
			})
			return frame, nil
		}

		if err := reader.Err(); err != nil && !errors.Is(err, io.EOF) {
			return frame, err
		}
	}
	return frame, nil
}

// newFrame builds a new Data Frame from an Arrow Schema.
func newFrame(schema *arrow.Schema) *data.Frame {
	fields := schema.Fields()
	df := &data.Frame{
		Fields: make([]*data.Field, len(fields)),
		Meta:   &data.FrameMeta{},
	}
	nullable := make([]bool, len(fields))
	for i, field := range fields {
		nullable[i] = field.Nullable
		switch field.Type.ID() {
		case arrow.STRING:
			if field.Nullable {
				var s []*string
				df.Fields[i] = data.NewField(field.Name, nil, s)
				continue
			}
			var s []string
			df.Fields[i] = data.NewField(field.Name, nil, s)
		case arrow.FLOAT64:
			if field.Nullable {
				var s []*float64
				df.Fields[i] = data.NewField(field.Name, nil, s)
				continue
			}
			var s []float64
			df.Fields[i] = data.NewField(field.Name, nil, s)
		case arrow.UINT32:
			if field.Nullable {
				var s []*uint32
				df.Fields[i] = data.NewField(field.Name, nil, s)
				continue
			}
			var s []uint32
			df.Fields[i] = data.NewField(field.Name, nil, s)
		case arrow.INT64:
			if field.Nullable {
				var s []*int64
				df.Fields[i] = data.NewField(field.Name, nil, s)
				continue
			}
			var s []int64
			df.Fields[i] = data.NewField(field.Name, nil, s)
		case arrow.BOOL:
			if field.Nullable {
				var s []*bool
				df.Fields[i] = data.NewField(field.Name, nil, s)
				continue
			}
			var s []bool
			df.Fields[i] = data.NewField(field.Name, nil, s)
		case arrow.TIMESTAMP:
			if field.Nullable {
				var s []*time.Time
				df.Fields[i] = data.NewField(field.Name, nil, s)
				continue
			}
			var s []time.Time
			df.Fields[i] = data.NewField(field.Name, nil, s)
		case arrow.DENSE_UNION:
			if field.Nullable {
				var s []*json.RawMessage
				df.Fields[i] = data.NewField(field.Name, nil, s)
				continue
			}
			var s []json.RawMessage
			df.Fields[i] = data.NewField(field.Name, nil, s)
		}
	}
	return df
}

// copyData copies the contents of an Arrow column into a Data Frame field.
func copyData(field *data.Field, col arrow.Array) error {
	defer func() {
		if r := recover(); r != nil {
			logErrorf("Panic: %v %v", r, string(debug.Stack()))
		}
	}()

	switch col.DataType().ID() {
	case arrow.STRING:
		v := array.NewStringData(col.Data())
		for i := 0; i < col.Len(); i++ {
			if field.Nullable() {
				if v.IsNull(i) {
					var s *string
					field.Append(s)
					continue
				}
				s := v.Value(i)
				field.Append(&s)
				continue
			}
			field.Append(v.Value(i))
		}
	case arrow.UINT32:
		v := array.NewUint32Data(col.Data())
		for i := 0; i < col.Len(); i++ {
			if field.Nullable() {
				if v.IsNull(i) {
					var s *uint32
					field.Append(s)
					continue
				}
				s := v.Value(i)
				field.Append(&s)
				continue
			}
			field.Append(v.Value(i))
		}
	case arrow.INT64:
		v := array.NewInt64Data(col.Data())
		for i := 0; i < col.Len(); i++ {
			if field.Nullable() {
				if v.IsNull(i) {
					var s *int64
					field.Append(s)
					continue
				}
				s := v.Value(i)
				field.Append(&s)
				continue
			}
			field.Append(v.Value(i))
		}
	case arrow.FLOAT64:
		v := array.NewFloat64Data(col.Data())
		for i := 0; i < col.Len(); i++ {
			if field.Nullable() {
				if v.IsNull(i) {
					var f *float64
					field.Append(f)
					continue
				}
				f := v.Value(i)
				field.Append(&f)
				continue
			}
			field.Append(v.Value(i))
		}
	case arrow.BOOL:
		v := array.NewBooleanData(col.Data())
		for i := 0; i < col.Len(); i++ {
			if field.Nullable() {
				if v.IsNull(i) {
					var b *bool
					field.Append(b)
					continue
				}
				b := v.Value(i)
				field.Append(&b)
				continue
			}
			field.Append(v.Value(i))
		}
	case arrow.TIMESTAMP:
		v := array.NewTimestampData(col.Data())
		for i := 0; i < col.Len(); i++ {
			if field.Nullable() {
				if v.IsNull(i) {
					var t *time.Time
					field.Append(t)
					continue
				}
				t := v.Value(i).ToTime(arrow.Nanosecond)
				field.Append(&t)
				continue
			}
			field.Append(v.Value(i).ToTime(arrow.Nanosecond))
		}
	case arrow.DENSE_UNION:
		v := array.NewDenseUnionData(col.Data())
		for i := 0; i < col.Len(); i++ {
			sc, err := scalar.GetScalar(v, i)
			if err != nil {
				return err
			}
			value := sc.(*scalar.DenseUnion).ChildValue()

			var data any
			switch value.DataType().ID() {
			case arrow.STRING:
				data = value.(*scalar.String).String()
			case arrow.BOOL:
				data = value.(*scalar.Boolean).Value
			case arrow.INT32:
				data = value.(*scalar.Int32).Value
			case arrow.INT64:
				data = value.(*scalar.Int64).Value
			case arrow.LIST:
				data = value.(*scalar.List).Value
			}
			b, err := json.Marshal(data)
			if err != nil {
				return err
			}
			field.Append(json.RawMessage(b))
		}
	}

	return nil
}
