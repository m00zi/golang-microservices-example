package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"github.com/Sirupsen/logrus"
	"flag"
	"pbouda/golang-microservices-example/discovery"
)

func init() {
	logger = logrus.New()
	logger.Out = os.Stdout
}

var (
	logger *logrus.Logger
	// Enable verbose level in HTTP Client
	debug bool = false
)

const (
	serviceName = "encryption"
)

func main() {
	flag.BoolVar(&debug, "debug", false, "Enable verbose level in HTTP Client.")
	etcdUrl := flag.String("etcd", discovery.DefaultEtcdUrl, "Etcd Server URL address")
	enableDiscovery := flag.Bool("discovery", discovery.DefaultEnableDiscovery, "Enable Service Discovery")
	host := flag.String("host", discovery.DefaultHost, "Datastore Host")
	port := flag.String("port", discovery.DefaultPort, "Datastore Port")

	router := createRouter()
	http.ListenAndServe(":" + *port, router)

	if *enableDiscovery {
		discovery.RegisterService(serviceName, *host, *host, *etcdUrl)
	}
}

func createRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/encrypt/{id}", compositeHandler).
		Methods("GET", "POST")
	return router
}
