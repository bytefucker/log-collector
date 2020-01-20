package main

import (
	"flag"
	"github.com/yihongzhi/logCollect/common/logger"
	"github.com/yihongzhi/logCollect/manager/routers"
	"net/http"
	"time"
)

var (
	router = routers.Router
	log    = logger.Instance
)

var (
	port     string
	debug    bool
	help     bool
	logLevel string
)

func init() {
	flag.StringVar(&port, "port", "8080", "server port")
	flag.StringVar(&logLevel, "log-level", "info", "log level")
	flag.BoolVar(&debug, "debug", false, "enable debug model")
	flag.BoolVar(&help, "h", false, "help message")

}

func main() {
	flagParse()
	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Infof("Server Listen On %s", srv.Addr)
	log.Error(srv.ListenAndServe())
}

func flagParse() {
	flag.Parse()
	logger.SetLevel(logLevel)
	if help {
		flag.PrintDefaults()
	}
}
