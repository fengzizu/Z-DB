package resp

import (
	"bytes"
	"testing"
)

func TestWriter_Write(t *testing.T) {
	tests := []struct {
		name     string
		value    Value
		expected string
	}{
		{
			name:     "Simple String",
			value:    Value{Type: "string", Str: "OK"},
			expected: "+OK\r\n",
		},
		{
			name:     "Error",
			value:    Value{Type: "error", Str: "Error"},
			expected: "-Error\r\n",
		},
		{
			name:     "Integer",
			value:    Value{Type: "integer", Num: 123},
			expected: ":123\r\n",
		},
		{
			name:     "Bulk String",
			value:    Value{Type: "bulk", Bulk: "hello"},
			expected: "$5\r\nhello\r\n",
		},
		{
			name:     "Null",
			value:    Value{Type: "null"},
			expected: "$-1\r\n",
		},
		{
			name: "Array",
			value: Value{
				Type: "array",
				Array: []Value{
					{Type: "bulk", Bulk: "hello"},
					{Type: "bulk", Bulk: "world"},
				},
			},
			expected: "*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			writer := NewWriter(buf)
			err := writer.Write(tt.value)
			if err != nil {
				t.Errorf("Write() error = %v", err)
				return
			}
			if got := buf.String(); got != tt.expected {
				t.Errorf("Write() = %q, want %q", got, tt.expected)
			}
		})
	}
}

