package main

import (
	"log"
	"net/http"
	"regexp"
	"github.com/gorilla/mux"
)


func main() {
	r := mux.NewRouter()
	r.HandleFunc("/auth/{node}/user/{user}", authCheck)
	r.HandleFunc("/", defaultResp)
	log.Fatal(http.ListenAndServe(":9001", r))
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
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(nodeAuth))
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
// This iis where real code could be inserted. 
func verifyAccess(node string, user string) (json string) {
	if node[0:1] == user[0:1] {
		json = `{"node":"` + node + `","user":"` + user + `","auth":true}`
		return
	} else {
		json = `{"node":"` + node + `","user":"` + user + `","auth":false}`
		return
	}
}
