package ctygoja

import (
	"testing"

	"github.com/dop251/goja"
	"github.com/zclconf/go-cty/cty"
)

func TestToCtyValue(t *testing.T) {
	tests := []struct {
		Src  string
		Want cty.Value
		Err  bool
	}{
		{
			Src:  "null",
			Want: cty.NullVal(cty.DynamicPseudoType),
		},
		{
			Src:  "undefined",
			Want: cty.NullVal(cty.DynamicPseudoType),
		},
		{
			Src:  "12",
			Want: cty.NumberIntVal(12),
		},
		{
			Src:  "12.5",
			Want: cty.NumberFloatVal(12.5),
		},
		{
			Src:  "true",
			Want: cty.True,
		},
		{
			Src:  "false",
			Want: cty.False,
		},
		{
			Src:  `""`,
			Want: cty.StringVal(""),
		},
		{
			Src:  `"hello"`,
			Want: cty.StringVal("hello"),
		},
		{
			Src:  `({})`,
			Want: cty.EmptyObjectVal,
		},
		{
			Src: `({a: "b"})`,
			Want: cty.ObjectVal(map[string]cty.Value{
				"a": cty.StringVal("b"),
			}),
		},
		{
			Src:  `[]`,
			Want: cty.EmptyTupleVal,
		},
		{
			Src: `[true]`,
			Want: cty.TupleVal([]cty.Value{
				cty.True,
			}),
		},
		{
			Src:  `(function () {})`,
			Want: cty.NullVal(cty.DynamicPseudoType),
		},
		{
			Src: `new Date(0)`,
			// The Date prototype includes a toJSON function which
			// produces a timestamp string.
			Want: cty.StringVal("1970-01-01T00:00:00.000Z"),
		},
		{
			Src: `JSON`,
			// The JSON object has no enumerable properties, so it appears
			// as an empty object after conversion.
			Want: cty.EmptyObjectVal,
		},
		{
			Src:  "NaN",
			Want: cty.NullVal(cty.DynamicPseudoType),
		},
		{
			Src: "Infinity",
			// Even though cty can represent positive infinity, JSON doesn't
			// and our mapping is via JSON and so the result is null.
			Want: cty.NullVal(cty.DynamicPseudoType),
		},
	}

	for _, test := range tests {
		t.Run(test.Src, func(t *testing.T) {
			js := goja.New()
			result, err := js.RunString(test.Src)
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			got, gotErr := ToCtyValue(result, js)
			if test.Err {
				if gotErr == nil {
					t.Errorf("wrong result\ngot:  %#v\nwant: (error)", got)
				}
				return
			}

			if gotErr != nil {
				t.Fatalf("unexpected error\ngot:  %s\nwant: %#v", gotErr, test.Want)
			}

			if !test.Want.RawEquals(got) {
				t.Errorf("wrong result\ngot:  %#v\nwant: %#v", got, test.Want)
			}
		})
	}
}
