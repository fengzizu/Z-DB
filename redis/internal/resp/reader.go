package resp

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

type Reader struct {
	r *bufio.Reader
}

func NewReader(rd io.Reader) *Reader {
	return &Reader{r: bufio.NewReader(rd)}
}

func (r *Reader) ReadLine() (line []byte, n int, err error) {
	for {
		b, err := r.r.ReadByte()
		if err != nil {
			return nil, 0, err
		}
		n += 1
		line = append(line, b)
		if len(line) >= 2 && line[len(line)-2] == '\r' {
			break
		}
	}
	return line[:len(line)-2], n, nil
}

func (r *Reader) ReadInteger() (x int, n int, err error) {
	line, n, err := r.ReadLine()
	if err != nil {
		return 0, 0, err
	}
	i64, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, n, err
	}
	return int(i64), n, nil
}

func (r *Reader) Read() (Value, error) {
	_type, err := r.r.ReadByte()
	if err != nil {
		return Value{}, err
	}

	switch _type {
	case ARRAY:
		return r.readArray()
	case BULK:
		return r.readBulk()
	case STRING:
		return r.readString()
	default:
		return Value{}, fmt.Errorf("unknown type: %v", string(_type))
	}
}

func (r *Reader) readArray() (Value, error) {
	v := Value{}
	v.Type = "array"

	// Read length of array
	len, _, err := r.ReadInteger()
	if err != nil {
		return v, err
	}

	// foreach line, parse
	v.Array = make([]Value, 0)
	for i := 0; i < len; i++ {
		val, err := r.Read()
		if err != nil {
			return v, err
		}
		v.Array = append(v.Array, val)
	}

	return v, nil
}

func (r *Reader) readBulk() (Value, error) {
	v := Value{}
	v.Type = "bulk"

	len, _, err := r.ReadInteger()
	if err != nil {
		return v, err
	}

	bulk := make([]byte, len)
	_, err = io.ReadFull(r.r, bulk)
	if err != nil {
		return v, err
	}

	v.Bulk = string(bulk)

	// Read trailing CRLF
	r.ReadLine()

	return v, nil
}

func (r *Reader) readString() (Value, error) {
	v := Value{}
	v.Type = "string"

	line, _, err := r.ReadLine()
	if err != nil {
		return v, err
	}

	v.Str = string(line)
	return v, nil
}
