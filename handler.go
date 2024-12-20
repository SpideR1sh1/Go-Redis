package main

import (
	"strings"
)

// Handlers map to associate commands with their functions
var Handlers = map[string]func([]Value) Value{
	"SET":  handleSet,
	"GET":  handleGet,
	"HSET": handleHSet,
	"HGET": handleHGet,
}

// In-memory storage for keys and values
var storage = make(map[string]string)
var hashStorage = make(map[string]map[string]string)

// handleSet processes the SET command
func handleSet(args []Value) Value {
	if len(args) < 2 || args[0].typ != "bulk" || args[1].typ != "bulk" {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'SET' command"}
	}

	key := args[0].bulk
	value := args[1].bulk
	storage[key] = value

	return Value{typ: "string", str: "OK"}
}

// handleGet processes the GET command
func handleGet(args []Value) Value {
	if len(args) < 1 || args[0].typ != "bulk" {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'GET' command"}
	}

	key := args[0].bulk
	value, ok := storage[key]
	if !ok {
		return Value{typ: "null"}
	}

	return Value{typ: "bulk", bulk: value}
}

// handleHSet processes the HSET command
func handleHSet(args []Value) Value {
	if len(args) < 3 || args[0].typ != "bulk" || args[1].typ != "bulk" || args[2].typ != "bulk" {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'HSET' command"}
	}

	hashKey := args[0].bulk
	field := args[1].bulk
	value := args[2].bulk

	if _, ok := hashStorage[hashKey]; !ok {
		hashStorage[hashKey] = make(map[string]string)
	}

	hashStorage[hashKey][field] = value
	return Value{typ: "string", str: "OK"}
}

// handleHGet processes the HGET command
func handleHGet(args []Value) Value {
	if len(args) < 2 || args[0].typ != "bulk" || args[1].typ != "bulk" {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'HGET' command"}
	}

	hashKey := args[0].bulk
	field := args[1].bulk

	if hash, ok := hashStorage[hashKey]; ok {
		if value, ok := hash[field]; ok {
			return Value{typ: "bulk", bulk: value}
		}
	}

	return Value{typ: "null"}
}

// Additional utility to display all stored keys (debugging purposes)
func handleKeys(args []Value) Value {
	if len(args) > 0 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'KEYS' command"}
	}

	keys := strings.Join(getAllKeys(), " ")
	return Value{typ: "bulk", bulk: keys}
}

// Helper function to get all keys from the in-memory storage
func getAllKeys() []string {
	keys := make([]string, 0, len(storage))
	for key := range storage {
		keys = append(keys, key)
	}
	return keys
}
