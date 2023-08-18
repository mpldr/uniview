package main

import (
	"fmt"
	"net"

	"git.sr.ht/~mpldr/uniview/internal/mansion"
	"git.sr.ht/~mpldr/uniview/internal/server"
	"git.sr.ht/~mpldr/uniview/protocol"
	"git.sr.ht/~poldi1405/glog"
	"google.golang.org/grpc"
)

func startServer() error {
	srv := grpc.NewServer()
	protocol.RegisterUniViewServer(srv, &server.Server{
		Rooms: mansion.New(),
	})

	glog.Debugf("starting listener on %s:%d", "0.0.0.0", 1558)
	lis, err := net.Listen("tcp", "0.0.0.0:1558")
	if err != nil {
		return fmt.Errorf("failed to start listener: %w", err)
	}

	return srv.Serve(lis)
}
