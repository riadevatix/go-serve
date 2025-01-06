package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/gorilla/handlers"
)

type config struct {
	Base     string
	Dir      string
	Port     uint
	Hostname string
}

var DefaultConfig config

func main() {
	flag.StringVar(&DefaultConfig.Base, "base", "", "Base path")
	flag.StringVar(&DefaultConfig.Dir, "dir", ".", "Directory to serve")
	flag.UintVar(&DefaultConfig.Port, "port", 8080, "Port")
	flag.StringVar(&DefaultConfig.Hostname, "host", "localhost", "Hostname")
	flag.Parse()

	base, _ := url.JoinPath("/", DefaultConfig.Base, "/")

	mux := http.NewServeMux()
	mux.Handle(base, http.StripPrefix(base, handlers.CombinedLoggingHandler(os.Stdout, http.FileServer(http.Dir(DefaultConfig.Dir)))))

	hostPort := fmt.Sprintf("%s:%d", DefaultConfig.Hostname, DefaultConfig.Port)

	server := http.Server{
		Addr:    hostPort,
		Handler: mux,
	}

	addr := url.URL{
		Scheme: "http",
		Host:   hostPort,
		Path:   base,
	}

	fmt.Println("Serving files on... " + addr.String())

	server.ListenAndServe()
}
