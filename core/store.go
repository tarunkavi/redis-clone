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
var expires map[*RedisObject]int64

func init() {
	store = make(map[string]*RedisObject)
	expires = make(map[*RedisObject]int64)
}

func NewValue(value interface{}, expiry int64, oType uint8, oEnc uint8) *RedisObject {
	val := &RedisObject{
		TypeEncoding: oType | oEnc,
		Value:        value,
		LastAccessAt: uint32(time.Now().Unix()), //this should have been 24 bits for efficiency since due to golang limitations we are going forward with 32 bits we can even combine this and TypeEncoding for efficiency. but for future i guess.
	}
	// expiry < 0 means "no expiry": the key stays absent from the expires map.
	if expiry >= 0 {
		SetExpiry(val, expiry)
	}
	return val
}

func SetExpiry(val *RedisObject, expiry int64) {
	expires[val] = expiry
}

// GetExpiry returns the absolute expiry (unix seconds) for val and whether one
// is set. A key with no expiry is simply absent from the expires map.
func GetExpiry(val *RedisObject) (int64, bool) {
	expiry, ok := expires[val]
	return expiry, ok
}

func Put(key string, val *RedisObject) {
	if len(store) > config.KeysLimit {
		evictMethod()
	}
	store[key] = val

	KeySpaceStat[0]["keys"]++
}

func Get(key string) *RedisObject {
	log.Println("I am in GET if starting")
	val := store[key]
	if val == nil {
		return nil
	}
	if expiry, ok := expires[val]; ok && hasExpired(expiry) {
		log.Println("I am in GET if state")

		Del(key)

		return nil // delete if the expired key is accessed
	}

	return val
}
func Del(key string) bool {
	val := store[key]
	if val != nil {
		delete(store, key)
		delete(expires, val)
		KeySpaceStat[0]["keys"]--
		return true
	}
	return false
}
