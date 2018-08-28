# Consul2SSH.
[![Travis Build](https://img.shields.io/travis/AAbouZaid/consul-ssh-conf-generator/master.svg?logo=travis)](https://travis-ci.org/AAbouZaid/consul-ssh-conf-generator)
[![Releases](https://img.shields.io/github/release/AAbouZaid/consul-ssh-conf-generator.svg?logo=github)](https://github.com/AAbouZaid/consul-ssh-conf-generator/releases/latest)
[![Docker Image](https://img.shields.io/microbadger/image-size/aabouzaid/consul2ssh.svg?logo=docker&label=docker%20img)](https://hub.docker.com/r/aabouzaid/consul2ssh/)
[![CLI](https://img.shields.io/badge/CLI-amd64%2Ci386%2Carm-blue.svg?longCache=true)](https://github.com/AAbouZaid/consul-ssh-conf-generator/releases/latest)

API with CLI to get hosts from Consul and format them in SSH config style.

## Why?
When you are working in a hybrid dynamic environment (e.g. on-premise and public cloud), it's hard to track nodes when you need to access them.

And if you already use Consul as a service discovery, you can use `consul2ssh` as a middleware to generate SSH config based on nodes registered in Consul.

Work still in progress.
