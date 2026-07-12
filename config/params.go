package config

// Config holds all configuration parameters for the Redis clone
var Host string
var Port int
var KeysLimit int = 5
var EvictionStrategy string = "simple-first"

var AOFFile string = "db-master.aof"
