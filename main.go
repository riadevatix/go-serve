package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"runtime/debug"

	"github.com/gorilla/handlers"
)

type config struct {
	Base     string
	Dir      string
	Port     uint
	Hostname string
}

var DefaultConfig config

func getVersionInfo() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return ""
	}
	mainVersion := info.Main.Version
	vcsVersion := ""
	for _, sett := range info.Settings {
		// vcs.revision is the git hash
		if sett.Key == "vcs.revision" {
			vcsVersion = sett.Value
		}
	}
	if mainVersion == "" || mainVersion == "(devel)" {
		return vcsVersion
	}

	return mainVersion
}

func main() {
	flag.StringVar(&DefaultConfig.Base, "base", "", "Base path")
	flag.StringVar(&DefaultConfig.Dir, "dir", ".", "Directory to serve")
	flag.UintVar(&DefaultConfig.Port, "port", 8080, "Port")
	flag.StringVar(&DefaultConfig.Hostname, "host", "localhost", "Hostname")
	version := flag.Bool("version", false, "Show version")

	flag.Parse()

	if *version {
		fmt.Println(getVersionInfo())
		os.Exit(0)
	}

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
