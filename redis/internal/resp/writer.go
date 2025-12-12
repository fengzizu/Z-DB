package resp

import (
	"io"
	"strconv"
)

// Writer wraps an io.Writer to serialize RESP Values.
// Writer 包装 io.Writer 以序列化 RESP 值。
type Writer struct {
	w io.Writer
}

// NewWriter creates a new RESP writer.
// NewWriter 创建一个新的 RESP 写入器。
func NewWriter(w io.Writer) *Writer {
	return &Writer{w: w}
}

// Write serializes a Value struct into RESP bytes and writes it to the underlying writer.
// Write 将 Value 结构体序列化为 RESP 字节流并写入底层 writer。
func (w *Writer) Write(v Value) error {
	var bytes []byte

	switch v.Type {
	case "string":
		// Simple String: +content\r\n
		bytes = []byte("+" + v.Str + "\r\n")
	case "error":
		// Error: -message\r\n
		bytes = []byte("-" + v.Str + "\r\n")
	case "integer":
		// Integer: :number\r\n
		bytes = []byte(":" + strconv.Itoa(v.Num) + "\r\n")
	case "bulk":
		// Bulk String: $<len>\r\n<content>\r\n
		bytes = []byte("$" + strconv.Itoa(len(v.Bulk)) + "\r\n" + v.Bulk + "\r\n")
	case "null":
		// Null Bulk String (nil): $-1\r\n
		bytes = []byte("$-1\r\n")
	case "array":
		// Array: *<len>\r\n...elements...
		bytes = []byte("*" + strconv.Itoa(len(v.Array)) + "\r\n")
		// Write array header first / 先写入数组头
		if _, err := w.w.Write(bytes); err != nil {
			return err
		}
		// Write each element recursively / 递归写入每个元素
		for _, val := range v.Array {
			if err := w.Write(val); err != nil {
				return err
			}
		}
		return nil // Array finished
	default:
		// Default to Null if type is unknown / 如果类型未知，默认为 Null
		bytes = []byte("$-1\r\n")
	}

	_, err := w.w.Write(bytes)
	return err
}
