# Consul2SSH.
[![Travis Build](https://img.shields.io/travis/AAbouZaid/consul-ssh-conf-generator/master.svg?logo=travis)](https://travis-ci.org/AAbouZaid/consul-ssh-conf-generator)
[![Releases](https://img.shields.io/github/release/AAbouZaid/consul-ssh-conf-generator.svg?logo=github)](https://github.com/AAbouZaid/consul-ssh-conf-generator/releases/latest)
[![Docker Image](https://img.shields.io/microbadger/image-size/aabouzaid/consul2ssh.svg?logo=docker&label=docker%20img)](https://hub.docker.com/r/aabouzaid/consul2ssh/)
[![CLI](https://img.shields.io/badge/CLI-amd64%2Ci386%2Carm-blue.svg?longCache=true)](https://github.com/AAbouZaid/consul-ssh-conf-generator/releases/latest)

API with CLI to get hosts from Consul and format them in SSH config style.


## Why?
When you are working in a hybrid dynamic environment (e.g. on-premise and public cloud), it's hard to track nodes when you need to access them.

So if you already use Consul as a service discovery, you can use `consul2ssh` as a middleware to generate SSH config based on nodes registered in Consul.

## Quick overview

### How it works?
Consul2SSH has 2 parts. It works as a middleware API where it call Consul to get registered members and return them in SSH config format.

<p align="center">
<img src="https://gist.githubusercontent.com/AAbouZaid/aee2010d4b0d0ff89adc517664b8f130/raw/f3d6f94ef331f28f6d64856cc040fc80e3ae83e3/consul2ssh_dia.png" width="320">
</p>

Since Consul2SSH is just a REST API, any client could be used to interact with it like `curl` or `postman`. Also Consul2SSH works as a client to call the API.


### Try it locally
You can try `consul2ssh` locally with no need to running Consul cluster as following:
```
# Mock Consul server.
go run sample/mock.go

# Run Consul2SSH API.
go run main.go listen

#
# Get nodes from Consul2SSH API.

# Using C2S CLI.
consul2ssh get

# Or using cURL.
curl http://localhost:8001/nodes -d '@sample/config.json'

---
# Output:
Host dv.bastion01
  ForwardAgent yes
  Hostname bastion01.fqdn
  Port 2222
  User aabouzaid

Host dv.node01
  Hostname node01.node.dev.consul
  Port 2222
  ProxyCommand ssh dv.bastion01 -W %h:%p
  User aabouzaid

Host dv.node02
  Hostname node02.node.dev.consul
  Port 2222
  ProxyCommand ssh dv.bastion01 -W %h:%p
  User aabouzaid

Host dv.cassandra-local-proxy
  Hostname bastion01.fqdn
  LocalForward 9042 node02:9042
  Port 2222
  TCPKeepAlive yes
  User aabouzaid
```

## More details

### Config file
C2S client sends JSON config file to API where it has all info that's needed to generate SSH config style from Consul.
In `samples` dir, there is an example for client config.

```
{
  "api": {
    "consul": "http://localhost:8501",
    "consul2ssh": "http://localhost:8001"
  },
  "main": {
    "prefix": "dv",
    "jumphost": "bastion01.fqdn",
    "domain": "consul"
  },
  "global": {
    "User": "aabouzaid",
    "Port": 2222
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
}
```

There are 5 main sections in the config file:
  - API: It has the address of `Consul` as well `consul2ssh`.
  - Main: It has configuration related to formating, naming, and accessing the hosts.
  - Global: SSH config as in `ssh_conf` file, those config are applied on all hosts.
  - PerNode: SSH config as in `ssh_conf` file, but applied (obviously!) pre node based on it's name.
  - Custom: This simply to access some internal services using `LocalForward` via Bastion host. 

Please note: Config file is meant to be used in client slide not API side. 


### API Service
Simply pull `consul2ssh` and run it, and it will listen to `8001` port.
Since configuration comes from the client, no configuration is need. You can set `LISTEN_HOST` and `LISTEN_PORT` as env vars.

You can test it using docker image:
```
docker pull aabouzaid/consul2ssh
docker run -p 8001:8001 -dt aabouzaid/consul2ssh
```

### Client agent
Any tool that works with HTTP/s like `curl` could interact with `consul2ssh`, however `consul2ssh` has a CLI client to make life easier.

By default C2S CLI will read config form `~/.consul2ssh/config.json` where it also has the URL for `consul2ssh` API.
So using C2S CLI all what you need is:
```
consul2ssh get -h
  -config-file string
    	Config file that will be used. (default "~/.consul2ssh/config.json")
  -url string
    	API URL for consul2ssh. (default "http://localhost:8001")

```

Also you can do the same with cURL:
```
curl http://consul2ssh.host:8081/nodes -d '@sample/config.json'
```

## To-do
More details at the project [Kanban board](https://github.com/AAbouZaid/consul-ssh-conf-generator/projects/1).
