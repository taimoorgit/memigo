package main

import (
	"fmt"
	"sort"
)

var data = make(map[string][]byte)

func get(key string) ([]byte, error) {
	value, ok := data[key]

	if !ok {
		return nil, fmt.Errorf("key not found: %s", key)
	}

	return value, nil
}

func set(key string, value []byte) {
	data[key] = value
}

func list() []string {
	keys := []string{}
	for key := range data {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
