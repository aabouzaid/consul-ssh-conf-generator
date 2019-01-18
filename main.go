package main

import (
	"fmt"
	"github.com/AAbouZaid/consul-ssh-conf-generator/consul2ssh"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func help() {
	// TODO: A better way to do this later.
	help_message := "Usage:" +
		"\n  API: Use \"listen\" to run as a deamon." +
		"\n  CLI: Use \"get\" to run as command line interface for the API.\n"
	fmt.Printf(help_message)
}

func api() {
	var (
		host, port = consul2ssh.GetEnvKey("LISTEN_HOST", "localhost"), consul2ssh.GetEnvKey("LISTEN_PORT", "8001")
		listenAddr = fmt.Sprintf("%s:%s", host, port)
	)
	router := mux.NewRouter()
	log.Println("Gorilla Mux started and listen:", listenAddr)
	router.HandleFunc("/nodes", consul2ssh.GetNodes).Methods("GET")
	log.Fatal(http.ListenAndServe(listenAddr, router))
}

func cmd() {
	consul2ssh.GetNodesCMD(os.Args[2:])
}

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		// Work as CLI.
		case "get":
			cmd()
		// Work as API.
		case "listen":
			api()
		// Help message.
		default:
			help()
		}
	} else {
		// Help message.
		help()
	}
}
