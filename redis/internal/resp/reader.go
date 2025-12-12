package resp

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

// Reader wraps a bufio.Reader to provide RESP parsing capabilities.
// Reader 包装了 bufio.Reader 以提供 RESP 协议解析能力。
type Reader struct {
	r *bufio.Reader
}

// NewReader creates a new RESP reader from an io.Reader.
// NewReader 从 io.Reader 创建一个新的 RESP 读取器。
func NewReader(rd io.Reader) *Reader {
	return &Reader{r: bufio.NewReader(rd)}
}

// ReadLine reads bytes until CRLF ('\r\n') and returns the line without CRLF.
// Used for parsing simple strings, errors, and integers.
// ReadLine 读取字节直到遇到 CRLF ('\r\n')，并返回不带 CRLF 的行内容。
// 用于解析简单字符串、错误和整数。
func (r *Reader) ReadLine() (line []byte, n int, err error) {
	for {
		b, err := r.r.ReadByte()
		if err != nil {
			return nil, 0, err
		}
		n += 1
		line = append(line, b)
		// Check for \r\n (CRLF) suffix / 检查是否以 \r\n 结尾
		if len(line) >= 2 && line[len(line)-2] == '\r' {
			break
		}
	}
	// Return line excluding the last 2 bytes (\r\n) / 返回去掉最后两个字节 (\r\n) 的内容
	return line[:len(line)-2], n, nil
}

// ReadInteger reads a line and parses it as an integer.
// ReadInteger 读取一行并将其解析为整数。
func (r *Reader) ReadInteger() (x int, n int, err error) {
	line, n, err := r.ReadLine()
	if err != nil {
		return 0, 0, err
	}
	// Parse base-10 integer / 解析十进制整数
	i64, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, n, err
	}
	return int(i64), n, nil
}

// Read reads the next full RESP value from the stream.
// Read 从流中读取下一个完整的 RESP 值。
func (r *Reader) Read() (Value, error) {
	// Read the first byte to determine type / 读取第一个字节以确定类型
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
		// Handle unknown types strictly for learning purposes / 为了学习目的严格处理未知类型
		fmt.Printf("Unknown type: %q (%d)\n", string(_type), _type)
		return Value{}, fmt.Errorf("unknown type: %v", string(_type))
	}
}

// readArray parses a RESP Array: *<len>\r\n...
// readArray 解析 RESP 数组：格式为 *<长度>\r\n...
func (r *Reader) readArray() (Value, error) {
	v := Value{}
	v.Type = "array"

	// 1. Read the length of the array / 读取数组长度
	len, _, err := r.ReadInteger()
	if err != nil {
		return v, err
	}

	// 2. Loop and read each element recursively / 循环递归读取每个元素
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

// readBulk parses a RESP Bulk String: $<len>\r\n<content>\r\n
// readBulk 解析 RESP 定长字符串：格式为 $<长度>\r\n<内容>\r\n
func (r *Reader) readBulk() (Value, error) {
	v := Value{}
	v.Type = "bulk"

	// 1. Read the length of the string / 读取字符串长度
	len, _, err := r.ReadInteger()
	if err != nil {
		return v, err
	}

	// 2. Read exact bytes based on length / 根据长度读取精确的字节数
	bulk := make([]byte, len)
	_, err = io.ReadFull(r.r, bulk)
	if err != nil {
		return v, err
	}

	v.Bulk = string(bulk)

	// 3. Read and discard the trailing CRLF / 读取并丢弃末尾的 CRLF
	r.ReadLine()

	return v, nil
}

// readString parses a RESP Simple String: +<content>\r\n
// readString 解析 RESP 简单字符串：格式为 +<内容>\r\n
func (r *Reader) readString() (Value, error) {
	v := Value{}
	v.Type = "string"

	// Read until CRLF / 读取直到 CRLF
	line, _, err := r.ReadLine()
	if err != nil {
		return v, err
	}

	v.Str = string(line)
	return v, nil
}
