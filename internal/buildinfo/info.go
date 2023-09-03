package buildinfo

// SPDX-FileCopyrightText: © Moritz Poldrack & AUTHORS
// SPDX-License-Identifier: AGPL-3.0-or-later

import (
	"fmt"
	"runtime/debug"
	"strings"
)

var (
	Version                               = "unset"
	BuiltFor                              = "unknown Linux distribution"
	GitRevision                           = "unset"
	goos, goarch                          string
	commit, modified                      string
	compiler, compilerVersion, cgoEnabled string
)

func init() {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}

	BuiltFor = strings.ReplaceAll(BuiltFor, "_", " ")

	compilerVersion = info.GoVersion

	for _, data := range info.Settings {
		switch data.Key {
		case "-compiler":
			compiler = data.Value
		case "GOOS":
			goos = data.Value
		case "GOARCH":
			goarch = data.Value
		case "vcs.revision":
			commit = data.Value
		case "vcs.modified":
			if data.Value != "true" {
				continue
			}
			modified = "(modified tree)"
		case "CGO_ENABLED":
			if data.Value != "1" {
				continue
			}
			cgoEnabled = "(CGO enabled)"
		}
	}
}

func VersionString() string {
	if Version == "unset" {
		return "unsupported build"
	}

	return Version
}

func BugReportVersion() string {
	return fmt.Sprintf(`uniview
package version: %s for %s
built for: %s/%s
commit: %s %s
built using: %s-%s %s
`, Version, BuiltFor, goos, goarch, commit, modified, compiler, compilerVersion, cgoEnabled)
}
