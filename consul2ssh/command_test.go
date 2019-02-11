package consul2ssh

import (
	"fmt"
	"github.com/kami-zh/go-capturer"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
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
		fileContent := readConfFile(tmpfile.Name())

		if string(fileContent) != tp.expected {
			t.Errorf("%v: got %q, want %q.", tp.message, fileContent, tp.expected)
		}
	}
}

func TestGetNodesCMD(t *testing.T) {

	// CMD side.
	config := []byte(`{
	  "api": {
	  },
	  "main": {
	    "prefix": "dv",
	    "jumphost": "bastion01.fqdn",
	    "domain": "consul"
	  },
	  "global": {
	    "User": "foo",
	    "Port": 22
	  },
	  "pernode": {
	    "bastion01": {
	      "ForwardAgent": "yes"
	    }
	  },
	  "custom": {
	    "cassandra-local-proxy": {
	      "TCPKeepAlive": "yes",
	      "LocalForward": [
	        "9042 node02:9042"
	      ]
	    }
	  }
	}`)

	// TODO: Find a better way rather than creating tmpfile.
	tmpfile, err := ioutil.TempFile("", "config.json")
	checkErrCMD(err)
	// Clean up later.
	defer os.Remove(tmpfile.Name())
	// Write content.
	if _, err := tmpfile.Write(config); err != nil {
		tmpfile.Close()
		checkErrCMD(err)
	}

	// API side.
	expected := "dummy data"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, expected)
	}))
	defer server.Close()

	// Run GetNodesCMD function and read from stdout.
	args := []string{
		"-config-file", tmpfile.Name(),
		"-url", server.URL + c2sNodesEndpoint,
	}

	actual := capturer.CaptureStdout(func() {
		GetNodesCMD(args)
	})

	// Assert.
	assert.Equal(t, expected, strings.TrimSpace(actual))

}
