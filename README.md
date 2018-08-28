# Consul2SSH.
API with CLI to get hosts from Consul and format them in SSH config style.

## Why?
When you are working in a hybrid dynamic environment (e.g. on-premise and public cloud), it's hard to track nodes when you need to access them.

And if you already use Consul as a service discovery, you can use `consul2ssh` as a middleware to generate SSH config based on nodes registered in Consul.

Work still in progress.
