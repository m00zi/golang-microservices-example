package discovery

import (
	"github.com/satori/go.uuid"
	"github.com/petrbouda/etcd_service_discovery"
	"log"
	"net/http"
	"io"
	"fmt"
	"errors"
)

var (
	// Client for register a service and getting other services
	DiscoveryClient discovery.RegistryClient
)

const (
	DefaultEtcdUrl = "http://127.0.0.1:4001"
	DefaultHost = "127.0.0.1"
	DefaultPort = "5000"
	DefaultEnableDiscovery = true
)

// Register the given service to discovery service.
func RegisterService(serviceName, host, port, etcd string) {
	// Create application configuration to send it to etcd
	registryConfig := discovery.EtcdRegistryConfig{
		EtcdEndpoints: []string{etcd},
		ServiceName:   serviceName,
		InstanceName:  serviceName + "-" + uuid.NewV4().String(),
		BaseURL:       "http://" + host + ":" + port,
	}

	// Creates etcd client, registers a current application and sends health-check requests
	client, err := discovery.New(registryConfig)
	if err != nil {
		log.Panicf("Cannot connect to etcd to register this service: %+v", err)
	}

	DiscoveryClient = client
	DiscoveryClient.Register()
}

// Register only client.
func RegisterClient(etcd string) discovery.RegistryClient {
	// Create application configuration to send it to etcd
	registryConfig := discovery.EtcdRegistryConfig{EtcdEndpoints: []string{etcd}}

	// Creates etcd client, registers a current application and sends health-check requests
	client, err := discovery.New(registryConfig)
	if err != nil {
		log.Panicf("Cannot connect to etcd to register this service: %+v", err)
	}

	DiscoveryClient = client
	return DiscoveryClient
}

// Creates a new request with service url which is automatically fetched from service discovery
// according to service name. Method concatenates fetched url with additional path, path parameter
// must start with '\' character.
func NewRequest(method, serviceName, path string, body io.Reader) (*http.Request, error) {
	// Check whether the application is already registered and has client assigned
	if DiscoveryClient == nil {
		return nil, errors.New("DiscoveryClient is not created, register an application or a discovery client.")
	}

	instances, err := DiscoveryClient.ServicesByName(serviceName)
	if err != nil || len(instances) <= 0 {
		return nil, fmt.Errorf("Cannot get instances of service: %s.", serviceName)
	}

	// Select first instance of registered instances belonging to the given service
	// TODO: create Round-robin mechanism if count is > 0
	serviceUrl := instances[0]

	// Create a new request with the URL of the service
	return http.NewRequest(method, serviceUrl + path, body)
}
