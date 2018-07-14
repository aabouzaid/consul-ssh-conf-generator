package main

import (
	"fmt"
	"github.com/AAbouZaid/consul-ssh-conf-generator/consul2ssh"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var (
	host, port = consul2ssh.GetEnvKey("LISTEN_HOST", ""), consul2ssh.GetEnvKey("LISTEN_PORT", "8001")
	listenAddr = fmt.Sprintf("%s:%s", host, port)
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/nodes", consul2ssh.GetNodes).Methods("GET")
	log.Fatal(http.ListenAndServe(listenAddr, router))
}
