package websocket

import (
	"encoding/json"
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"
)

func TestUnitCommandTypeUnmarshalJSON(t *testing.T) {
	testCases := map[string]struct {
		identity func(string) bool
		config   *quick.Config
	}{
		"Should unmarshal successfully": {
			identity: func(data string) bool {
				marshaled, err := json.Marshal(data)
				if err != nil {
					t.Errorf("Failed to marshal %s", data)
					return false
				}
				var commandType CommandType
				if err := json.Unmarshal(marshaled, &commandType); err != nil {
					t.Errorf("Unexpected error occurred %s", err)
					return false
				}
				return true
			},
			config: &quick.Config{
				Values: func(arguments []reflect.Value, rand *rand.Rand) {
					types := []string{string(GetBlockchainHeightCommand)}
					val := types[rand.Intn(len(types))]
					arguments[0] = reflect.ValueOf(val)
				},
			},
		},
		"Should not unmarshal": {
			identity: func(data string) bool {
				marshaled, err := json.Marshal(data)
				if err != nil {
					t.Errorf("Failed to marshal %s", data)
					return false
				}
				var commandType CommandType
				if err := json.Unmarshal(marshaled, &commandType); err == nil {
					t.Error("Unexpected error to occur")
					return false
				}
				return true
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			if err := quick.Check(tc.identity, tc.config); err != nil {
				t.Fatalf("Unexpected error occurred %s", err)
			}
		})
	}
}
