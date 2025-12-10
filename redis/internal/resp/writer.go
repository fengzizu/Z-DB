package resp

import (
	"io"
	"strconv"
)

type Writer struct {
	w io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{w: w}
}

func (w *Writer) Write(v Value) error {
	var bytes []byte

	switch v.Type {
	case "string":
		bytes = []byte("+" + v.Str + "\r\n")
	case "error":
		bytes = []byte("-" + v.Str + "\r\n")
	case "integer":
		bytes = []byte(":" + strconv.Itoa(v.Num) + "\r\n")
	case "bulk":
		bytes = []byte("$" + strconv.Itoa(len(v.Bulk)) + "\r\n" + v.Bulk + "\r\n")
	case "null":
		bytes = []byte("$-1\r\n")
	case "array":
		bytes = []byte("*" + strconv.Itoa(len(v.Array)) + "\r\n")
		if _, err := w.w.Write(bytes); err != nil {
			return err
		}
		for _, val := range v.Array {
			if err := w.Write(val); err != nil {
				return err
			}
		}
		return nil // Array parts written
	default:
		// Default to null or error?
		// Let's assume null for safety
		bytes = []byte("$-1\r\n")
	}

	_, err := w.w.Write(bytes)
	return err
}

