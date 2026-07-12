package core

import (
	"fmt"
	"log"
	"os"
	"strings"
)

/*
TODO: AOF file needs to written for all the write commands (SET,UPDATE,DELETE)
TODO: AOF file is written by automic renmame(write to temp.aof then rename to db-file.aof)
TODO: Need to run this in a seperate process
*/
type Aof struct {
	file *os.File
}

var aofInstance *Aof

func InitAof(filepath string) {
	log.Println(filepath)
	fp, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("error initializing the file", err)
	}

	aofInstance = &Aof{
		file: fp,
	}
}

// // for all write operations SET,DEL,Expire

func writeToAof(key string, val *RedisObject) error {
	writeString := fmt.Sprintf("SET %s %s", key, val.Value)
	tokens := strings.Split(writeString, " ")
	aofInstance.file.Write(Encode(tokens, false))

	return nil
}
func DumpAll() {

	for k, v := range store {
		err := writeToAof(k, v)
		if err != nil {
			log.Println("error while writing %v", err)
		}
	}

}
