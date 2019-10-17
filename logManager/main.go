package main

import (
	"flag"
	"logManager/logger"
	"logManager/routers"
	"net/http"
	"time"
)

var (
	router = routers.Instance
	log    = logger.Instance
)

var (
	port  string
	debug bool
	help  bool
)

func init() {
	flag.StringVar(&port, "port", "8080", "server port")
	flag.BoolVar(&debug, "debug", false, "enable debug model")
	flag.BoolVar(&help, "help", false, "help message")
}

func main() {
	flag.Parse()
	if help {
		flag.PrintDefaults()
	}
	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Infof("Server Listen On %s", srv.Addr)
	log.Error(srv.ListenAndServe())
}
