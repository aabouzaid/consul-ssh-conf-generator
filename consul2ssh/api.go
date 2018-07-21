package consul2ssh

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type MapInterface map[string]interface{}

type Conf struct {
	Main    MainConf                `json:"main"`
	Global  MapInterface            `json:"global"`
	PerNode map[string]MapInterface `json:"pernode"`
	Custom  map[string]MapInterface `json:"custom"`
}

type MainConf struct {
	BaseURL    string `json:"baseurl"`
	Format     string `json:"format"`
	Prefix     string `json:"prefix"`
	JumpHost   string `json:"jumphost"`
	Domain     string `json:"domain"`
	Datacenter string `json:"datacenter"`
}

func (c *Conf) Get(reqBody io.Reader) error {
	err := json.NewDecoder(reqBody).Decode(&c)
	if err != nil {
		return err
	}
	return nil
}

type ConsulNodes []ConsulNode

type ConsulNode struct {
	Name       string `json:"node"`
	Datacenter string `json:"datacenter"`
}

func (c *ConsulNodes) GetJSON(apiURL string) error {

	// TODO: Better HTTP/s client.
	_, err := url.Parse(apiURL)
	if err != nil {
		log.Printf("E %s", err)
		return err
	}

	resp, err := http.Get(apiURL)
	if err != nil {
		log.Printf("E %s", err)
		return err
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("E! %s %s", apiURL, resp.Status)
		return err
	}

	if err := json.NewDecoder(resp.Body).Decode(&c); err != nil {
		log.Printf("E! %s", err)
		return err
	}

	return nil
}
