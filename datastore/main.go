package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"github.com/Sirupsen/logrus"
	"flag"
	"pbouda/yoti/discovery"
)

var logger *logrus.Logger

func init() {
	logger = logrus.New()
	logger.Out = os.Stdout
}

const (
	serviceName = "datastore"
)

func main() {
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
	router.HandleFunc("/kv/{key}", compositeHandler).
		Methods("GET", "POST", "DELETE")
	return router
}
