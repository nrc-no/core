package types

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestStringOrArrayMarshaling tests that we can successfully Marshal and Unmarshal
// a StringOrArray to/from a string or an array of string, depending on the underlying value.
func TestStringOrArrayMarshaling(t *testing.T) {

	tests := []struct {
		name      string
		json      string
		val       StringOrArray
		expectErr bool
	}{
		{
			name: "array",
			json: `["val1","val2"]`,
			val: StringOrArray{
				Kind:       ArrayValue,
				ArrayValue: []string{"val1", "val2"},
			},
		}, {
			name: "string",
			json: `"val1"`,
			val: StringOrArray{
				Kind:        StringValue,
				StringValue: "val1",
			},
		}, {
			name: "null string",
			json: `null`,
			val: StringOrArray{
				Kind: NullValue,
			},
		}, {
			name:      "bad value",
			json:      `123`,
			expectErr: true,
		}, {
			name:      "bad string",
			json:      `"ab`,
			expectErr: true,
		}, {
			name:      "bad slice",
			json:      `["a","b"`,
			expectErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			jsonBytes, err := json.Marshal(test.val)
			if test.expectErr {
				assert.Error(t, err)
				return
			}
			if !assert.NoError(t, err) {
				return
			}
			assert.Equal(t, test.json, string(jsonBytes))
			var val StringOrArray
			if err := json.Unmarshal([]byte(test.json), &val); !assert.NoError(t, err) {
				return
			}
			assert.Equal(t, test.val, val)
		})
	}

}
