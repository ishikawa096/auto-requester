package scheduler

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestSelectRandomElement(t *testing.T) {
	tests := []struct {
		name     string
		body     []byte
		expected []byte
	}{
		{
			name:     "Valid JSON",
			body:     []byte(`{"name":"apple","children":[{"name":"banana"},{"name":"cherry"}]}`),
			expected: []byte(`{"name":"apple","children":[{"name":"banana"},{"name":"cherry"}]}`),
		},
		{
			name:     "Valid JSON array",
			body:     []byte(`[{"name":"apple"},{"name":"banana"},{"name":"cherry"}]`),
			expected: nil, // We can't predict the exact output, so we will check if it's one of the elements
		},
		{
			name:     "Empty JSON array",
			body:     []byte(`[]`),
			expected: []byte(`[]`),
		},
		{
			name:     "Invalid JSON",
			body:     []byte(`invalid json`),
			expected: []byte(`invalid json`),
		},
		{
			name:     "Nil",
			body:     nil,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := selectRandomElement(tt.body)

			if tt.expected == nil && tt.body != nil {
				// Check if the result is one of the elements in the original array
				var jsonArray []interface{}
				if err := json.Unmarshal(tt.body, &jsonArray); err == nil {
					found := false
					for _, elem := range jsonArray {
						elemBytes, _ := json.Marshal(elem)
						if reflect.DeepEqual(got, elemBytes) {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("selectRandomElement() = %s, want one of %s", got, tt.body)
					}
				}
				return
			}

			// IF the expected result is not nil, we can compare the results
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("selectRandomElement() = %s, want %s", got, tt.expected)
			}
		})
	}
}
