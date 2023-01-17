package flightsql

import (
	"fmt"
	"strings"
	"testing"

	"github.com/apache/arrow/go/v10/arrow"
	"github.com/apache/arrow/go/v10/arrow/array"
	"github.com/apache/arrow/go/v10/arrow/memory"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewQueryDataResponse(t *testing.T) {
	alloc := memory.DefaultAllocator
	schema := arrow.NewSchema(
		[]arrow.Field{
			{Name: "f1-i64", Type: arrow.PrimitiveTypes.Int64},
			{Name: "f2-f64", Type: arrow.PrimitiveTypes.Float64},
		},
		nil,
	)

	i64s, _, err := array.FromJSON(
		alloc,
		&arrow.Int64Type{},
		strings.NewReader(`[1, 2, 3]`),
	)
	require.NoError(t, err)
	f64s, _, err := array.FromJSON(
		alloc,
		&arrow.Float64Type{},
		strings.NewReader(`[1.1, 2.2, 3.3]`),
	)
	require.NoError(t, err)

	record := array.NewRecord(schema, []arrow.Array{i64s, f64s}, -1)
	records := []arrow.Record{record}
	reader, err := array.NewRecordReader(schema, records)
	require.NoError(t, err)

	resp := newQueryDataResponse(errReader{RecordReader: reader}, "")
	require.NoError(t, resp.Error)
	require.Len(t, resp.Frames, 1)
	require.Len(t, resp.Frames[0].Fields, 2)

	frame := resp.Frames[0]
	field0 := frame.Fields[0]
	assert.Equal(t, field0.Name, "f1-i64")
	assert.Equal(t, field0.Type(), data.FieldTypeInt64)
	assert.Equal(t, []any{int64(1), int64(2), int64(3)}, extractFieldValues(field0))

	field1 := frame.Fields[1]
	assert.Equal(t, field1.Name, "f2-f64")
	assert.Equal(t, field1.Type(), data.FieldTypeFloat64)
	assert.Equal(t, []any{1.1, 2.2, 3.3}, extractFieldValues(field1))
}

func TestNewQueryDataResponse_Error(t *testing.T) {
	alloc := memory.DefaultAllocator
	schema := arrow.NewSchema(
		[]arrow.Field{
			{Name: "f1-i64", Type: arrow.PrimitiveTypes.Int64},
			{Name: "f2-f64", Type: arrow.PrimitiveTypes.Float64},
		},
		nil,
	)

	i64s, _, err := array.FromJSON(
		alloc,
		&arrow.Int64Type{},
		strings.NewReader(`[1, 2, 3]`),
	)
	require.NoError(t, err)
	f64s, _, err := array.FromJSON(
		alloc,
		&arrow.Float64Type{},
		strings.NewReader(`[1.1, 2.2, 3.3]`),
	)
	require.NoError(t, err)

	record := array.NewRecord(schema, []arrow.Array{i64s, f64s}, -1)
	records := []arrow.Record{record}
	reader, err := array.NewRecordReader(schema, records)
	require.NoError(t, err)

	wrappedReader := errReader{
		RecordReader: reader,
		err:          fmt.Errorf("explosion!"),
	}
	resp := newQueryDataResponse(wrappedReader, "")
	require.Error(t, resp.Error)
	require.Equal(t, fmt.Errorf("explosion!"), resp.Error)
}

func TestNewQueryDataResponse_WideTable(t *testing.T) {
	alloc := memory.DefaultAllocator
	schema := arrow.NewSchema(
		[]arrow.Field{
			{Name: "time", Type: &arrow.TimestampType{}},
			{Name: "label", Type: &arrow.StringType{}},
			{Name: "value", Type: arrow.PrimitiveTypes.Int64},
		},
		nil,
	)

	times, _, err := array.FromJSON(
		alloc,
		&arrow.TimestampType{},
		strings.NewReader(`["2023-01-01T00:00:00Z", "2023-01-01T00:00:01Z", "2023-01-01T00:00:02Z"]`),
	)
	require.NoError(t, err)
	strs, _, err := array.FromJSON(
		alloc,
		&arrow.StringType{},
		strings.NewReader(`["foo", "bar", "baz"]`),
	)
	require.NoError(t, err)
	i64s, _, err := array.FromJSON(
		alloc,
		arrow.PrimitiveTypes.Int64,
		strings.NewReader(`[1, 2, 3]`),
	)
	require.NoError(t, err)

	record := array.NewRecord(schema, []arrow.Array{times, strs, i64s}, -1)
	records := []arrow.Record{record}
	reader, err := array.NewRecordReader(schema, records)
	require.NoError(t, err)

	resp := newQueryDataResponse(errReader{RecordReader: reader}, "")
	require.NoError(t, resp.Error)
	require.Len(t, resp.Frames, 1)
	require.Equal(t, 3, resp.Frames[0].Rows())
	require.Len(t, resp.Frames[0].Fields, 4)

	frame := resp.Frames[0]
	assert.Equal(t, "time", frame.Fields[0].Name)

	// label=bar
	assert.Equal(t, "value", frame.Fields[1].Name)
	assert.Equal(t, data.Labels{"label": "bar"}, frame.Fields[1].Labels)
	assert.Equal(t, []any{int64(0), int64(2), int64(0)}, extractFieldValues(frame.Fields[1]))

	// label=baz
	assert.Equal(t, "value", frame.Fields[2].Name)
	assert.Equal(t, data.Labels{"label": "baz"}, frame.Fields[2].Labels)
	assert.Equal(t, []any{int64(0), int64(0), int64(3)}, extractFieldValues(frame.Fields[2]))

	// label=foo
	assert.Equal(t, "value", frame.Fields[3].Name)
	assert.Equal(t, data.Labels{"label": "foo"}, frame.Fields[3].Labels)
	assert.Equal(t, []any{int64(1), int64(0), int64(0)}, extractFieldValues(frame.Fields[3]))
}

func extractFieldValues(field *data.Field) []any {
	values := make([]any, 0, field.Len())
	for i := 0; i < cap(values); i++ {
		values = append(values, field.CopyAt(i))
	}
	return values
}

type errReader struct {
	array.RecordReader
	err error
}

func (r errReader) Err() error {
	return r.err
}
