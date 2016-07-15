package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"io/ioutil"
	"strings"
	"encoding/base64"
	"io"
	"pbouda/golang-microservices-example/discovery"
)

const (
	textContentType = "text/plain"
	contentTypeHeader = "Content-Type"

	datastoreService = "datastore"
)

func compositeHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getHandler(w, r)
	case "POST":
		postHandler(w, r)
	}
}

// Process incoming request to store a value under the given id
func postHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// id fetching from path
	id := vars["id"]

	// Value fetching from body
	valueBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Errorf("A value reading failed: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	aesKey := generateAesKey()
	encrypted, err := aesEncrypt(aesKey, valueBytes)
	if err != nil {
		logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	base64encoded := base64.URLEncoding.EncodeToString(encrypted)
	bodyReader := strings.NewReader(base64encoded)
	req, err := discovery.NewRequest("PUT", datastoreService, "/kv/" + id, bodyReader)
	if err != nil {
		logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp, err := httpClient.Do(req)
	if err != nil || resp == nil {
		logger.Errorf("Invocation of the service '%s' failed: %+v", datastoreService, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	// Read all bytes to ensure the client reusability
	io.Copy(ioutil.Discard, resp.Body)

	w.Header().Set(contentTypeHeader, textContentType)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(base64.URLEncoding.EncodeToString(aesKey)))
}

// Process incoming request to retrieve value using the given id
func getHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// id fetching from path
	id := vars["id"]

	// encryption key fetching from query params
	queryParams := r.URL.Query()
	keys := queryParams["key"]
	if len(keys) != 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error: One encryption key is allowed"))
		return
	}

	key, err := base64.URLEncoding.DecodeString(keys[0])
	if err != nil {
		logger.Errorf("Error during Base64 decoding aes key: '%+v'", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	logger.Errorf("Keyyyyyyy: '%s'", id)
	req, err := discovery.NewRequest("GET", datastoreService, "/kv/" + id, nil)
	if err != nil {
		logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp, err := httpClient.Do(req)
	if err != nil || resp == nil {
		logger.Errorf("Invocation of the service '%s' failed: %+v", datastoreService, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read all bytes to ensure the client reusability
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Errorf("Reading entity form the service '%s' failed: %+v", datastoreService, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	decodedValue, err := base64.URLEncoding.DecodeString(string(body))
	if err != nil {
		logger.Errorf("Decoding the value failed: %+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	value, err := aesDecrypt(key, decodedValue)
	if err != nil {
		logger.Errorf("Decoding the value failed: %+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set(contentTypeHeader, textContentType)
	w.WriteHeader(http.StatusOK)
	w.Write(value)
}