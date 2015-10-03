package main

import (
	"encoding/json"
	"log"
	"net/http"
)

var version struct {
	Major  int `json:"major"`
	Minor  int `json:"minor"`
	Hotfix int `json:"hotfix"`
}

func VersionHandler(w http.ResponseWriter, req *http.Request) {
	err := json.NewEncoder(w).Encode(version)
	if err != nil {
		log.Println("ERROR: ", err)
	}
}

func main() {
	version.Major = 0
	version.Minor = 0
	version.Hotfix = 0

	http.HandleFunc("/version", VersionHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalln("ListenAndServe: ", err)
	}
}
