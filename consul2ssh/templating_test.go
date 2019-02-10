package consul2ssh

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFmtSSHElems(t *testing.T) {
	input := mapInterface{
		"key01": "val",
		"key02": true,
		"key03": 1,
	}
	output := fmtSSHElems(input)

	// Loop over all elements in input and make sure they converted to
	// concatinated string version.
	inputLen := make([]int, len(input))
	for keyID := range inputLen {
		key := fmt.Sprintf("key%02d", keyID+1)
		elemKeyValue := fmt.Sprintf("%v %v", key, input[key])
		assert.Equal(t, elemKeyValue, output[keyID])
	}
}
