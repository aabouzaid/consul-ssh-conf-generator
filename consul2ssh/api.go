package consul2ssh

import (
	"encoding/json"
	"io"
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
