package main

// SPDX-FileCopyrightText: © Moritz Poldrack & AUTHORS
// SPDX-License-Identifier: AGPL-3.0-or-later

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"git.sr.ht/~mpldr/uniview/internal/buildinfo"
	"git.sr.ht/~mpldr/uniview/internal/config"
	"git.sr.ht/~poldi1405/glog"
	cli "github.com/jawher/mow.cli"
)

func main() {
	defer glog.PanicHandler()
	glog.SetLevel(glog.INFO)

	if filepath.Base(os.Args[0]) == "univiewd" {
		glog.Debug("starting in server mode")

		err := config.Load(&config.Server, config.ServerPaths)
		if err != nil {
			glog.Warnf("no config loaded: %v", err)
		}

		err = startServer()
		if err != nil {
			glog.Errorf("failed to start server: %v", err)
			os.Exit(1)
		}
		glog.Debug("server has shut down")
		return
	}

	app := cli.App("uniview", "synchronise video playback")
	app.Spec = "([--insecure] SERVER ROOM FILE) | URL | --version"

	var server, room, file, rawurl string
	var insecure, ver bool
	app.StringArgPtr(&server, "SERVER", "", "the server to connect to")
	app.StringArgPtr(&room, "ROOM", "", "the room to join")
	app.StringArgPtr(&file, "FILE", "", "the file to open")
	app.BoolOptPtr(&insecure, "i insecure", false, "do not validate the server certificate")
	app.StringArgPtr(&rawurl, "URL", "", "the uniview:// address to connect to")
	app.BoolOptPtr(&ver, "v version", false, "show version and quit")

	app.Action = func() {
		if ver {
			fmt.Println(buildinfo.BugReportVersion())
			return
		}
		glog.Infof("starting uniview version %s", buildinfo.VersionString())

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
