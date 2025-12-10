package resp

// RESP Type Constants
const (
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
)

// Value represents a decoded RESP value
type Value struct {
	Type  string  // "string", "error", "integer", "bulk", "array"
	Str   string  // For Simple Strings, Errors
	Num   int     // For Integers
	Bulk  string  // For Bulk Strings
	Array []Value // For Arrays
}
