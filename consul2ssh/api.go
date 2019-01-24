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

const (
	consulNodesEndpoint = "/v1/catalog/nodes"
)

var sshConfTemplate = `
Host {{ .Host }}
{{- range $item := fmtSSHElems .Main }}
  {{ $item }}
{{- end }}
`

type MapInterface map[string]interface{}

type Conf struct {
	API     APIConf                 `json:"api"`
	Main    MainConf                `json:"main"`
	Global  MapInterface            `json:"global"`
	PerNode map[string]MapInterface `json:"pernode"`
	Custom  map[string]MapInterface `json:"custom"`
}

type APIConf struct {
	ConsulURL string `json:"consul"`
	C2SURL    string `json:"consul2ssh"`
}

type MainConf struct {
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

type ConsulNodes []struct {
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

// GetNodes - Get Consul nodes.
func GetNodes(w http.ResponseWriter, r *http.Request) {

	// Read body.
	var conf Conf
	if err := conf.Get(r.Body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Misconfigured or no body provided!\n"))
		log.Printf("%s %s %s 400 error=%s\n", r.RemoteAddr, r.Method, r.URL, err)
		return
	}

	var (
		// Main config groups.
		ac = conf.API
		mc = conf.Main
		gc = conf.Global
		nc = conf.PerNode
		cc = conf.Custom

		// Common config.
		jumphostName        = strings.Split(mc.JumpHost, ".")[0]
		consulNodesEndpoint = ac.ConsulURL + consulNodesEndpoint
	)

	// Get nodes from Consul API, and format the output.
	var nodesList ConsulNodes
	nodesList.GetJSON(consulNodesEndpoint)
	// Put datacenter name as part of main config.
	mc.Datacenter = setStrVal(mc.Datacenter, nodesList[0].Datacenter)
	// Use datacenter name as prefix if there is no prefix.
	mc.Prefix = setStrVal(mc.Prefix, mc.Datacenter)

	// Nodes from
	for _, node := range nodesList {

		// Load global conf.
		nodeConf := MapInterface{}
		mergeMaps(gc, nodeConf)

		// Any special handling.
		if node.Name == jumphostName {
			nodeConf["Hostname"] = mc.JumpHost
		} else {
			nodeConf["Hostname"] = fmt.Sprintf("%s.node.%s.%s", node.Name, node.Datacenter, mc.Domain)
			nodeConf["ProxyCommand"] = fmt.Sprintf("ssh %s.%s -W %%h:%%p", mc.Prefix, jumphostName)
		}

		// Overwrite global if any.
		if nodeSpecialConf, hasSpecialConf := nc[node.Name]; hasSpecialConf {
			for key, value := range nodeSpecialConf {
				nodeConf[key] = value
			}
		}

		// Generate the template.
		sshConf := sshNodeConf{
			Host: fmt.Sprintf("%s.%s", mc.Prefix, node.Name),
			Main: nodeConf,
		}
		sshConf.buildTemplate(w, sshConfTemplate)

	}

	// Not real nodes but to access internal services using jumphost.
	for nodeHost, nodeCustomConf := range cc {
		// Load global conf.
		nodeConf := MapInterface{}
		mergeMaps(gc, nodeConf)
		mergeMaps(nodeCustomConf, nodeConf)
		nodeConf["Hostname"] = mc.JumpHost

		// Generate the template.
		sshConf := sshNodeConf{
			Host: fmt.Sprintf("%s.%s", mc.Prefix, nodeHost),
			Main: nodeConf,
		}
		sshConf.buildTemplate(w, sshConfTemplate)
	}

	// Log request.
	numberOfNodes := int(len(nodesList) + len(cc))
	log.Printf("%s %s %s 200 nodesNum=%d\n", r.RemoteAddr, r.Method, r.URL, numberOfNodes)
}
