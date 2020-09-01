package parser

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
)

const StringType = 1
const ErrorType = 2
const IntegerType = 3
const BulkStringType = 4
const ArrayType = 5
const UnknownType = 6

type RawResponse = string
type ResultType = int
type BoxedValue struct {
	resultType ResultType
	value	   interface{}
}

var emptyArrayResult = []BoxedValue{}

func NewBoxedValue(resultType ResultType, value interface{}) BoxedValue {
	return BoxedValue{
		resultType: resultType,
		value: value,
	}
}

func GetType(rawType uint8) ResultType {
	switch rawType {
	case '+': return StringType
	case '-': return ErrorType
	case ':': return IntegerType
	case '$': return BulkStringType
	case '*': return ArrayType
	default:
		return UnknownType
	}
}

func typeToString(resultType ResultType) string {
	switch resultType {
	case IntegerType: return "IntegerType"
	case ErrorType: return "ErrorType"
	case StringType: return "StringType"
	case BulkStringType: return "BulkStringType"
	case ArrayType: return "ArrayType"
	default:
		return "UnknownType"
	}
}

func consumeContent(buffer *bytes.Buffer) (string, error) {
	content, err := buffer.ReadString('\r')
	if err != nil {
		return "", err
	}
	return content[0: len(content) -1], err
}

func consumeInt(buffer *bytes.Buffer) (int, error) {
	content, err := consumeContent(buffer)
	if err != nil {
		return -1, err
	}
	intValue, err := strconv.Atoi(content)
	if err != nil {
		return -1, errors.New(fmt.Sprintf("Content of the response is not an integer got %s", content))
	}
	return intValue, nil
}

func consumeString(buffer *bytes.Buffer) (string, error) {
	content, err := consumeContent(buffer)
	if err != nil {
		return "", err
	}
	return content, err
}

func consumeType(buffer *bytes.Buffer) (ResultType, error) {
	rawResultType, err := buffer.ReadByte()
	if err != nil {
		return UnknownType, err
	}
	resultType := GetType(rawResultType)
	if resultType == UnknownType {
		return UnknownType, errors.New(fmt.Sprintf("Unknown type detected got %c", rawResultType))
	}
	return resultType, nil
}

func consumeBulkString(buffer *bytes.Buffer) (*string, error) {
	size, err := consumeInt(buffer)
	if err != nil {
		return nil, errors.New("Cannot read size of the bulk element ")
	}
	if size == -1 {
		return nil, nil
	}
	buffer.Next(1)
	content := string(buffer.Next(size))
	return &content, nil
}

func readError(buffer *bytes.Buffer) error {
	errMsg, err := consumeContent(buffer)
	if err != nil {
		return errors.New("Cannot read error message from the response ")
	}
	return errors.New(errMsg)
}

func handleError(buffer *bytes.Buffer, expected ResultType, got ResultType) error {
	if got == ErrorType {
		return readError(buffer)
	}
	return errors.New(fmt.Sprintf("Expected %s type got type %s", typeToString(expected), typeToString(got)))
}

func ReadInt(response RawResponse) (int, error) {
	buf := bytes.NewBufferString(response)
	resultType, err := consumeType(buf)
	if err != nil {
		return -1, err
	}
	if resultType != IntegerType {
		return -1, handleError(buf, IntegerType, resultType)
	}
	intValue, err := consumeInt(buf)
	if err != nil {
		return -1, err
	}
	return intValue, err
}

func ReadString(response RawResponse) (string, error) {
	buf := bytes.NewBufferString(response)
	resultType, err := consumeType(buf)
	if err != nil {
		return "", err
	}
	if resultType != StringType {
		return "", handleError(buf, StringType, resultType)
	}
	str, err := consumeString(buf)
	if err != nil {
		return "", err
	}
	return str, nil
}

func ReadBulkString(response RawResponse) (*string, error) {
	buf := bytes.NewBufferString(response)
	resultType, err := consumeType(buf)
	if err != nil {
		return nil, err
	}
	if resultType != BulkStringType {
		return nil, handleError(buf, BulkStringType, resultType)
	}
	bulkStr, err := consumeBulkString(buf)
	if err != nil {
		return nil, err
	}
	return bulkStr, nil
}

func ReadArray(response RawResponse) (*[]BoxedValue, error) {
	buf := bytes.NewBufferString(response)
	resultType, err := consumeType(buf)
	if err != nil {
		return &emptyArrayResult, err
	}
	if resultType != ArrayType {
		return &emptyArrayResult, handleError(buf, ArrayType, resultType)
	}
	size, err := consumeInt(buf)
	if err != nil {
		return &emptyArrayResult, err
	}
	buf.Next(1)
	array := make([]BoxedValue, size, size)
	idx := 0
	eof := false
	for !eof {
		resultType, err := consumeType(buf)
		if errors.Is(err, io.EOF) {
			eof = true
			break
		}
		var value interface{}
		switch resultType {
		case IntegerType:
			value, err = consumeInt(buf)
		case StringType:
			value, err = consumeString(buf)
		case BulkStringType:
			value, err = consumeBulkString(buf)
		default:
			return &emptyArrayResult, errors.New("Only integer, string and bulk string are expected ")
		}

		array[idx] = NewBoxedValue(resultType, value)
		idx++
		buf.Next(1)
	}
	return &array, nil
}