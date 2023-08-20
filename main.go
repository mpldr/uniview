package main

import (
	"os"
	"path/filepath"

	"git.sr.ht/~poldi1405/glog"
)

// Version is filled with the programs version at compile time
var Version = "0.1.0"

func main() {
	defer glog.PanicHandler()

	glog.Infof("starting uniview version %s", Version)
	if filepath.Base(os.Args[0]) == "univiewd" {
		glog.Debug("starting in server mode")
		err := startServer()
		if err != nil {
			glog.Errorf("failed to start server: %v", err)
			os.Exit(1)
		}
		glog.Debug("server has shut down")
		return
	}

	err := startClient()
	if err != nil {
		glog.Errorf("failed to start server: %v", err)
		os.Exit(1)
	}
}
