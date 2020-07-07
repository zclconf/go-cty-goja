package ctygoja

import (
	"testing"

	"github.com/dop251/goja"
	"github.com/zclconf/go-cty/cty"
)

func TestFromCtyValue(t *testing.T) {
	js := goja.New()

	tests := []struct {
		given cty.Value
		test  string
	}{
		{
			cty.NullVal(cty.DynamicPseudoType),
			`if (v !== null) throw new Error('want null, but got '+v)`,
		},
		{
			cty.StringVal("hello"),
			`if (v !== 'hello') throw new Error('want "hello", but got '+v)`,
		},
		{
			cty.True,
			`if (!v) throw new Error('want true, but got '+v)`,
		},
		{
			cty.False,
			`if (v) throw new Error('want false, but got '+v)`,
		},
		{
			cty.NumberIntVal(0),
			`if (v !== 0) throw new Error('want 0, but got '+v)`,
		},
		{
			cty.NumberIntVal(1),
			`if (v !== 1) throw new Error('want 1, but got '+v)`,
		},
		{
			cty.NumberFloatVal(1.5),
			`if (v !== 1.5) throw new Error('want 1.5, but got '+v)`,
		},
		{
			cty.ObjectVal(map[string]cty.Value{
				"name": cty.StringVal("Ermintrude"),
				"foo":  cty.True,
			}),
			`
			if (v.name !== 'Ermintrude') throw new Error('wrong name');
			if (v.foo !== true) throw new Error('wrong foo');
			var keys = [];
			for (k in v) {
				keys.push(k);
			}
			if (keys.length != 2 || keys[0] != 'foo' || keys[1] != 'name')
				throw new Error('wrong keys')
			`,
		},
		{
			cty.MapVal(map[string]cty.Value{
				"name": cty.StringVal("Ermintrude"),
			}),
			`
			if (v.name !== 'Ermintrude') throw new Error('wrong name');
			var keys = [];
			for (k in v) {
				keys.push(k);
			}
			if (keys.length != 1 || keys[0] != 'name')
				throw new Error('wrong keys')
			`,
		},
		{
			cty.EmptyObjectVal,
			`
			if (JSON.stringify(v) != '{}') throw new Error('wrong result');
			`,
		},
		{
			cty.MapValEmpty(cty.String),
			`
			if (JSON.stringify(v) != '{}') throw new Error('wrong result');
			`,
		},
		{
			cty.TupleVal([]cty.Value{cty.True, cty.False}),
			`
			if (JSON.stringify(v) != '[true,false]') throw new Error('wrong result');
			`,
		},
		{
			cty.ListVal([]cty.Value{cty.True, cty.False}),
			`
			if (JSON.stringify(v) != '[true,false]') throw new Error('wrong result');
			`,
		},
		{
			cty.SetVal([]cty.Value{cty.StringVal("b"), cty.StringVal("a")}),
			`
			if (JSON.stringify(v) != '["a","b"]') throw new Error('wrong result');
			`,
		},
		{
			cty.EmptyTupleVal,
			`
			if (JSON.stringify(v) != '[]') throw new Error('wrong result');
			`,
		},
		{
			cty.ListValEmpty(cty.String),
			`
			if (JSON.stringify(v) != '[]') throw new Error('wrong result');
			`,
		},
		{
			cty.SetValEmpty(cty.String),
			`
			if (JSON.stringify(v) != '[]') throw new Error('wrong result');
			`,
		},
	}
	for _, test := range tests {
		t.Run(test.given.GoString(), func(t *testing.T) {
			got := FromCtyValue(test.given, js)
			testJS := goja.New()
			testJS.Set("v", got)
			_, err := testJS.RunString(test.test)
			if err != nil {
				gotObj := got.ToObject(js)
				repr, jsonErr := gotObj.MarshalJSON()
				if jsonErr != nil {
					repr = []byte(got.String())
				}
				t.Errorf("assertion failed\nGot:   %s\n%s", repr, err.Error())
			}
		})
	}
}
