package main

import (
	"bytes"
	"reflect"
	"strings"
	"testing"
)

func TestRunExpression(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		initialStore   map[string][]byte
		expectedOutput []byte
		expectedError  string
		finalStore     map[string][]byte
	}{
		{
			name:           "Set a key-value pair",
			input:          "set key1 value1",
			initialStore:   map[string][]byte{},
			expectedOutput: nil,
			expectedError:  "",
			finalStore:     map[string][]byte{"key1": []byte("value1")},
		},
		{
			name:           "Set a key-value pair (bad key name)",
			input:          "set key! value1",
			initialStore:   map[string][]byte{},
			expectedOutput: nil,
			expectedError:  "key contains invalid character",
			finalStore:     map[string][]byte{},
		},
		{
			name:           "Get an existing key",
			input:          "get key1",
			initialStore:   map[string][]byte{"key1": []byte("value1")},
			expectedOutput: []byte("value1"),
			expectedError:  "",
			finalStore:     map[string][]byte{"key1": []byte("value1")},
		},
		{
			name:           "Get an existing key (bad key name)",
			input:          "get key,",
			initialStore:   map[string][]byte{},
			expectedOutput: nil,
			expectedError:  "key contains invalid character",
			finalStore:     map[string][]byte{},
		},
		{
			name:           "Get a non-existent key",
			input:          "get key2",
			initialStore:   map[string][]byte{},
			expectedOutput: nil,
			expectedError:  "key not found: key2",
			finalStore:     map[string][]byte{},
		},
		{
			name:           "List keys",
			input:          "list",
			initialStore:   map[string][]byte{"key1": []byte("value1"), "key2": []byte("value2")},
			expectedOutput: []byte("key1, key2"),
			expectedError:  "",
			finalStore:     map[string][]byte{"key1": []byte("value1"), "key2": []byte("value2")},
		},
		{
			name:           "Invalid syntax for set",
			input:          "set key1",
			initialStore:   map[string][]byte{},
			expectedOutput: nil,
			expectedError:  "invalid syntax for set",
			finalStore:     map[string][]byte{},
		},
		{
			name:           "Unknown command",
			input:          "unknown",
			initialStore:   map[string][]byte{},
			expectedOutput: nil,
			expectedError:  "unknown command: 'unknown'",
			finalStore:     map[string][]byte{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Reset the store for each test
			data = make(map[string][]byte)
			for key, value := range tc.initialStore {
				data[key] = value
			}

			// Execute the function
			output, err := runExpression(tc.input)

			// Check for expected error
			if tc.expectedError != "" {
				if err == nil || !strings.Contains(err.Error(), tc.expectedError) {
					t.Fatalf("expected error: %q, got: %v", tc.expectedError, err)
				}
			} else if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Check for expected output
			if !bytes.Equal(output, tc.expectedOutput) {
				t.Errorf("expected output: %q, got: %q", tc.expectedOutput, output)
			}

			// Use reflect.DeepEqual for final store comparison
			if !reflect.DeepEqual(data, tc.finalStore) {
				t.Errorf("expected final store: %v, got: %v", tc.finalStore, data)
			}
		})
	}
}
