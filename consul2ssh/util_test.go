package consul2ssh

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"os/user"
	"testing"
)

func TestGetEnvKey(t *testing.T) {
	var testParameters = []struct {
		message  string
		keyVal   string
		defVal   string
		expected string
	}{
		{"Test non exists key", "", "def_key_val", "def_key_val"},
		{"Test exists key", "key_val", "def_key_val", "key_val"},
	}
	for _, tp := range testParameters {
		if tp.keyVal != "" {
			os.Setenv("FOO", tp.keyVal)
		}
		val := GetEnvKey("FOO", tp.defVal)
		if val != tp.expected {
			t.Errorf("%v: got %q, want %q.", tp.message, val, tp.expected)
		}
	}

}

func TestGetFilePath(t *testing.T) {
        var testParameters = []struct {
                message  string
                filePath string
                expected string
        }{
                {"Test file with tilde", "~/dummy_file.txt",
		func() string {user, _ := user.Current(); userHome := user.HomeDir; return userHome}() + "/dummy_file.txt"},
                {"Test file with tilde", "~",
		func() string {user, _ := user.Current(); userHome := user.HomeDir; return userHome}()},
                {"Test file absolute path", "/tmp/dummy_file.txt", "/tmp/dummy_file.txt"},
                {"Test file relative path", "dummy_file.txt", "dummy_file.txt"},
                {"Test file relative path with current dir", "./dummy_file.txt", "./dummy_file.txt"},
        }
        for _, tp := range testParameters {
		val := getFilePath(tp.filePath)
		if val != tp.expected {
			t.Errorf("%v: got %q, want %q.", tp.message, val, tp.expected)
		}
        }
}

func TestMergeMaps(t *testing.T) {
	sourceMap := mapInterface{
		"key01": "val",
		"key02": true,
		"key03": 1,
	}
	distMap := mapInterface{}
	mergeMaps(sourceMap, distMap)
	sourceMapLen := make([]int, len(sourceMap))
	for keyID := range sourceMapLen {
		key := fmt.Sprintf("key%02d", keyID)
		assert.Equal(t, sourceMap[key], distMap[key])
	}
}

func TestSetStrVal(t *testing.T) {
	var testParameters = []struct {
		message  string
		inVal    string
		defVal   string
		expected string
	}{
		{"Test empty val", "", "def_val", "def_val"},
		{"Test exists val", "val", "", "val"},
	}
	for _, tp := range testParameters {
		val := setStrVal(tp.inVal, tp.defVal)
		if val != tp.expected {
			t.Errorf("%v: got %q, want %q.", tp.message, val, tp.expected)
		}
	}

}
