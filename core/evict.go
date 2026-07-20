package core

import "redis-clone/config"

// Lets implement LFU with moriss counter later
func evictFirst() {
	//evict the
	for k := range store {
		Del(k)
		return
	}
}

func evictRandom() {
	evictNumber := float64(config.KeysLimit) * config.EvictionRatio
	for k := range store {
		evictNumber--
		Del(k)
		if evictNumber <= 0 {
			return
		}
	}
}

func evictLRU() {

}

func evictMethod() {
	switch config.EvictionStrategy {
	case "simple-first":
		evictFirst()
	case "allkeys-random":
		evictRandom()
	case "lru":
		evictLRU()
	}

}
