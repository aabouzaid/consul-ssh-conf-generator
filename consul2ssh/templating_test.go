package consul2ssh

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestFmtSSHElems(t *testing.T) {
	input := mapInterface{
		"key01": "val",
		"key02": true,
		"key03": 1,
		"key04": []interface{}{
			"subkey4.1",
			"subkey4.2",
		},
	}
	output := fmtSSHElems(input)

	// Loop over all elements in input and make sure they converted to
	// concatinated string version.
	inputLen := make([]int, len(input))
	for keyID := range inputLen {
		key := fmt.Sprintf("key%02d", keyID+1)
		value := input[key]
		rt := reflect.TypeOf(value)

		switch rt.Kind() {
		// If the element is iterable we use the same main key but with element value.
		case reflect.Slice, reflect.Array:
			valueList := value.([]interface{})
			for subID, item := range valueList {
				elemKeyValue := fmt.Sprintf("%v %v", key, item)
				assert.Equal(t, elemKeyValue, output[keyID+subID])
			}
		// Non-iterable elements.
		default:
			elemKeyValue := fmt.Sprintf("%v %v", key, value)
			assert.Equal(t, elemKeyValue, output[keyID])
		}
	}
}
