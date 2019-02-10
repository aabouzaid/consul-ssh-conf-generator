package consul2ssh

import (
	"os"
	"testing"
	"io/ioutil"
)

func TestReadConfFile(t *testing.T) {
	var testParameters = []struct {
		message  string
		fileDir  string
		fileName string
		expected string
	}{
		{"Test file with tilde", "~/", "dummy_file.txt", "foo text"},
		{"Test file absolute path", "/tmp", "dummy_file.txt", "foo text"},
		{"Test file relative path", "", "dummy_file.txt", "foo text"},
		{"Test file relative path with dir", "./", "dummy_file.txt", "foo text"},
	}
	for _, tp := range testParameters {

		// TODO: Find a better way rather than creating tmpfile.
		content := []byte(tp.expected)
		dir := getFilePath(tp.fileDir)
		tmpfile, err := ioutil.TempFile(dir, tp.fileName)
		checkErrCMD(err)
		// Clean up later.
		defer os.Remove(tmpfile.Name())
		// Write content.
		if _, err := tmpfile.Write(content); err != nil {
			tmpfile.Close()
			checkErrCMD(err)
		}
		// Read content.
		fileContent, err := ioutil.ReadFile(tmpfile.Name())

		if string(fileContent) != tp.expected {
			t.Errorf("%v: got %q, want %q.", tp.message, fileContent, tp.expected)
		}
	}
}
