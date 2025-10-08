package core

import (
	"log"
	"redis-clone/config"
	"time"
)

//key-->string
//val -->
// can we make this into a singleton instance

var store map[string]*Value

/*
Without a separate dictionary, Redis would need to store expiration metadata alongside every key in the main keyspace,
increasing memory overhead for non-expiring keys and complicating lookup

For active deletion (background expiration task), Redis samples keys from the expires dictionary directly,
ensuring it only checks keys that have expiration times.This avoids wasting cycles on non-expiring key
*/
type Value struct {
	value  interface{}
	expiry int64
}

func init() {
	store = make(map[string]*Value)
}

func NewValue(value interface{}, expiry int64) *Value {
	val := &Value{
		value:  value,
		expiry: expiry,
	}
	return val
}
func Put(key string, val *Value) {
	if len(store) > config.KeysLimit {
		evictMethod()
	}
	store[key] = val
}

func Get(key string) *Value {
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
