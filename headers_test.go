package headers

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

type testStruct struct {
	String  string       `header:"x-string"`
	Slices  []string     `header:"x-slice"`
	Int     int          `header:"x-int"`
	Int8    int8         `header:"x-int8"`
	Int16   int16        `header:"x-int16"`
	Int32   int32        `header:"x-int32"`
	Int64   int64        `header:"x-int64"`
	Uint    uint         `header:"x-uint"`
	Uint8   uint8        `header:"x-uint8"`
	Uint16  uint16       `header:"x-uint16"`
	Uint32  uint32       `header:"x-uint32"`
	Uint64  uint64       `header:"x-uint64"`
	Float32 float32      `header:"x-float32"`
	Float64 float64      `header:"x-float64"`
	Bool    bool         `header:"x-bool"`
	Ints    CustomInts   `header:"x-ints"`
	Floats  CustomFloats `header:"x-floats"`

	SkipKey struct{} `header:"-"`
	SkipIt  struct{}
}

type CustomInts []int

func (ci *CustomInts) UnmarshalHeader(values []string) error {
	for i := 0; i < len(values); i++ {
		tmp, err := strconv.Atoi(values[i])
		if err != nil {
			return err
		}

		*ci = append(*ci, tmp)
	}

	return nil
}

type CustomFloats []float64

func (cf *CustomFloats) UnmarshalHeader(values []string) error {
	for i := 0; i < len(values); i++ {
		tmp, err := strconv.ParseFloat(values[i], 64)
		if err != nil {
			return err
		}

		*cf = append(*cf, tmp)
	}

	return nil
}

func setTestHeaders() http.Header {
	out := make(http.Header)

	out.Set("x-string", "some-string")
	out.Add("x-slice", "some")
	out.Add("x-slice", "slice")
	out.Set("x-int", "-1")
	out.Set("x-int8", "-2")
	out.Set("x-int16", "-3")
	out.Set("x-int32", "-4")
	out.Set("x-int64", "-5")
	out.Set("x-uint", "1")
	out.Set("x-uint8", "2")
	out.Set("x-uint16", "3")
	out.Set("x-uint32", "4")
	out.Set("x-uint64", "5")
	out.Set("x-float32", "0.123")
	out.Set("x-float64", "1.234")
	out.Set("x-bool", "true")

	for i := 0; i < 10; i++ {
		out.Add("x-ints", strconv.Itoa(i))
	}

	for i := 0; i < 10; i++ {
		out.Add("x-floats", strconv.Itoa(i))
	}

	return out
}

func TestUnknownType(t *testing.T) {
	var out struct {
		UnknownType [3]byte `header:"x-unknown-type"`
	}

	hdr := http.Header{"X-Unknown-Type": []string{"xxx"}}

	require.Equal(t, ErrUnknownType, UnmarshalHeaders(&out, hdr))
}

func TestUnmarshalHeaders(t *testing.T) {
	cases := []struct {
		name string

		val testStruct
		inp http.Header

		exp testStruct
		err error
	}{
		{name: "empty"},
		{name: "multiple", val: testStruct{}, inp: setTestHeaders(), exp: testStruct{
			String:  "some-string",
			Slices:  []string{"some", "slice"},
			Int:     -1,
			Int8:    -2,
			Int16:   -3,
			Int32:   -4,
			Int64:   -5,
			Uint:    1,
			Uint8:   2,
			Uint16:  3,
			Uint32:  4,
			Uint64:  5,
			Float32: 0.123,
			Float64: 1.234,
			Bool:    true,
			Ints:    CustomInts{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			Floats:  CustomFloats{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		}},
		{name: "int-errors", inp: http.Header{"X-Int": []string{"bad"}},
			err: &strconv.NumError{Func: "ParseInt", Num: "bad", Err: strconv.ErrSyntax}},
		{name: "uint-errors", inp: http.Header{"X-Uint": []string{"bad"}},
			err: &strconv.NumError{Func: "ParseUint", Num: "bad", Err: strconv.ErrSyntax}},
		{name: "float-errors", inp: http.Header{"X-Float32": []string{"bad"}},
			err: &strconv.NumError{Func: "ParseFloat", Num: "bad", Err: strconv.ErrSyntax}},
		{name: "bool-errors", inp: http.Header{"X-Bool": []string{"bad"}},
			err: &strconv.NumError{Func: "ParseBool", Num: "bad", Err: strconv.ErrSyntax}},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.err, UnmarshalHeaders(&tt.val, tt.inp))
			require.Equal(t, tt.exp, tt.val)
		})
	}
}
