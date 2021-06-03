package parser

import "testing"

func TestGetType(t *testing.T) {
	parser := NewParser()
	if parser.GetType('+') != StringType {
		t.Error("Expected a string type to be returned")
	}
	if parser.GetType('-') != ErrorType {
		t.Error("Expected an error type to be returned")
	}
	if parser.GetType(':') != IntegerType {
		t.Error("Expected a integer type to be returned")
	}
	if parser.GetType('$') != BulkStringType {
		t.Error("Expected a bulk string type to be returned")
	}
	if parser.GetType('*') != ArrayType {
		t.Error("Expected an array type to be returned")
	}
	if parser.GetType('@') != UnknownType {
		t.Error("Expected an unknown type to be returned")
	}
}

func TestReadIntResult(t *testing.T) {
	parser := NewParser()
	intValue, err := parser.ReadInt(":14\r\n")
	if err != nil {
		t.Error(err.Error())
	}
	if intValue != 14 {
		t.Error("int value should be equal to 14")
	}
}

func TestReadString(t *testing.T) {
	parser := NewParser()
	strValue, err := parser.ReadString("+hello world\r\n")
	if err != nil {
		t.Error(err.Error())
	}
	if strValue != "hello world" {
		t.Error("str value should be equal to hello world")
	}
}

func TestReadBulkString(t *testing.T) {
	parser := NewParser()
	strValue, err := parser.ReadBulkString("$6\r\nfoobar\r\n")
	if err != nil {
		t.Error(err.Error())
	}
	if strValue != nil && *strValue != "foobar" {
		t.Error("str value should be equal to foobar")
	}
}

func TestReadNullBulkString(t *testing.T) {
	parser := NewParser()
	strValue, err := parser.ReadBulkString("$-1\r\n")
	if err != nil {
		t.Error(err.Error())
	}
	if strValue != nil {
		t.Error("str value should be a nil string")
	}
}

func TestReadEmptyString(t *testing.T) {
	parser := NewParser()
	strValue, err := parser.ReadBulkString("$0\r\n\r\n")
	if err != nil {
		t.Error(err.Error())
	}
	if strValue != nil && *strValue != "" {
		t.Error("str value should be a nil string")
	}
}

func TestReadArray(t *testing.T) {
	parser := NewParser()
	result, err := parser.ReadArray("*3\r\n:1\r\n:2\r\n:3\r\n")
	if err != nil {
		t.Error(err)
	}
	if result == nil {
		t.Error("OK")
	}
}

func BenchmarkReadIntResult(b *testing.B) {
	parser := NewParser()
	for i := 0; i < b.N; i++ {
		parser.ReadInt(":2567065")
	}
	b.ReportAllocs()
}

func BenchmarkReadString(b *testing.B) {
	parser := NewParser()
	for i := 0; i < b.N; i++ {
		parser.ReadString("+Hello world Test")
	}
	b.ReportAllocs()
}

func BenchmarkReadBulkString(b *testing.B) {
	parser := NewParser()
	for i := 0; i < b.N; i++ {
		parser.ReadBulkString("$6\r\nfoobar\r\n")
	}
	b.ReportAllocs()
}

func BenchmarkReadNullBulkString(b *testing.B) {
	parser := NewParser()
	for i := 0; i < b.N; i++ {
		parser.ReadBulkString("$-1\r\n")
	}
	b.ReportAllocs()
}

func BenchmarkReadEmptyString(b *testing.B) {
	parser := NewParser()
	for i := 0; i < b.N; i++ {
		parser.ReadBulkString("$0\r\n\r\n")
	}
	b.ReportAllocs()
}

func BenchmarkReadArray(b *testing.B) {
	parser := NewParser()
	for i := 0; i < b.N; i++ {
		parser.ReadArray("*4\r\n:1\r\n:2\r\n:3\r\n:4\r\n")
	}
	b.ReportAllocs()
}
