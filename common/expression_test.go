package common

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
		finalLog       string
	}{
		{
			name:           "Set a key-value pair (simple string data)",
			input:          "set key1 value1",
			initialStore:   map[string][]byte{},
			expectedOutput: nil,
			expectedError:  "",
			finalStore:     map[string][]byte{"key1": []byte("value1")},
			finalLog:       "set key1 value1\n",
		},
		{
			name:           "Set a key-value pair (base64 encoded data)",
			input:          "set key1 aGVsbG8gd29ybGQ=",
			initialStore:   map[string][]byte{},
			expectedOutput: nil,
			expectedError:  "",
			finalStore:     map[string][]byte{"key1": []byte("aGVsbG8gd29ybGQ=")},
			finalLog:       "set key1 aGVsbG8gd29ybGQ=\n",
		},
		{
			name:           "Set a key-value pair (bad key name)",
			input:          "set key! value1",
			initialStore:   map[string][]byte{},
			expectedOutput: nil,
			expectedError:  "key contains invalid character",
			finalStore:     map[string][]byte{},
			finalLog:       "",
		},
		{
			name:           "Get an existing key",
			input:          "get key1",
			initialStore:   map[string][]byte{"key1": []byte("value1")},
			expectedOutput: []byte("value1"),
			expectedError:  "",
			finalStore:     map[string][]byte{"key1": []byte("value1")},
			finalLog:       "",
		},
		{
			name:           "Get an existing key (bad key name)",
			input:          "get key,",
			initialStore:   map[string][]byte{},
			expectedOutput: nil,
			expectedError:  "key contains invalid character",
			finalStore:     map[string][]byte{},
			finalLog:       "",
		},
		{
			name:           "Get a non-existent key",
			input:          "get key2",
			initialStore:   map[string][]byte{},
			expectedOutput: nil,
			expectedError:  "key not found: key2",
			finalStore:     map[string][]byte{},
			finalLog:       "",
		},
		{
			name:           "List keys",
			input:          "list",
			initialStore:   map[string][]byte{"key1": []byte("value1"), "key2": []byte("value2")},
			expectedOutput: []byte("key1, key2"),
			expectedError:  "",
			finalStore:     map[string][]byte{"key1": []byte("value1"), "key2": []byte("value2")},
			finalLog:       "",
		},
		{
			name:           "Invalid syntax for set",
			input:          "set key1",
			initialStore:   map[string][]byte{},
			expectedOutput: nil,
			expectedError:  "invalid syntax for set",
			finalStore:     map[string][]byte{},
			finalLog:       "",
		},
		{
			name:           "Unknown command",
			input:          "unknown",
			initialStore:   map[string][]byte{},
			expectedOutput: nil,
			expectedError:  "unknown command: 'unknown'",
			finalStore:     map[string][]byte{},
			finalLog:       "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Reset the store for each test
			var buffer bytes.Buffer
			store := NewStore(&buffer)
			for key, value := range tc.initialStore {
				store.data[key] = value
			}

			output, err := runExpression(tc.input, store)

			if tc.expectedError != "" {
				if err == nil || !strings.Contains(err.Error(), tc.expectedError) {
					t.Fatalf("expected error: %q, got: %v", tc.expectedError, err)
				}
			} else if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !bytes.Equal(output, tc.expectedOutput) {
				t.Errorf("expected output: %q, got: %q", tc.expectedOutput, output)
			}

			if actual := buffer.String(); actual != tc.finalLog {
				t.Errorf("expected final log: %v, got: %v", tc.finalLog, actual)
			}

			if !reflect.DeepEqual(store.data, tc.finalStore) {
				t.Errorf("expected final store: %v, got: %v", tc.finalStore, store.data)
			}
		})
	}
}
