package main

import (
	"os"

	"git.sr.ht/~poldi1405/glog"
)

// Version is filled with the programs version at compile time
var Version = "devel"

func main() {
	defer glog.PanicHandler()

	glog.Infof("starting uniview version %s", Version)
	if os.Args[0] == "univiewd" {
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
