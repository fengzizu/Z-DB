package resp

// RESP Type Constants
// These prefixes determine the data type in the RESP protocol.
// RESP 类型常量
// 这些前缀决定了 RESP 协议中的数据类型。
const (
	STRING  = '+' // Simple String (e.g., "+OK\r\n") / 简单字符串
	ERROR   = '-' // Error (e.g., "-Error message\r\n") / 错误
	INTEGER = ':' // Integer (e.g., ":1000\r\n") / 整数
	BULK    = '$' // Bulk String (e.g., "$5\r\nhello\r\n") / 定长字符串（用于二进制安全数据）
	ARRAY   = '*' // Array (e.g., "*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n") / 数组
)

// Value represents a decoded RESP value.
// It acts as a container for any type of data transferred in Redis protocol.
// Value 结构体表示一个解码后的 RESP 值。
// 它是 Redis 协议传输中任意类型数据的容器。
type Value struct {
	Type  string  // Type of value: "string", "error", "integer", "bulk", "array" / 值类型
	Str   string  // Holds value for Simple Strings and Errors / 用于存储简单字符串和错误信息
	Num   int     // Holds value for Integers / 用于存储整数
	Bulk  string  // Holds value for Bulk Strings / 用于存储定长字符串内容
	Array []Value // Holds value for Arrays / 用于存储数组元素
}
