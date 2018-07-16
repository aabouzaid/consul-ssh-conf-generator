package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var nodesJSON = `[
    {
        "Address": "10.0.0.10",
        "CreateIndex": 6517915,
        "Datacenter": "dev",
        "ID": "92e862d7-61c5-8892-e910-a3d4f8770b1e",
        "Meta": {
            "consul-network-segment": ""
        },
        "ModifyIndex": 6517937,
        "Node": "bastion01",
        "TaggedAddresses": {
            "lan": "10.0.0.10",
            "wan": "10.0.0.10"
        }
    },
    {
        "Address": "10.0.0.11",
        "CreateIndex": 6517844,
        "Datacenter": "dev",
        "ID": "74db79c0-1233-98d8-50a7-c4a1770dc6ae",
        "Meta": {
            "consul-network-segment": ""
        },
        "ModifyIndex": 6517844,
        "Node": "node01",
        "TaggedAddresses": {
            "lan": "10.0.0.11",
            "wan": "10.0.0.11"
        }
    },
    {
        "Address": "10.0.0.12",
        "CreateIndex": 6517844,
        "Datacenter": "dev",
        "ID": "51646f08-2542-4ea9-4e38-4050d858cdae",
        "Meta": {
            "consul-network-segment": ""
        },
        "ModifyIndex": 6517844,
        "Node": "node02",
        "TaggedAddresses": {
            "lan": "10.0.0.12",
            "wan": "10.0.0.12"
        }
    }
]
`

func nodes(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, nodesJSON)
	log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/v1/catalog/nodes", nodes).Methods("GET")
	log.Fatal(http.ListenAndServe(":8501", router))
}
