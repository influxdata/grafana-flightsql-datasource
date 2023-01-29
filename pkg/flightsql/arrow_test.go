package flightsql

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/apache/arrow/go/v10/arrow"
	"github.com/apache/arrow/go/v10/arrow/array"
	"github.com/apache/arrow/go/v10/arrow/memory"
	"github.com/google/go-cmp/cmp"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/grafana-plugin-sdk-go/data/sqlutil"
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

	query := sqlutil.Query{Format: sqlutil.FormatOptionTable}
	resp := newQueryDataResponse(errReader{RecordReader: reader}, query)
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
	query := sqlutil.Query{Format: sqlutil.FormatOptionTable}
	resp := newQueryDataResponse(wrappedReader, query)
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

	resp := newQueryDataResponse(errReader{RecordReader: reader}, sqlutil.Query{})
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

func TestNewFrame(t *testing.T) {
	schema := arrow.NewSchema([]arrow.Field{
		{
			Name:     "name",
			Type:     &arrow.StringType{},
			Nullable: false,
			Metadata: arrow.NewMetadata(nil, nil),
		},
		{
			Name:     "time",
			Type:     &arrow.TimestampType{},
			Nullable: false,
			Metadata: arrow.NewMetadata(nil, nil),
		},
		{
			Name:     "extra",
			Type:     &arrow.Int64Type{},
			Nullable: true,
			Metadata: arrow.NewMetadata(nil, nil),
		},
	}, nil)

	actual := newFrame(schema)
	expected := &data.Frame{
		Fields: []*data.Field{
			data.NewField("name", nil, []string{}),
			data.NewField("time", nil, []time.Time{}),
			data.NewField("extra", nil, []*int64{}),
		},
	}
	if !cmp.Equal(expected, actual, cmp.Comparer(cmpFrame)) {
		log.Fatalf(cmp.Diff(expected, actual))
	}
}

func cmpFrame(a, b data.Frame) bool {
	if len(a.Fields) != len(b.Fields) {
		return false
	}
	for i := 0; i < len(a.Fields); i++ {
		if a.Fields[i].Name != b.Fields[i].Name {
			return false
		}
		if a.Fields[i].Nullable() != b.Fields[i].Nullable() {
			return false
		}
	}
	return true
}

func TestCopyData_String(t *testing.T) {
	field := data.NewField("field", nil, []string{})
	builder := array.NewStringBuilder(memory.DefaultAllocator)
	builder.Append("joe")
	builder.Append("john")
	builder.Append("jackie")
	copyData(field, builder.NewArray())
	require.Equal(t, "joe", field.CopyAt(0))
	require.Equal(t, "john", field.CopyAt(1))
	require.Equal(t, "jackie", field.CopyAt(2))

	field = data.NewField("field", nil, []*string{})
	builder = array.NewStringBuilder(memory.DefaultAllocator)
	builder.Append("joe")
	builder.AppendNull()
	builder.Append("jackie")
	copyData(field, builder.NewArray())
	require.Equal(t, "joe", *(field.CopyAt(0).(*string)))
	require.Equal(t, (*string)(nil), field.CopyAt(1))
	require.Equal(t, "jackie", *(field.CopyAt(2).(*string)))
}

func TestCopyData_Timestamp(t *testing.T) {
	start, _ := time.Parse(time.RFC3339, "2023-01-01T01:01:01Z")

	field := data.NewField("field", nil, []time.Time{})
	builder := array.NewTimestampBuilder(memory.DefaultAllocator, &arrow.TimestampType{})
	builder.Append(arrow.Timestamp(start.Add(time.Hour).UnixNano()))
	builder.Append(arrow.Timestamp(start.Add(2 * time.Hour).UnixNano()))
	builder.Append(arrow.Timestamp(start.Add(3 * time.Hour).UnixNano()))
	copyData(field, builder.NewArray())
	require.Equal(t, start.Add(time.Hour), field.CopyAt(0))
	require.Equal(t, start.Add(2*time.Hour), field.CopyAt(1))
	require.Equal(t, start.Add(3*time.Hour), field.CopyAt(2))

	field = data.NewField("field", nil, []*time.Time{})
	builder = array.NewTimestampBuilder(memory.DefaultAllocator, &arrow.TimestampType{})
	builder.Append(arrow.Timestamp(start.Add(time.Hour).UnixNano()))
	builder.AppendNull()
	builder.Append(arrow.Timestamp(start.Add(3 * time.Hour).UnixNano()))
	copyData(field, builder.NewArray())
	require.Equal(t, start.Add(time.Hour), *field.CopyAt(0).(*time.Time))
	require.Equal(t, (*time.Time)(nil), field.CopyAt(1))
	require.Equal(t, start.Add(3*time.Hour), *field.CopyAt(2).(*time.Time))
}

func TestCopyData_Boolean(t *testing.T) {
	field := data.NewField("field", nil, []bool{})
	builder := array.NewBooleanBuilder(memory.DefaultAllocator)
	builder.Append(true)
	builder.Append(false)
	builder.Append(true)
	copyData(field, builder.NewArray())
	require.Equal(t, true, field.CopyAt(0))
	require.Equal(t, false, field.CopyAt(1))
	require.Equal(t, true, field.CopyAt(2))

	field = data.NewField("field", nil, []*bool{})
	builder = array.NewBooleanBuilder(memory.DefaultAllocator)
	builder.Append(true)
	builder.AppendNull()
	builder.Append(true)
	copyData(field, builder.NewArray())
	require.Equal(t, true, *field.CopyAt(0).(*bool))
	require.Equal(t, (*bool)(nil), field.CopyAt(1))
	require.Equal(t, true, *field.CopyAt(2).(*bool))
}

func TestCopyData_Int64(t *testing.T) {
	field := data.NewField("field", nil, []int64{})
	builder := array.NewInt64Builder(memory.DefaultAllocator)
	builder.Append(1)
	builder.Append(2)
	builder.Append(3)
	copyData(field, builder.NewArray())
	require.Equal(t, int64(1), field.CopyAt(0))
	require.Equal(t, int64(2), field.CopyAt(1))
	require.Equal(t, int64(3), field.CopyAt(2))

	field = data.NewField("field", nil, []*int64{})
	builder = array.NewInt64Builder(memory.DefaultAllocator)
	builder.Append(1)
	builder.AppendNull()
	builder.Append(3)
	arr := builder.NewArray()
	copyData(field, arr)
	require.Equal(t, int64(1), *field.CopyAt(0).(*int64))
	require.Equal(t, (*int64)(nil), field.CopyAt(1))
	require.Equal(t, int64(3), *field.CopyAt(2).(*int64))
}

func TestCopyData_Float64(t *testing.T) {
	field := data.NewField("field", nil, []float64{})
	builder := array.NewFloat64Builder(memory.DefaultAllocator)
	builder.Append(1.1)
	builder.Append(2.2)
	builder.Append(3.3)
	copyData(field, builder.NewArray())
	require.Equal(t, float64(1.1), field.CopyAt(0))
	require.Equal(t, float64(2.2), field.CopyAt(1))
	require.Equal(t, float64(3.3), field.CopyAt(2))

	field = data.NewField("field", nil, []*float64{})
	builder = array.NewFloat64Builder(memory.DefaultAllocator)
	builder.Append(1.1)
	builder.AppendNull()
	builder.Append(3.3)
	copyData(field, builder.NewArray())
	require.Equal(t, float64(1.1), *field.CopyAt(0).(*float64))
	require.Equal(t, (*float64)(nil), field.CopyAt(1))
	require.Equal(t, float64(3.3), *field.CopyAt(2).(*float64))
}
