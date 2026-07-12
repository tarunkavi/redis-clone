package core

import (
	"log"
	"redis-clone/config"
	"time"
)

//key-->string
//val -->
// can we make this into a singleton instance

var store map[string]*RedisObject

func init() {
	store = make(map[string]*RedisObject)
}

func NewValue(value interface{}, expiry int64, oType uint8, oEnc uint8) *RedisObject {
	val := &RedisObject{
		TypeEncoding: oType | oEnc,
		Value:        value,
		expiry:       expiry,
	}
	return val
}
func Put(key string, val *RedisObject) {
	if len(store) > config.KeysLimit {
		evictMethod()
	}
	store[key] = val
}

func Get(key string) *RedisObject {
	log.Println("I am in GET if starting")
	val := store[key]
	if val == nil {
		return nil
	}
	if val.expiry != -1 && val.expiry <= time.Now().Unix() {
		log.Println("I am in GET if state")

		delete(store, key)

		return nil // delete if the expired key is accessed
	}
	log.Println("I am in GET")

	return val
}
func Del(key string) bool {
	if store[key] != nil {
		delete(store, key)
		return true
	}
	return false
}
