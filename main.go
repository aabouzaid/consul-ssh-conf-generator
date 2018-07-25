package main

import (
	"fmt"
	"github.com/AAbouZaid/consul-ssh-conf-generator/consul2ssh"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func api() {
	var (
		host, port = consul2ssh.GetEnvKey("LISTEN_HOST", ""), consul2ssh.GetEnvKey("LISTEN_PORT", "8001")
		listenAddr = fmt.Sprintf("%s:%s", host, port)
	)
	router := mux.NewRouter()
	router.HandleFunc("/nodes", consul2ssh.GetNodes).Methods("GET")
	log.Fatal(http.ListenAndServe(listenAddr, router))
}

func cmd() {
	consul2ssh.GetNodesCMD(os.Args[2:])
}

func main() {
	// Work as command line.
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "get":
			cmd()
		}
	// Work as API.
	} else {
		api()
	}
}
