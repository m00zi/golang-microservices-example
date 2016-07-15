package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"github.com/Sirupsen/logrus"
	"flag"
	"pbouda/golang-microservices-example/discovery"
	"github.com/petrbouda/golang-http-client"
	"log"
	"fmt"
)

func init() {
	logger = logrus.New()
	logger.Out = os.Stdout
}

var (
	logger *logrus.Logger
	// Shared http Client for all REST Calls
	httpClient *http.Client
)

const (
	serviceName = "encryption"

	DefaultHost = "127.0.0.1"
	DefaultPort = "4000"
)

func main() {
	debug := flag.Bool("debug", false, "Enable verbose level in HTTP Client.")
	etcdUrl := flag.String("etcd", discovery.DefaultEtcdUrl, "Etcd Server URL address")
	enableDiscovery := flag.Bool("discovery", discovery.DefaultEnableDiscovery, "Enable Service Discovery")
	host := flag.String("host", DefaultHost, "Encryption Server Host")
	port := flag.String("port", DefaultPort, "Encrypton Server Port")

	// Initialize Shared HTTP Client
	httpClient = http_client.NewHttpClient(*debug)

	if *enableDiscovery {
		discovery.RegisterService(serviceName, *host, *port, *etcdUrl)
		fmt.Printf("Encryption service registered in Service Discovery: %s:%s\n", *host, *port)
	}

	println("Encryption-Server started!")
	router := createRouter()
	log.Fatal(http.ListenAndServe(":" + *port, router))
}

func createRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/encrypt/{id}", compositeHandler).
		Methods("GET", "POST")
	return router
}
