package core

import "redis-clone/config"

// Lets implement LFU with moriss counter later
func evictFirst() {
	//evict the
	for k := range store {
		delete(store, k)
		return
	}
}
func evictMethod() {
	switch config.EvictionStrategy {
	case "simple-first":
		evictFirst()
	}

}
