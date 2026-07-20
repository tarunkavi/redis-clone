package core

/*
Type is first four bits and encoding is next 4 bits since golang does not have capability to
assign 4 bits to type and encoding seperately
TYPE ENCODING
0000 0000 ->
*/

type RedisObject struct {
	TypeEncoding uint8
	Value        interface{}
	LastAccessAt uint32
}

// The below may be uint8 but they should be treated as 4 bits because type and encoding seperately are 4 bits
var OBJ_TYPE_STRING uint8 = 0 << 4

// Following redis conventions from here
// https://github.com/redis/redis/blob/0f39801756b21464b2fade19dbe213c67b293790/src/server.h#L1020
var OBJ_ENCODING_RAW uint8 = 0
var OBJ_ENCODING_INT uint8 = 1
var OBJ_ENCODING_EMBSTR uint8 = 8
