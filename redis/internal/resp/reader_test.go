package resp

import (
	"bytes"
	"testing"
)

func TestReader_ReadSimpleString(t *testing.T) {
	input := "+OK\r\n"
	reader := NewReader(bytes.NewBufferString(input))
	val, err := reader.Read()
	if err != nil {
		t.Fatalf("Failed to read: %v", err)
	}
	if val.Type != "string" {
		t.Errorf("Expected string, got %v", val.Type)
	}
	if val.Str != "OK" {
		t.Errorf("Expected OK, got %v", val.Str)
	}
}

func TestReader_ReadBulkString(t *testing.T) {
	input := "$5\r\nhello\r\n"
	reader := NewReader(bytes.NewBufferString(input))
	val, err := reader.Read()
	if err != nil {
		t.Fatalf("Failed to read: %v", err)
	}
	if val.Type != "bulk" {
		t.Errorf("Expected bulk, got %v", val.Type)
	}
	if val.Bulk != "hello" {
		t.Errorf("Expected hello, got %v", val.Bulk)
	}
}

func TestReader_ReadArray(t *testing.T) {
	input := "*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n"
	reader := NewReader(bytes.NewBufferString(input))
	val, err := reader.Read()
	if err != nil {
		t.Fatalf("Failed to read: %v", err)
	}

	if val.Type != "array" {
		t.Errorf("Expected array, got %v", val.Type)
	}

	if len(val.Array) != 2 {
		t.Errorf("Expected array length 2, got %d", len(val.Array))
	}

	if val.Array[0].Bulk != "hello" {
		t.Errorf("Expected first element hello, got %v", val.Array[0].Bulk)
	}
	if val.Array[1].Bulk != "world" {
		t.Errorf("Expected second element world, got %v", val.Array[1].Bulk)
	}
}
