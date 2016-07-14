package main

import (
	"net/http/httptest"
	"net/http"
	"testing"
	"strings"
	"os"
	"github.com/gorilla/handlers"
	"pbouda/golang-microservices-example/discovery"
	"io/ioutil"
)

func TestNoInstances(t *testing.T) {
	// Service Discovery does not contain any instance of datastore service
	discovery.DiscoveryClient = &StubEtcdRegistryClient{
		instances: make([]string, 0),
	}

	router := handlers.LoggingHandler(os.Stdout, createRouter())

	postReq, err := http.NewRequest("POST", "/encrypt/my-post", strings.NewReader("my-value"))
	checkRequest(err, t)

	postResp := httptest.NewRecorder()
	router.ServeHTTP(postResp, postReq)
	checkStatus(postResp.Code, http.StatusInternalServerError, t)
}

func TestPutValue(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		t.Logf("Incoming request to Datastore server for path '%s'", request.URL.Path)

		if request.Method == "POST" {
			if request.URL.Path == "/kv/my-key" {
				t.Logf("Datastore server returned value for path '%s'", request.URL.Path)
				response.Header().Set("Location", request.URL.Path)
				response.Write([]byte("my-value"))
			}
			response.WriteHeader(http.StatusOK)
			return
		}
	}))

	// Return datastore URL when the Registry client is called
	discovery.DiscoveryClient = &StubEtcdRegistryClient{
		instances: []string{server.URL},
	}

	router := handlers.LoggingHandler(os.Stdout, createRouter())

	postReq, err := http.NewRequest("POST", "/encrypt/my-key", strings.NewReader("my-value"))
	checkRequest(err, t)

	postResp := httptest.NewRecorder()
	router.ServeHTTP(postResp, postReq)
	checkStatus(postResp.Code, http.StatusCreated, t)

	respBody := postResp.Body.String()
	t.Logf("Returned key: '%s'", postResp.Body.String())
	if respBody == "" {
		t.Fatal("There is no returned key from encryption server.")
	}
}

func TestGetValue(t *testing.T) {
	key := "my-key-2"

	var storedValue string
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		t.Logf("Incoming request to Datastore server for path '%s', method '%s'", request.URL.Path, request.Method)

		if request.Method == "PUT" {
			if request.URL.Path == "/kv/" + key {
				t.Logf("Datastore server returned value for POST path '%s'", request.URL.Path)
				bodyBytes, err := ioutil.ReadAll(request.Body)
				if err != nil {
					t.Fatal("Failed during reading POST method.")
				}
				storedValue = string(bodyBytes)
				t.Logf("Datastore server stored value '%s'", storedValue)
				response.Header().Set("Location", request.URL.Path)
			}
			response.WriteHeader(http.StatusOK)
		} else if request.Method == "GET" {
			if request.URL.Path == "/kv/" + key {
				t.Logf("Datastore server returned value for GET path '%s'", request.URL.Path)
				t.Logf("Datastore server sent value '%s'", storedValue)
				response.Write([]byte(storedValue))
			}
			response.WriteHeader(http.StatusOK)
		}
	}))

	// Return datastore URL when the Registry client is called
	discovery.DiscoveryClient = &StubEtcdRegistryClient{
		instances: []string{server.URL},
	}

	router := handlers.LoggingHandler(os.Stdout, createRouter())

	postReq, err := http.NewRequest("POST", "/encrypt/my-key-2", strings.NewReader("my-value"))
	checkRequest(err, t)

	postResp := httptest.NewRecorder()
	router.ServeHTTP(postResp, postReq)
	returnedKey := postResp.Body.String()

	getReq, err := http.NewRequest("GET", "/encrypt/my-key-2?key=" + returnedKey, nil)
	checkRequest(err, t)

	getResp := httptest.NewRecorder()
	router.ServeHTTP(getResp, getReq)
	returnedValue := getResp.Body.String()
	t.Logf("Returned value: '%s'", returnedValue)
	if returnedValue != "my-value" {
		t.Fatalf("There is not the expected value: %s", returnedValue)
	}
}

func checkStatus(actual int, expected int, t *testing.T) {
	if actual != expected {
		t.Fatal("Server error: Returned ", actual, " instead of ", expected)
	}
}

func checkRequest(err error, t *testing.T) {
	if err != nil {
		t.Fatalf("Creating request failed: %+v", err)
	}
}