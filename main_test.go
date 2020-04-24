package main

// Fix the tests, make sure auth true and false are returned correctly

import (
	// "fmt"
	"github.com/MarkGibbons/chefapi_lib"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"testing"
)

// authNodeCheck( w http.ResponseWriter, r *http.Request) {
func TestAuthNodeCheck(t *testing.T) {
	// Check the status code and response body - authorized case
	req, err := http.NewRequest("GET", "/auth/mynode/user/myuser", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	// Invoke authNodeCheck
	newAuthNodeCheckServer().ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("AuthNodeCheck status code is not ok. Got: %v want: %v\n", status, http.StatusOK)
	}
	wantBody := `{"auth":true,"group":"mugworts","node":"mynode","user":"myuser"}`
	if rr.Body.String() != wantBody {
		t.Errorf("AuthNodeCheck unexpected json returned. Expected: %v Got: %v\n", wantBody, rr.Body.String())
	}

	// Check the status code and response body - authorized case
	req, err = http.NewRequest("GET", "/auth/mynode/user/otheruser", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	// Invoke authNodeCheck
	newAuthNodeCheckServer().ServeHTTP(rr, req)
	// Check the status code and response body - unauthorized case
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("AuthNodeCheck status code is not ok. Got: %v want: %v\n", status, http.StatusOK)
	}
	wantBody = `{"auth":false,"group":"mugworts","node":"mynode","user":"otheruser"}`
	if rr.Body.String() != wantBody {
		t.Errorf("AuthNodeCheck unexpected json returned. Expected: %v Got: %v\n", wantBody, rr.Body.String())
	}

	// Check the status code and response body - invalid request
	req, err = http.NewRequest("GET", "/auth/mynode/user/other&user", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	// Invoke authNodeCheck
	newAuthNodeCheckServer().ServeHTTP(rr, req)
	// Check the status code and response body - unauthorized case
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("AuthNodeCheck status code is not expected. Got: %v want: %v\n", status, http.StatusBadRequest)
	}
	wantBody = `{"message":"Bad url input value"}`
	if rr.Body.String() != wantBody {
		t.Errorf("AuthNodeCheck unexpected json returned. Expected: %v Got: %v\n", wantBody, rr.Body.String())
	}
}

// cleanInput(in string) (out string, err error) {
func TestCleanInput(t *testing.T) {
	expected_in := "mynode"
	in, ok := cleanInput(expected_in)
	if !ok {
		t.Errorf("Error cleaning: %+v Err: %+v\n", expected_in, ok)
	}
	if in != expected_in {
		t.Errorf("In and out of cleanInput do not match in: %+v, out: %+v\n", expected_in, in)
	}
	expected_in = "\nbounceit"
	in, ok = cleanInput(expected_in)
	if ok {
		t.Errorf("CleanInput did not receive expected error cleaning: %+v Err: %+v\n", expected_in, ok)
	}
	if in != expected_in {
		t.Errorf("CleanInout in and out of cleanInput do not match in: %+v, out: %+v\n", expected_in, in)
	}
}

// inputerror(w *http.ResponseWriter) {
func TestInputerror(t *testing.T) {
	// Check the status code and response body - invalid request invoked inputerror
	req, err := http.NewRequest("GET", "/auth/mynode/user/other&user", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	// Invoke authNodeCheck
	newAuthNodeCheckServer().ServeHTTP(rr, req)
	// Check the status code and response body - unauthorized case
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("AuthNodeCheck status code is not expected. Got: %v want: %v\n", status, http.StatusBadRequest)
	}
	wantBody := `{"message":"Bad url input value"}`
	if rr.Body.String() != wantBody {
		t.Errorf("AuthNodeCheck unexpected json returned. Expected: %v Got: %v\n", wantBody, rr.Body.String())
	}
}

// verifyAccess(node string, user string) (json string) {
func TestVerifyAccess(t *testing.T) {
	expected_hit := chefapi_lib.Auth{Node: "mynode", Group: "mugworts", User: "myuser", Auth: true}
	hit := verifyAccess("mynode", "myuser")
	if hit != expected_hit {
		t.Errorf("Verify Access unexpected return: Wanted: %+v Got: %+v", expected_hit, hit)
	}
	expected_miss := chefapi_lib.Auth{Node: "mynode", Group: "mugworts", User: "otheruser", Auth: false}
	miss := verifyAccess("mynode", "otheruser")
	if hit != expected_hit {
		t.Errorf("Verify Access unexpected return: Wanted: %+v Got: %+v", expected_miss, miss)
	}
}

func newAuthNodeCheckServer() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/auth/{node}/user/{user}", authNodeCheck)
	return r
}
