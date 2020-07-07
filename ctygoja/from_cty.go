package ctygoja

import (
	"fmt"
	"math/big"

	"github.com/dop251/goja"
	"github.com/zclconf/go-cty/cty"
)

// FromCtyValue takes a cty.Value and returns the equivalent goja.Value
// belonging to the given goja Runtime.
//
// Only known values can be converted to goja.Value. If you pass an unknown
// value then this function will panic. This function cannot convert
// capsule-typed values and will panic if you pass one.
//
// This function must not be called concurrently with other use of the given
// runtime.
func FromCtyValue(v cty.Value, js *goja.Runtime) goja.Value {
	ty := v.Type()
	switch {
	case !v.IsKnown():
		panic("ctygoja.FromCtyValue on unknown value")
	case v.IsNull():
		return goja.Null()
	case ty.IsObjectType() || ty.IsMapType():
		return fromCtyValueObject(v, js)
	default:
		raw := fromCtyValueReflect(v, js)
		return js.ToValue(raw)
	}
}

func fromCtyValueReflect(v cty.Value, js *goja.Runtime) interface{} {
	ty := v.Type()
	switch {
	case ty == cty.Bool:
		return v.True()
	case ty == cty.String:
		return v.AsString()
	case ty == cty.Number:
		// The JavaScript engine only has integer and float64 types, so
		// this can potentially be lossy.
		raw := v.AsBigFloat()
		if rawInt64, acc := raw.Int64(); acc == big.Exact {
			return rawInt64
		}
		rawFloat, _ := raw.Float64()
		return rawFloat
	case ty.IsListType() || ty.IsSetType() || ty.IsTupleType():
		raw := make([]interface{}, 0, v.LengthInt())
		for it := v.ElementIterator(); it.Next(); {
			_, v := it.Element()
			gojaV := FromCtyValue(v, js)
			raw = append(raw, gojaV)
		}
		return raw
	default:
		panic(fmt.Sprintf("ctygoja.FromCtyValue doesn't know how to convert %#v", v))
	}
}

func fromCtyValueObject(v cty.Value, js *goja.Runtime) *goja.Object {
	ret := js.NewObject()
	for it := v.ElementIterator(); it.Next(); {
		k, v := it.Element()
		rawK := k.AsString()
		gojaV := FromCtyValue(v, js)
		ret.DefineDataProperty(rawK, gojaV, goja.FLAG_FALSE, goja.FLAG_FALSE, goja.FLAG_TRUE)
	}
	return ret
}
