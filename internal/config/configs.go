package config

// SPDX-FileCopyrightText: © Moritz Poldrack & AUTHORS
// SPDX-License-Identifier: AGPL-3.0-or-later

import (
	"os"
	"path"
)

// Server contains the configuration for the server
var Server = struct {
	General struct {
		Bind string `toml:"bind-to"`
		Host string `toml:"host"`
	} `toml:"general"`
	Advanced struct {
		EnableInstrumentation bool `toml:"enable-instrumentation"`
		JSONLog               bool `toml:"json-log"`
	} `toml:"advanced"`
}{
	General: struct {
		Bind string `toml:"bind-to"`
		Host string `toml:"host"`
	}{
		Bind: "0.0.0.0:1558",
	},
}

var Client = struct {
	General struct {
		Player string `toml:"player"`
	}
	Media struct {
		Directories []string `toml:"directories"`
	} `toml:"media"`
	Advanced struct {
		JSONLog bool `toml:"json-log"`
	} `toml:"advanced"`
}{
	Media: struct {
		Directories []string `toml:"directories"`
	}{
		Directories: []string{},
	},
}

func init() {
	home, err := os.UserHomeDir()
	if err == nil {
		Client.Media.Directories = append(Client.Media.Directories, path.Join(home, "Videos"))
	}

	cfg, err := os.UserConfigDir()
	if err == nil {
		ClientPaths = append(ClientPaths, path.Join(cfg, "uniview.toml"))
	}
}
