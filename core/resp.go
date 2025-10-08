package core

import (
	"errors"
	"fmt"
	"log"
)

var RESPNIL = []byte("$-1\r\n")

func DecodeAll(data []byte) (interface{}, int, error) {
	if len(data) == 0 {
		return nil, 0, errors.New("no data")
	}

	log.Println(string(data))
	switch data[0] {
	case '+':
		return readSimpleString(data)
	case '-':
		return readError(data)
	case ':':
		return readInteger(data)
	case '$':
		return readBulkStrings(data)
	case '*':
		return readArrays(data)
		// 	// default:
		// 	// 	return nil, errors.New("unknown RESP type")
	}

	return nil, 0, nil

}

// +OK\r\n
func readSimpleString(data []byte) (string, int, error) {
	pos := 1

	for ; data[pos] != '\r'; pos++ {
	}

	return string(data[1:pos]), pos + 2, nil

}

// -ERR unknown command\r\n
func readError(data []byte) (string, int, error) {

	return readSimpleString(data)
}

// :1000\r\n
func readInteger(data []byte) (int64, int, error) {
	pos := 1
	var value int64 = 0
	var sign int64 = 1
	if data[pos] == '-' {
		sign = -1
		pos = pos + 1
	}
	for ; data[pos] != '\r'; pos++ {
		value = value*10 + int64(data[pos]-'0')
	}

	return (sign) * value, pos + 2, nil
}

func readBulkStrings(data []byte) (string, int, error) {
	size, pos, err := readInteger(data)
	if err != nil {
		return "", pos, err
	}
	// $3\r\nfoo\r\n
	value := string(data[pos : pos+int(size)])

	return value, pos + int(size) + 2, nil
}
func readArrays(data []byte) ([]interface{}, int, error) {
	arraySize, pos, err := readInteger(data)
	if err != nil {
		return nil, pos, err
	}
	if arraySize == -1 {
		return nil, pos, err
	}

	arrayOutput := make([]interface{}, 0)

	for i := 0; i < int(arraySize); i++ {
		output, newpos, err := DecodeAll(data[pos:])
		if err != nil {
			return nil, newpos, err
		}
		pos = pos + newpos
		arrayOutput = append(arrayOutput, output)
	}
	return arrayOutput, pos, nil
}
func Decode(data []byte) ([]interface{}, error) {
	if len(data) == 0 {
		return nil, errors.New("empty data stream")
	}
	log.Println(string(data))
	var values []interface{} = make([]interface{}, 0)
	var index int = 0
	for index < len(data) {
		value, delta, err := DecodeAll(data[index:])
		if err != nil {
			return values, err
		}
		if delta == 0 {
			return values, nil
		}
		index = index + delta
		values = append(values, value)

	}

	return values, nil

}
func DecodeArrayString(output interface{}) ([]string, error) {
	ts := output.([]interface{})
	tokens := make([]string, 0)
	for _, v := range ts {
		tokens = append(tokens, v.(string))
	}
	return tokens, nil
}

func EncodeError(data error) []byte {
	return []byte(fmt.Sprintf("-%s\r\n", data))
}

// After Eval
func Encode(data interface{}, isSimple bool) []byte {
	switch v := data.(type) {
	case string:
		if isSimple {
			return []byte(fmt.Sprintf("+%s\r\n", v))
		}
		return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(v), v))
	case int64:
		return []byte(fmt.Sprintf(":%d\r\n", v))

	}

	return []byte{}
}
