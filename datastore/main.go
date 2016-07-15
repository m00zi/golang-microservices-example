package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"github.com/Sirupsen/logrus"
	"flag"
	"pbouda/golang-microservices-example/discovery"
	"log"
	"fmt"
)

var logger *logrus.Logger

func init() {
	logger = logrus.New()
	logger.Out = os.Stdout
}

const (
	serviceName = "datastore"

	DefaultHost = "127.0.0.1"
	DefaultPort = "5000"
)

func main() {
	etcdUrl := flag.String("etcd", discovery.DefaultEtcdUrl, "Etcd Server URL address")
	enableDiscovery := flag.Bool("discovery", discovery.DefaultEnableDiscovery, "Enable Service Discovery")
	host := flag.String("host", DefaultHost, "Datastore Host")
	port := flag.String("port", DefaultPort, "Datastore Port")

	if *enableDiscovery {
		discovery.RegisterService(serviceName, *host, *port, *etcdUrl)
		fmt.Printf("Datastore service registered in Service Discovery: %s:%s\n", *host, *port)
	}

	println("Datastore Application started!")
	router := createRouter()
	log.Fatal(http.ListenAndServe(":" + *port, router))
}

func createRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/kv/{key}", compositeHandler).
		Methods("GET", "PUT", "DELETE")
	return router
}
