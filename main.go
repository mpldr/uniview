package main

// SPDX-FileCopyrightText: © Moritz Poldrack & AUTHORS
// SPDX-License-Identifier: AGPL-3.0-or-later

import (
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"path/filepath"

	"git.sr.ht/~mpldr/uniview/internal/buildinfo"
	"git.sr.ht/~mpldr/uniview/internal/client"
	"git.sr.ht/~mpldr/uniview/internal/config"
	"git.sr.ht/~poldi1405/glog"
	cli "github.com/jawher/mow.cli"
)

var levels = map[string]slog.Level{
	"debug": slog.LevelDebug,
	"info":  slog.LevelInfo,
	"warn":  slog.LevelWarn,
	"error": slog.LevelError,
}

func main() {
	defer glog.PanicHandler()

	loglevel, ok := os.LookupEnv("LOG_LEVEL")
	if !ok {
		loglevel = "info"
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: loglevel == "debug",
		Level:     levels[loglevel],
	})))

	if filepath.Base(os.Args[0]) == "univiewd" {
		err := config.Load(&config.Server, config.ServerPaths)
		if err != nil {
			slog.Warn("no config loaded", "error", err)
		}

		if config.Server.Advanced.JSONLog {
			slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
				AddSource: loglevel == "debug",
				Level:     levels[loglevel],
			})))
		}

		slog.Debug("starting in server mode", "loaded_config", config.Server)

		err = startServer()
		if err != nil {
			slog.Error("failed to start server", "error", err)
			os.Exit(1)
		}
		slog.Debug("server has shut down")
		return
	}

	app := cli.App("uniview", "synchronise video playback")
	app.Spec = "[--log-level] (([--insecure] SERVER ROOM FILE) | URL | --version)"

	var server, room, file, rawurl string
	var insecure, ver bool
	app.StringOptPtr(&loglevel, "l log-level", loglevel, "sets the loglevel between debug, info, warn, and error")
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

		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: loglevel == "debug",
			Level:     levels[loglevel],
		})))

		err := config.Load(&config.Client, config.ClientPaths)
		if err != nil {
			slog.Warn("no config loaded", "error", err)
		}

		if config.Client.Advanced.JSONLog {
			slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
				AddSource: loglevel == "debug",
				Level:     levels[loglevel],
			})))
		}

		slog.Info("starting uniview", "version", buildinfo.VersionString())

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
			slog.Error("failed to parse URL", "url", u, "error", err)
			os.Exit(1)
		}

		slog.Debug("connecting", "url", u)
		err = client.StartClient(u)
		if err != nil {
			slog.Error("failed to start server", "error", err)
			os.Exit(1)
		}
	}

	app.Run(os.Args)
}
