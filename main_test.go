package main

import (
	// "fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/gorilla/mux"
)

// authCheck( w http.ResponseWriter, r *http.Request) {
func TestAuthCheck(t *testing.T) {
	// Check the status code and response body - authorized case
	req, err := http.NewRequest("GET", "/auth/mynode/user/myuser", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	// Invoke authCheck
	newAuthCheckServer().ServeHTTP(rr, req)
	if status := rr.Code; status!= http.StatusOK {
		t.Errorf("AuthCheck status code is not ok. Got: %v want: %v\n", status, http.StatusOK)
	}
	wantBody := `{"node":"mynode","user":"myuser","auth":true}`
	if rr.Body.String() != wantBody {
		t.Errorf("AuthCheck unexpected json returned. Expected: %v Got: %v\n", wantBody, rr.Body.String())
	}

	// Check the status code and response body - authorized case
	req, err = http.NewRequest("GET", "/auth/mynode/user/otheruser", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	// Invoke authCheck
	newAuthCheckServer().ServeHTTP(rr, req)
	// Check the status code and response body - unauthorized case
	if status := rr.Code; status!= http.StatusOK {
		t.Errorf("AuthCheck status code is not ok. Got: %v want: %v\n", status, http.StatusOK)
	}
	wantBody = `{"node":"mynode","user":"otheruser","auth":false}`
	if rr.Body.String() != wantBody {
		t.Errorf("AuthCheck unexpected json returned. Expected: %v Got: %v\n", wantBody, rr.Body.String())
	}

	// Check the status code and response body - invalid request
	req, err = http.NewRequest("GET", "/auth/mynode/user/other&user", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	// Invoke authCheck
	newAuthCheckServer().ServeHTTP(rr, req)
	// Check the status code and response body - unauthorized case
	if status := rr.Code; status!= http.StatusBadRequest {
		t.Errorf("AuthCheck status code is not expected. Got: %v want: %v\n", status, http.StatusBadRequest)
	}
	wantBody = `{"message":"Bad url input value"}`
	if rr.Body.String() != wantBody {
		t.Errorf("AuthCheck unexpected json returned. Expected: %v Got: %v\n", wantBody, rr.Body.String())
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

// defaultResp(w http.ResponseWriter, r *http.Request) {
func TestDefaultResp(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(defaultResp)
	handler.ServeHTTP(rr, req)

	// Check the status code and response body
	if status := rr.Code; status!= http.StatusBadRequest {
		t.Errorf("Status code is not expected. Got: %v want: %v\n", status, http.StatusBadRequest)
	}
	wantBody := `{"message":"GET /auth/NODE/user/USER is the only valid method"}`
	if rr.Body.String() != wantBody {
		t.Errorf("AuthCheck unexpected json returned. Expected: %v Got: %v\n", wantBody, rr.Body.String())
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
	// Invoke authCheck
	newAuthCheckServer().ServeHTTP(rr, req)
	// Check the status code and response body - unauthorized case
	if status := rr.Code; status!= http.StatusBadRequest {
		t.Errorf("AuthCheck status code is not expected. Got: %v want: %v\n", status, http.StatusBadRequest)
	}
	wantBody := `{"message":"Bad url input value"}`
	if rr.Body.String() != wantBody {
		t.Errorf("AuthCheck unexpected json returned. Expected: %v Got: %v\n", wantBody, rr.Body.String())
	}
}

// verifyAccess(node string, user string) (json string) {
func TestVerifyAccess(t *testing.T) {
	expected_hit := `{"node":"mynode","user":"myuser","auth":true}`
	hit := verifyAccess("mynode", "myuser")
        if hit != expected_hit {
		t.Errorf("Unexpected return: Wanted: %+v Got: %+v", expected_hit, hit)
	}
	expected_miss := `{"node":"mynode","user":"otheruser","auth":false}`
	miss := verifyAccess("mynode", "otheruser")
        if hit != expected_hit {
		t.Errorf("Unexpected return: Wanted: %+v Got: %+v", expected_miss, miss)
	}
}

func newAuthCheckServer() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/auth/{node}/user/{user}", authCheck)
	return r
}
