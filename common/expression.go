package common

import (
	"fmt"
	"strings"
)

func validateKeyName(key string) error {
	for _, char := range key {
		if !(char >= 'a' && char <= 'z' || char >= 'A' && char <= 'Z' || char >= '0' && char <= '9' || char == '_' || char == '-') {
			return fmt.Errorf("key contains invalid character: %q", char)
		}
	}
	return nil
}

func runExpression(input string, store *Store) ([]byte, error) {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return []byte{}, fmt.Errorf("empty command")
	}

	switch parts[0] {
	case "get":
		if len(parts) != 2 {
			return []byte{}, fmt.Errorf("invalid syntax for get: expected 'get <key>'")
		}

		err := validateKeyName(parts[1])
		if err != nil {
			return []byte{}, err
		}

		value, err := store.GetKey(parts[1])
		if err != nil {
			return []byte{}, err
		}

		return value, nil
	case "set":
		if len(parts) != 3 {
			return []byte{}, fmt.Errorf("invalid syntax for set: expected 'set <key> <value>'")
		}

		err := validateKeyName(parts[1])
		if err != nil {
			return []byte{}, err
		}

		// TODO: currently, set can fail if could not write to binary log
		value := []byte(parts[2])
		err = store.SetKey(parts[1], value)
		if err != nil {
			return []byte{}, err
		}

		return []byte{}, nil
	case "list":
		if len(parts) != 1 {
			return []byte{}, fmt.Errorf("invalid syntax for list: expected 'list'")
		}

		keys := store.ListKeys()

		return []byte(strings.Join(keys, ", ")), nil
	default:
		return []byte{}, fmt.Errorf("unknown command: '%s'", parts[0])
	}
}
