package consul2ssh

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
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
		"http://localhost:8001", "Config file that will be used.")
}

func readCMDArgs(args []string) *cmd {
	c := &cmd{}
	c.init()
	c.flags.Parse(args)
	return c
}

func readConfFile(file string) []byte {
	fileContent, err := ioutil.ReadFile(file)
	checkErrCMD(err)
	return fileContent
}
