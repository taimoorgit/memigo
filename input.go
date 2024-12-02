package main

import (
	"fmt"
	"strings"
)

func runExpression(input string) ([]byte, error) {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return []byte{}, fmt.Errorf("empty command")
	}

	switch parts[0] {
	case "get":
		if len(parts) != 2 {
			return []byte{}, fmt.Errorf("invalid syntax for get: expected 'get <key>'")
		}

		value, err := get(parts[1])
		if err != nil {
			return []byte{}, err
		}

		return value, nil
	case "set":
		if len(parts) != 3 {
			return []byte{}, fmt.Errorf("invalid syntax for set: expected 'set <key> <value>'")
		}

		value := []byte(parts[2])
		set(parts[1], value)

		return []byte{}, nil
	case "list":
		if len(parts) != 1 {
			return []byte{}, fmt.Errorf("invalid syntax for list: expected 'list'")
		}

		keys := list()

		return []byte(strings.Join(keys, ", ")), nil
	default:
		return []byte{}, fmt.Errorf("unknown command: '%s'", parts[0])
	}
}
