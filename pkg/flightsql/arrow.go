package flightsql

import (
	"runtime/debug"
	"time"

	"github.com/apache/arrow/go/v10/arrow"
	"github.com/apache/arrow/go/v10/arrow/array"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

func newFrame(schema *arrow.Schema) *data.Frame {
	// TODO(brett): Make use of schema.Metadata.
	log.DefaultLogger.Info("Metadata", schema.Metadata())

	fields := schema.Fields()
	df := &data.Frame{
		Fields: make([]*data.Field, len(fields)),
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
