//go:build tools

package main

// SPDX-FileCopyrightText: © nobody
// SPDX-License-Identifier: CC0-1.0

import (
	_ "github.com/josephspurrier/goversioninfo/cmd/goversioninfo"
	_ "github.com/ogen-go/ogen/cmd/ogen"
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
)
