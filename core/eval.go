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
	val := NewValue(value, -1)
	for i := 2; i < len(cmd.Args); i++ {
		switch cmd.Args[i] {
		case "EX", "ex":
			if i+1 >= len(cmd.Args) {
				return nil, errors.New("ERR syntax error")
			}
			expiry, err := strconv.ParseInt(cmd.Args[i+1], 10, 64)
			if err != nil {
				return nil, errors.New("ERR value is not an integer or out of range")
			}
			val.expiry = time.Now().Unix() + expiry

			i++
		default:
			return nil, errors.New("(error) ERR syntax error")
		}
	}
	Put(key, val)
	return Encode("OK", true), nil
}
func evalGET(cmd *Cmd) ([]byte, error) {
	if len(cmd.Args) < 1 {
		return nil, errors.New("ERR Wrong number of arguments for GET command")
	}
	key := cmd.Args[0]
	val := Get(key)
	log.Println("Tarun", val)
	if val == nil {
		return RESPNIL, nil
	}
	return Encode(val.value, false), nil
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
	if val.expiry == -1 {
		return Encode(val.value, false), nil
	}
	TTL := val.expiry - time.Now().Unix()
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
	//expiry shoul be updated inside the PUT
	val.expiry = int64(expiry + time.Now().Unix())
	return Encode(int64(1), false), nil
}

func Eval(cmd *Cmd) ([]byte, error) {
	switch strings.ToUpper(cmd.Cmd) {
	case "PING":
		return evalPing(cmd)
	case "GET":
		return evalGET(cmd)
	case "SET":
		return evalSET(cmd)
	case "TTL":
		return evalTTL(cmd)
	case "DEL":
		return evalDEL(cmd)
	case "EXPIRE":
		return evalEXPIRE(cmd)
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
