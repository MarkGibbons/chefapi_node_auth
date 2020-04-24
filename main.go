package main

// TODO: TLS

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/MarkGibbons/chefapi_lib"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"regexp"
)

type restInfo struct {
	Cert string
	Key  string
	Port string
}

var flags restInfo

func main() {
	flagInit()

	r := mux.NewRouter()
	r.HandleFunc("/auth/{user}/node/{node}", authNodeCheck)
	r.HandleFunc("/auth/{user}/org/{org}", authOrgCheck)
	log.Fatal(http.ListenAndServe(":"+flags.Port, r))
}

func authNodeCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	node, ok := cleanInput(vars["node"])
	if !ok {
		inputerror(&w)
		return
	}
	user, ok := cleanInput(vars["user"])
	if !ok {
		inputerror(&w)
		return
	}
	nodeAuth := verifyAccess(node, user)
	nodeJson, _ := json.Marshal(nodeAuth)
	fmt.Printf("Node authorization %+v\n", nodeJson)
	w.WriteHeader(http.StatusOK)
	w.Write(nodeJson)
	return
}

func authOrgCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	org, ok := cleanInput(vars["org"])
	if !ok {
		inputerror(&w)
		return
	}
	user, ok := cleanInput(vars["user"])
	if !ok {
		inputerror(&w)
		return
	}
	orgAuth := chefapi_lib.Auth{}
	orgAuth.Org = org
	orgAuth.User = user
	orgAuth.Auth = false
	// Most users cannot join admin or pci* organizations
	if org != "admin" && org[0:2] != "pci" {
		orgAuth.Auth = true
	}
	// Users that start with a can do anything
	if user[0:0] == "a" {
		orgAuth.Auth = true
	}

	orgJson, _ := json.Marshal(orgAuth)
	fmt.Printf("Org authorization %+v\n", orgJson)
	w.WriteHeader(http.StatusOK)
	w.Write(orgJson)
	return
}

func cleanInput(in string) (out string, matched bool) {
	matched, _ = regexp.MatchString("^[[:word:]]+$", in)
	out = in
	return
}

func inputerror(w *http.ResponseWriter) {
	(*w).WriteHeader(http.StatusBadRequest)
	(*w).Write([]byte(`{"message":"Bad url input value"}`))
}

// verifyAccess checks to see if the user is authorized to change the node
// This is where real authorization code could be inserted.
func verifyAccess(node string, user string) (auth chefapi_lib.Auth) {
	auth.Auth = false
	auth.Node = node
	auth.User = user
	if node[0:1] == user[0:1] {
		auth.Auth = true
	} else {
		auth.Auth = false
	}
	if user[0:1] == "a" {
		auth.Auth = true
	}
	switch node[0:1] {
	case "r":
		auth.Group = "ravens"
	case "s":
		auth.Group = "slyfolks"
	default:
		auth.Group = "mugworts"
	}
	return
}

func flagInit() {
	restcert := flag.String("restcert", "", "Rest Certificate File")
	restkey := flag.String("restkey", "", "Rest Key File")
	restport := flag.String("restport", "9001", "Rest interface https port")
	flag.Parse()
	flags.Cert = *restcert
	flags.Key = *restkey
	flags.Port = *restport
	return
}
