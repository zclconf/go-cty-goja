// Package ctygoja is an adapter layer for converting between cty values and
// values used by the pure-Go JavaScript engine "goja". This can be used to
// create interfaces between a cty-based application and a JavaScript
// context.
//
// The philosophy for this package is to behave as if values were converted
// via JSON using the mappings in the ctyjson package, relying on the fact
// that the cty type system (capsule types notwithstanding) is a superset
// of JSON.
package ctygoja
