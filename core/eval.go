package core

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"time"
)

func evalPing(cmd *Cmd) ([]byte, error) {
	log.Println("evalPing", cmd)

	switch len(cmd.Args) {
	case 0:
		// log.Println("evalPing", cmd)
		return Encode("PONG", true), nil
	case 1:
		return Encode(cmd.Args[0], false), nil
	default:
		return nil, errors.New("ERR Wrong number of arguments for ping command")
	}

}
func evalSET(cmd *Cmd) ([]byte, error) {
	if len(cmd.Args) < 2 {
		return nil, errors.New("ERR Wrong number of arguments for SET command")
	}
	key, value := cmd.Args[0], cmd.Args[1]
	oType, enc := getTypeEncoding(value)
	var expiry int64 = -1
	//when type keeps changing then can include that logic insidde getTypeEncoding
	for i := 2; i < len(cmd.Args); i++ {
		switch cmd.Args[i] {
		case "EX", "ex":
			if i+1 >= len(cmd.Args) {
				return nil, errors.New("ERR syntax error")
			}
			seconds, err := strconv.ParseInt(cmd.Args[i+1], 10, 64)
			if err != nil {
				return nil, errors.New("ERR value is not an integer or out of range")
			}
			expiry = time.Now().Unix() + seconds

			i++
		default:
			return nil, errors.New("(error) ERR syntax error")
		}
	}
	val := NewValue(value, expiry, oType, enc)
	Put(key, val)
	return RESP_OK, nil
}
func evalGET(cmd *Cmd) ([]byte, error) {
	if len(cmd.Args) < 1 {
		return nil, errors.New("ERR Wrong number of arguments for GET command")
	}
	key := cmd.Args[0]
	val := Get(key)
	log.Println("Tarun", val)
	if val == nil {
		return RESP_NIL, nil
	}
	return Encode(val.Value, false), nil
}
func evalTTL(cmd *Cmd) ([]byte, error) {
	if len(cmd.Args) < 1 {
		return nil, errors.New("ERR Wrong number of arguments for TTL command")
	}

	key := cmd.Args[0]
	val := Get(key)
	fmt.Println("val:", val)
	if val == nil {
		return Encode(int64(-2), false), nil
	}
	expiry, ok := GetExpiry(val)
	if !ok {
		return Encode(int64(-1), false), nil
	}
	TTL := expiry - time.Now().Unix()
	if TTL < 0 {
		return Encode(int64(-2), false), nil
	}
	return Encode(TTL, false), nil

}

func evalDEL(cmd *Cmd) ([]byte, error) {
	sucessCount := 0
	for i := 0; i < len(cmd.Args); i++ {
		if Del(cmd.Args[i]) {
			sucessCount++
		}
	}
	return Encode(int64(sucessCount), false), nil
}

func evalEXPIRE(cmd *Cmd) ([]byte, error) {
	key := cmd.Args[0]
	if len(cmd.Args) < 2 {
		return nil, errors.New("ERR Wrong number of arguments for ping command")
	}
	expiry, err := strconv.ParseInt(cmd.Args[1], 10, 64)
	if err != nil {
		return nil, err
	}
	val := Get(key)
	if val == nil {
		return Encode(int64(0), false), nil
	}
	SetExpiry(val, expiry+time.Now().Unix())
	return Encode(int64(1), false), nil
}

func evalBGREWRITEAOF() ([]byte, error) {
	DumpAll()
	return RESP_OK, nil
}

func evalINCR(cmd *Cmd) ([]byte, error) {
	key := cmd.Args[0]

	val := Get(key)
	if val == nil {
		val = NewValue("0", -1, OBJ_TYPE_STRING, OBJ_ENCODING_INT)
		Put(key, val)
	}
	if err := assertType(val.TypeEncoding, OBJ_TYPE_STRING); err != nil {
		return nil, err
	}
	if err := assertEncoding(val.TypeEncoding, OBJ_ENCODING_INT); err != nil {
		return nil, err
	}
	valInt, err := strconv.ParseInt(val.Value.(string), 10, 64)
	if err != nil {

	}
	valInt = valInt + 1
	val.Value = strconv.Itoa(int(valInt))
	Put(key, val)
	return Encode(valInt, false), nil
}

func evalINFO() ([]byte, error) {
	var b []byte

	buff := bytes.NewBuffer(b)

	buff.WriteString("# Keyspace\r\n")

	for index, metrics := range KeySpaceStat {
		buff.WriteString(fmt.Sprintf("db%d:keys=%d,expires=%d,avg_ttl=0\r\n", index, metrics["keys"], metrics["expires"]))
	}

	return Encode(buff.String(), false), nil
}

// LATENCY LATEST / HISTOGRAM are probed by redis_exporter. Real Redis replies
// with an array; redigo's redis.Values() requires that, so we return an empty
// array (*0\r\n) since this clone tracks no latency data.
func evalLATENCY() ([]byte, error) {
	return Encode([]string{}, false), nil
}

func Eval(cmd *Cmd) ([]byte, error) {
	switch strings.ToUpper(cmd.Cmd) {
	case "PING":
		return evalPing(cmd)
	case "GET":
		return evalGET(cmd)
	case "SET":
		return evalSET(cmd)
	case "INCR":
		return evalINCR(cmd)
	case "TTL":
		return evalTTL(cmd)
	case "DEL":
		return evalDEL(cmd)
	case "EXPIRE":
		return evalEXPIRE(cmd)
	case "BGREWRITEAOF":
		return evalBGREWRITEAOF()
	case "INFO":
		return evalINFO()
	case "LATENCY":
		return evalLATENCY()
	}
	return Encode("#{cmd.Cmd} Unknown", true), nil
}

func EvalAndRespond(cmds Cmds, c io.ReadWriter) {
	var response []byte
	buf := bytes.NewBuffer(response)
	for _, cmd := range cmds {
		output, err := Eval(cmd)
		var errWrite error
		if err != nil {
			buf.Write(EncodeError(err))
		} else {
			buf.Write(output)
		}
		if errWrite != nil {
			log.Println("err:% write", err)
		}
	}
	c.Write(buf.Bytes())
}
