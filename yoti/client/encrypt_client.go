package client

import (
	"pbouda/golang-microservices-example/discovery"
	"bytes"
	"net/http"
	"io/ioutil"
)

const encryptServiceName = "encryption"

type EncryptClient struct{
	HttpClient *http.Client
}

// Store accepts an id and a payload in bytes and requests that the
// encryption-server stores them in its data store
func (c *EncryptClient) Store(id, payload []byte) ([]byte, error) {
	req, err := discovery.NewRequest("POST", encryptServiceName, "/encrypt/" + id, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}

	return c.invoke(req)
}

// Retrieve accepts an id and an AES key, and requests that the
// encryption-server retrieves the original (decrypted) bytes stored
// with the provided id
func (c *EncryptClient) Retrieve(id, aesKey []byte) ([]byte, error) {
	req, err := discovery.NewRequest("GET", encryptServiceName, "/encrypt/" + id, nil)
	if err != nil {
		return nil, err
	}

	return c.invoke(req)
}

func (c *EncryptClient) invoke(request *http.Request) ([]byte, error) {
	resp, err := c.HttpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// return entire response body in bytes
	return ioutil.ReadAll(resp.Body)
}

