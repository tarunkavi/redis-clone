package core

import (
	"errors"
	"fmt"
	"strconv"
)

/*
TYPE ENCODING
0101 1010
to get only TYPE AND with 1111 0000 ==> 0101 0000
*/
func getType(t uint8) uint8 {
	return t & 0b11110000
}
func getEncoding(t uint8) uint8 {
	return t & 0b00001111
}

func assertEncoding(a uint8, b uint8) error {
	if getEncoding(a) != b {
		fmt.Println(getEncoding(a))
		fmt.Println(b)
		return errors.New("the operation is not permitted on this type")
	}
	return nil
}

func assertType(a uint8, b uint8) error {
	if getType(a) != b {
		fmt.Println(getType(a))
		fmt.Println(b)
		return errors.New("the operation is not permitted on this type")
	}
	return nil
}
func getValueEncoding(val string) uint8 {
	_, err := strconv.ParseInt(val, 10, 64)
	if err == nil {
		return OBJ_ENCODING_INT
	}
	if len(val) < 44 {
		return OBJ_ENCODING_EMBSTR
	}
	return OBJ_ENCODING_RAW
}
