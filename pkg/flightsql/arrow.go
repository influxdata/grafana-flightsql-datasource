package flightsql

import (
	"errors"
	"fmt"
	"io"
	"runtime/debug"
	"time"

	"github.com/apache/arrow/go/v10/arrow"
	"github.com/apache/arrow/go/v10/arrow/array"
	"github.com/apache/arrow/go/v10/arrow/flight"
	"github.com/apache/arrow/go/v10/arrow/memory"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

// newFrame creates a new Grafana Frame who's schema matches an Arrow Schema for
// a specific SQL query statement. The SQL statement provided populates the
// Frame's DataTopic.
func newFrame(schema *arrow.Schema, sql string) *data.Frame {
	log.DefaultLogger.Info(fmt.Sprintf(
		"Schema: metadata=%v fields=%v",
		schema.Metadata(),
		schema.Fields(),
	))

	fields := schema.Fields()
	df := &data.Frame{
		Fields: make([]*data.Field, len(fields)),
		Meta: &data.FrameMeta{
			ExecutedQueryString: sql,
			DataTopic:           data.DataTopic(sql),
		},
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
		}
	}
	return df
}

// copyData copies the data from an Arrow Array into a Grafana Field.
func copyData(field *data.Field, col arrow.Array) {
	defer func() {
		if r := recover(); r != nil {
			log.DefaultLogger.Error("Panic", r, string(debug.Stack()))
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
	case arrow.INT64:
		v := array.NewInt64Data(col.Data())
		for i := 0; i < col.Len(); i++ {
			if field.Nullable() {
				if v.IsNull(i) {
					var s *int64
					field.Append(s)
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
	}
}

// readTableSchemas scans a flight.Reader from CommandGetTables call that
// includes the schema information and produces a map of Table Name -> Arrow
// Schema.
func readTableSchemas(reader *flight.Reader) (map[string]*arrow.Schema, error) {
	schemas := map[string]*arrow.Schema{}

	if !reader.Next() {
		return nil, fmt.Errorf("no records in the response")
	}

	record := reader.Record()
	record.Retain()
	defer record.Release()

	namesIdx := record.Schema().FieldIndices("table_name")
	if len(namesIdx) == 0 {
		return nil, fmt.Errorf("table_name field not found")
	}
	schemasIdx := record.Schema().FieldIndices("table_schema")
	if len(schemasIdx) == 0 {
		return nil, fmt.Errorf("table_schema field not found")
	}

	schemaCol := record.Column(schemasIdx[0])
	nameCol := record.Column(namesIdx[0])
	for i := 0; i < schemaCol.Len(); i++ {
		name := array.NewStringData(nameCol.Data()).Value(i)
		serializedSchema := array.NewStringData(schemaCol.Data()).Value(i)
		schema, err := flight.DeserializeSchema([]byte(serializedSchema), memory.DefaultAllocator)
		if err != nil {
			return nil, fmt.Errorf("deserialize schema: %w", err)
		}
		schemas[name] = schema
	}

	if err := reader.Err(); err != nil {
		return nil, err
	}

	return schemas, nil
}

// newDataResponse creates a backend.DataResponse (with a single Frame) from the
// contents of an Arrow Table obtained by scanning a flight.Reader.
func newDataResponse(reader *flight.Reader) backend.DataResponse {
	var resp backend.DataResponse
	frame := newFrame(reader.Schema(), "")
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
