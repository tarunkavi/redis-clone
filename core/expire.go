package core

import (
	"log"
	"time"
)

/*
this works on principle of random sampling which says that if you pick 20 keys randomly and the expiry percentage
among them is more than 25 percent then we can say the entire set also have more than 25 percent expired keys
here.

Here golang hash map when iterated gives keys randomly

redis does this differntly where they have seperate maps for expired keys and only that map is iterated
*/

func expireSample() float32 {
	var limit int = 20
	var expiredCount int = 0
	for key, val := range store {
		expiry, ok := expires[val]
		if ok {
			limit--
			if expiry <= time.Now().Unix() {
				// log.Println("Logging for satisfaction")
				delete(store, key)
				delete(expires, val)
				expiredCount++
			}
		}

		if limit == 0 {
			break
		}

	}
	return float32(expiredCount) / float32(20.0)
}

func DeleteExpiredKeys() {
	for {
		frac := expireSample()
		if frac < 0.25 {
			break
		}
	}
	log.Println("deleted Expired keys len after deletion", len(store))
}

func hasExpired(expiry int64) bool {
	return expiry <= time.Now().Unix()
}
