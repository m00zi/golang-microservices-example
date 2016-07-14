package main

import (
	"fmt"
	"sync"
)

type KeyNotExistError struct {
	key string
}

func (e *KeyNotExistError) Error() string {
	return fmt.Sprintf("Key '%s' does not exist in the datastore.", e.key)
}

// A simple in-memory store for values based on key.
// All instances of this interface and their methods are thread-safe.
type DataStore struct {

	// In-memory storage of all Combined Feeds with their downloaded entries
	storage map[string]*interface{}

	// Synchronization of the access to in-memory datastore. Map is not thread-safe
	// and must be synchronized in the case of writers and readers existence.
	mutex   *sync.RWMutex
}

// Stores new data or replace an existing data stored under
// the key with the new one
func (ds *DataStore) Put(key string, data *interface{}) (*interface{}) {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()

	ds.storage[key] = data
	return data
}

// Retrieves data from the data store using a key parameter.
// Returns retrieved entity or an error if the key does not
// exist in datastore.
func (ds *DataStore) Get(key string) (*interface{}, error) {
	ds.mutex.RLock()
	defer ds.mutex.RUnlock()

	data := ds.storage[key]
	if data == nil {
		return nil, &KeyNotExistError{ key }
	}

	return data, nil
}

// Deletes the entity stored under the specified key.
// Returns error if the key does not exist in the datastore.
func (ds *DataStore) Delete(key string) error {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()

	data := ds.storage[key]
	if data == nil {
		return &KeyNotExistError{ key }
	}

	delete(ds.storage, key)
	return nil
}
