package config

// SPDX-FileCopyrightText: © Moritz Poldrack & AUTHORS
// SPDX-License-Identifier: AGPL-3.0-or-later

// Server contains the configuration for the server
var Server = struct {
	General struct {
		Bind string "toml:\"bind-to\""
		Host string "toml:\"host\""
	} "toml:\"general\""
}{
	General: struct {
		Bind string "toml:\"bind-to\""
		Host string "toml:\"host\""
	}{
		Bind: "127.1.2.4:1558",
	},
}
