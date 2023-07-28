package main

import (
	"os"
	"path/filepath"

	"git.sr.ht/~poldi1405/glog"
)

// Version is filled with the programs version at compile time
var Version = "devel"

func main() {
	defer glog.PanicHandler()

	glog.SetLevel(glog.INFO)
	glog.Infof("starting uniview version %s", Version)
	if filepath.Base(os.Args[0]) == "univiewd" {
		// TODO: start server
		err := startServer()
		if err != nil {
			glog.Errorf("failed to start server: %v", err)
			os.Exit(1)
		}
		glog.Debug("server has shut down")
		return
	}
	// TODO: start client
	os.Exit(0)
}
