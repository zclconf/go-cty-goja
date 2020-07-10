package ctygoja

import (
	"github.com/dop251/goja"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/json"
)

// ToCtyValue attempts to find a cty.Value that is equivalent to the given
// goja.Value, returning an error if no conversion is possible.
//
// Although cty is a superset of JSON and thus all cty values can be converted
// to JavaScript by way of a JSON-like mapping, JavaScript's type system
// includes many types that have no equivalent in cty, such as functions.
//
// For predictability and consistency, the conversion from JavaScript to cty
// is defined as a conversion from JavaScript to JSON using the same rules
// as JavaScript's JSON.stringify function, followed by interpretation of that
// result in cty using the same rules as the cty/json package follows.
//
// This function therefore fails in the cases where JSON.stringify would fail.
// Because neither cty nor JSON have an equivalent of "undefined", in cases
// where JSON.stringify would return undefined ToCtyValue returns a cty
// null value.
func ToCtyValue(v goja.Value, js *goja.Runtime) (cty.Value, error) {
	// There are some exceptions for things that can't be turned into a
	// goja.Object, because they don't have associated boxing prototypes.
	if goja.IsNull(v) || goja.IsUndefined(v) {
		return cty.NullVal(cty.DynamicPseudoType), nil
	}

	// For now at least, the implementation is literally to go via JSON
	// encoding, because goja offers a convenient interface to the same
	// behavior as JSON.stringify.
	src, err := v.ToObject(js).MarshalJSON()
	if err != nil {
		return cty.NilVal, err
	}

	ty, err := json.ImpliedType(src)
	if err != nil {
		// It'd be weird to end up here because that would suggest that
		// goja's MarshalJSON produced an invalid result, but we'll return
		// it out anyway.
		return cty.NilVal, err
	}

	return json.Unmarshal(src, ty)
}
