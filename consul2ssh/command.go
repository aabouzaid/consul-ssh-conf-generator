package consul2ssh

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

const (
	c2sNodesEndpoint = "/nodes"
)

type cmd struct {
	flags    *flag.FlagSet
	confFile string
	url      string
}

func (c *cmd) init() {
	c.flags = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	c.flags.StringVar(&c.confFile, "config-file",
		"~/.consul2ssh/config.json", "Config file that will be used.")
	c.flags.StringVar(&c.url, "url",
		"http://localhost:8001", "API URL for consul2ssh.")
}

func readCMDArgs(args []string) *cmd {
	c := &cmd{}
	c.init()
	c.flags.Parse(args)
	return c
}

func readConfFile(file string) []byte {
	user, _ := user.Current()
	userHome := user.HomeDir

	// Expand tilde to home directory.
	if file == "~" {
		file = userHome
	} else if strings.HasPrefix(file, "~/") {
		file = filepath.Join(userHome, file[2:])
	}

	fileContent, err := ioutil.ReadFile(file)
	checkErrCMD(err)
	return fileContent
}

// GetNodesCMD - get nodes from API in SSH conf format.
func GetNodesCMD(args []string) {

	//
	cmdConf := readCMDArgs(args)
	var confData s2cConf
	jsonPayload := readConfFile(cmdConf.confFile)
	confData.get(bytes.NewReader(jsonPayload))
	c2sNodesURL := setStrVal(cmdConf.url, confData.API.C2SURL) + c2sNodesEndpoint

	//
	req, err := http.NewRequest("GET", c2sNodesURL, bytes.NewBuffer(jsonPayload))
	checkErrCMD(err)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	checkErrCMD(err)
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Print(string(respBody))
}
