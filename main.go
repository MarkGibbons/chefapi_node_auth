package main

// TODO: TLS

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"regexp"
	"github.com/gorilla/mux"
	"github.com/MarkGibbons/chefapi_lib"
)

type restInfo struct {
        Cert string
        Key string
        Port string
}

var flags restInfo

func main() {
	flagInit()

	r := mux.NewRouter()
	r.HandleFunc("/auth/{node}/user/{user}", authCheck)
	// Send in a json body with an array of nodes?
	r.HandleFunc("/", defaultResp)
	log.Fatal(http.ListenAndServe(":"+flags.Port, r))
}

func authCheck( w http.ResponseWriter, r *http.Request) {
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
	nodeAuth  :=  verifyAccess(node, user)
	nodeJson, _  := json.Marshal(nodeAuth)
	w.WriteHeader(http.StatusOK)
	w.Write(nodeJson)
	return
}

func cleanInput(in string) (out string, matched bool) {
        matched, _ = regexp.MatchString("^[[:word:]]+$", in)
	out = in
	return
}

func defaultResp(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(`{"message":"GET /auth/NODE/user/USER is the only valid method"}`))
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
