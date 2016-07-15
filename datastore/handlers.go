package main

import (
	"net/http"
	"sync"
	"github.com/gorilla/mux"
	"io/ioutil"
	"fmt"
)

const (
	TextContentType = "text/plain"
	ContentTypeHeader = "Content-Type"
	LocationHeader = "Location"
)

// In-memory datastore for storing key/value entries.
var dataStore *DataStore = &DataStore{
	storage: make(map[string]*interface{}),
	mutex: &sync.RWMutex{},
}

func compositeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Errorf("----- jsme u")
	switch r.Method {
	case "GET":
		getHandler(w, r)
	case "PUT":
		putHandler(w, r)
	case "DELETE":
		deleteHandler(w, r)
	}
}

func putHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	if key == "" {
		logger.Errorf("Key %s is not provided.", key)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Errorf("A value reading failed: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

 	var pointer interface{} = string(bytes)
	dataStore.Put(key, &pointer)

	w.Header().Set(ContentTypeHeader, TextContentType)
	w.Header().Set(LocationHeader, getEntityLocation(r, key))
	w.WriteHeader(http.StatusOK)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	err := dataStore.Delete(key)
	if _, ok := err.(*KeyNotExistError); ok {
		logger.Errorf("Key %s does not exist.", key)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	value, err := dataStore.Get(key)
	if _, ok := err.(*KeyNotExistError); ok {
		logger.Errorf("Key %s does not exist.", key)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set(ContentTypeHeader, TextContentType)
	w.WriteHeader(http.StatusOK)
	body := (*value).(string);
	w.Write([]byte(body))
}

func getEntityLocation(r *http.Request, entityId string) string {
	return r.RequestURI + "/" + entityId
}
