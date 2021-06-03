package parser

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
)

type DataType int

const (
	StringType DataType = iota
	ErrorType
	IntegerType
	BulkStringType
	ArrayType
	UnknownType
)

type RawResponse = string
type BoxedValue struct {
	resultType DataType
	value      interface{}
}

var emptyArrayResult = []BoxedValue{}

type Parser struct{}

func NewBoxedValue(resultType DataType, value interface{}) BoxedValue {
	return BoxedValue{
		resultType: resultType,
		value:      value,
	}
}

func (parser *Parser) GetType(rawType uint8) DataType {
	switch rawType {
	case '+':
		return StringType
	case '-':
		return ErrorType
	case ':':
		return IntegerType
	case '$':
		return BulkStringType
	case '*':
		return ArrayType
	default:
		return UnknownType
	}
}

func typeToString(resultType DataType) string {
	switch resultType {
	case IntegerType:
		return "IntegerType"
	case ErrorType:
		return "ErrorType"
	case StringType:
		return "StringType"
	case BulkStringType:
		return "BulkStringType"
	case ArrayType:
		return "ArrayType"
	default:
		return "UnknownType"
	}
}

func (parser *Parser) consumeContent(buffer *bytes.Buffer) (string, error) {
	content, err := buffer.ReadString('\r')
	if err != nil {
		return "", err
	}
	return content[0 : len(content)-1], err
}

func (parser *Parser) consumeInt(buffer *bytes.Buffer) (int, error) {
	content, err := parser.consumeContent(buffer)
	if err != nil {
		return -1, err
	}
	intValue, err := strconv.Atoi(content)
	if err != nil {
		return -1, fmt.Errorf("content of the response is not an integer got %s", content)
	}
	return intValue, nil
}

func (parser *Parser) consumeString(buffer *bytes.Buffer) (string, error) {
	content, err := parser.consumeContent(buffer)
	if err != nil {
		return "", err
	}
	return content, err
}

func (parser *Parser) consumeType(buffer *bytes.Buffer) (DataType, error) {
	rawResultType, err := buffer.ReadByte()
	if err != nil {
		return UnknownType, err
	}
	resultType := parser.GetType(rawResultType)
	if resultType == UnknownType {
		return UnknownType, fmt.Errorf("unknown type detected got %c", rawResultType)
	}
	return resultType, nil
}

func (parser *Parser) consumeBulkString(buffer *bytes.Buffer) (*string, error) {
	size, err := parser.consumeInt(buffer)
	if err != nil {
		return nil, errors.New("cannot read size of the bulk element ")
	}
	if size == -1 {
		return nil, nil
	}
	buffer.Next(1)
	content := string(buffer.Next(size))
	return &content, nil
}

func (parser *Parser) readError(buffer *bytes.Buffer) error {
	errMsg, err := parser.consumeContent(buffer)
	if err != nil {
		return errors.New("cannot read error message from the response ")
	}
	return errors.New(errMsg)
}

func (parser *Parser) handleError(buffer *bytes.Buffer, expected DataType, got DataType) error {
	if got == ErrorType {
		return parser.readError(buffer)
	}
	return fmt.Errorf("expected %s type got type %s", typeToString(expected), typeToString(got))
}

func (parser *Parser) ReadInt(response RawResponse) (int, error) {
	buf := bytes.NewBufferString(response)
	resultType, err := parser.consumeType(buf)
	if err != nil {
		return -1, err
	}
	if resultType != IntegerType {
		return -1, parser.handleError(buf, IntegerType, resultType)
	}
	intValue, err := parser.consumeInt(buf)
	if err != nil {
		return -1, err
	}
	return intValue, err
}

func (parser *Parser) ReadString(response RawResponse) (string, error) {
	buf := bytes.NewBufferString(response)
	resultType, err := parser.consumeType(buf)
	if err != nil {
		return "", err
	}
	if resultType != StringType {
		return "", parser.handleError(buf, StringType, resultType)
	}
	str, err := parser.consumeString(buf)
	if err != nil {
		return "", err
	}
	return str, nil
}

func (parser *Parser) ReadBulkString(response RawResponse) (*string, error) {
	buf := bytes.NewBufferString(response)
	resultType, err := parser.consumeType(buf)
	if err != nil {
		return nil, err
	}
	if resultType != BulkStringType {
		return nil, parser.handleError(buf, BulkStringType, resultType)
	}
	bulkStr, err := parser.consumeBulkString(buf)
	if err != nil {
		return nil, err
	}
	return bulkStr, nil
}

func (parser *Parser) ReadArray(response RawResponse) (*[]BoxedValue, error) {
	buf := bytes.NewBufferString(response)
	var err error
	resultType, err := parser.consumeType(buf)
	if err != nil {
		return &emptyArrayResult, err
	}
	if resultType != ArrayType {
		return &emptyArrayResult, parser.handleError(buf, ArrayType, resultType)
	}
	size, err := parser.consumeInt(buf)
	if err != nil {
		return &emptyArrayResult, err
	}
	buf.Next(1)
	array := make([]BoxedValue, size)
	idx := 0
	eof := false
	for !eof {
		resultType, err := parser.consumeType(buf)
		if errors.Is(err, io.EOF) {
			eof = true
			break
		}
		var value interface{}
		switch resultType {
		case IntegerType:
			value, err = parser.consumeInt(buf)
		case StringType:
			value, err = parser.consumeString(buf)
		case BulkStringType:
			value, err = parser.consumeBulkString(buf)
		default:
			return &emptyArrayResult, errors.New("only integer, string and bulk string are expected ")
		}
		if err != nil {
			return &emptyArrayResult, err
		}
		array[idx] = NewBoxedValue(resultType, value)
		idx++
		buf.Next(1)
	}
	return &array, nil
}

func NewParser() *Parser {
	return &Parser{}
}
