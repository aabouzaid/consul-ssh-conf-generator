package consul2ssh

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"reflect"
	"regexp"
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

func TestBuildTemplate(t *testing.T) {
	sshConf := sshNodeConf{
		Host: "foo_host",
		Main: mapInterface{
			"Hostname": "foo-fqdn",
			"User":     "foo_user",
			"Port":     22,
		},
	}

	expected := `
	  Host foo_host
	    Hostname foo-fqdn
	    Port 22
	    User foo_user
	`

	// Mock http request (http.ResponseWriter).
	mockRequest := httptest.NewRecorder()
	sshConf.buildTemplate(mockRequest, sshConfTemplate)
	actual := mockRequest.Body.String()

	// Strip all new lines and empty space to make it easy to match.
	r, _ := regexp.Compile("(\n|\\s)")
	expectedTrimmed := r.ReplaceAllString(expected, "")
	actualTrimmed := r.ReplaceAllString(actual, "")

	assert.Equal(t, expectedTrimmed, actualTrimmed)
}
