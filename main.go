package main

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"git.sr.ht/~poldi1405/glog"
	cli "github.com/jawher/mow.cli"
)

// Version is filled with the programs version at compile time
var Version = "0.1.0"

func main() {
	defer glog.PanicHandler()
	glog.SetLevel(glog.INFO)

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

	app := cli.App("uniview", "synchronise video playback")
	app.Spec = "[--insecure] SERVER ROOM FILE | URL"

	var server, room, file, rawurl string
	var insecure bool
	app.StringArgPtr(&server, "SERVER", "", "the server to connect to")
	app.StringArgPtr(&room, "ROOM", "", "the room to join")
	app.StringArgPtr(&file, "FILE", "", "the file to open")
	app.BoolOptPtr(&insecure, "i insecure", false, "do not validate the server certificate")
	app.StringArgPtr(&rawurl, "URL", "", "the uniview:// address to connect to")

	app.Action = func() {
		var err error
		var u *url.URL
		if len(rawurl) == 0 {
			rawurl = fmt.Sprintf("uniview://%s/%s", server, room)
			u, err = url.Parse(rawurl)
			q := u.Query()
			q.Add("file", file)
			if insecure {
				q.Add("insecure", "")
			}
			u.RawQuery = q.Encode()
		} else {
			u, err = url.Parse(rawurl)
		}
		if err != nil {
			glog.Errorf("failed to parse URL: %v", err)
			os.Exit(1)
		}

		glog.Debugf("connecting to %s", u)
		err = startClient(u)
		if err != nil {
			glog.Errorf("failed to start server: %v", err)
			os.Exit(1)
		}
	}

	app.Run(os.Args)
}
