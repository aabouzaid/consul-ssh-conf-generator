package main

import (
	"github.com/AAbouZaid/consul-ssh-conf-generator/consul2ssh"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/nodes", consul2ssh.GetNodes).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", router))
}
