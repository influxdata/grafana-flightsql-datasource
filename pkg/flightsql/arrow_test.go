package flightsql

import (
	"log"
	"testing"
	"time"

	"github.com/apache/arrow/go/v10/arrow"
	"github.com/apache/arrow/go/v10/arrow/array"
	"github.com/apache/arrow/go/v10/arrow/memory"
	"github.com/google/go-cmp/cmp"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/stretchr/testify/require"
)

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

	actual := newFrame(schema, "select * from data")
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
