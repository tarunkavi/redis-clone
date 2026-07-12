package config

// Config holds all configuration parameters for the Redis clone
var Host string
var Port int
var KeysLimit int = 100
var EvictionStrategy string = "allkeys-random"

var AOFFile string = "db-master.aof"

// ratio of evicted keys
var EvictionRatio float64 = 0.40
