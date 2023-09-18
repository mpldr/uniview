package http

// SPDX-FileCopyrightText: © Moritz Poldrack & AUTHORS
// SPDX-License-Identifier: AGPL-3.0-or-later

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"sync/atomic"
	templates "webinterface"

	"git.sr.ht/~mpldr/uniview/internal/config"
	"git.sr.ht/~mpldr/uniview/internal/mansion"
	"github.com/prometheus/client_golang/prometheus"
)

type Server struct {
	grpc      http.Handler
	templates []*template.Template
	requestID *atomic.Uint64

	lobbyPagesServed prometheus.Counter
	roomPagesServed  prometheus.Counter
}

const (
	TemplateIndex = iota
	TemplateRoom
)

var templateList = []string{
	"dist/index.html",
	"dist/room.html",
}

func NewServer(g http.Handler, m *mansion.Mansion) (*Server, error) {
	srv := &Server{
		grpc:      g,
		requestID: &atomic.Uint64{},

		lobbyPagesServed: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "uniview",
			Subsystem: "httpserver",
			Name:      "lobby_served",
			Help:      "number of times the lobby (/) has been served",
		}),
		roomPagesServed: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "uniview",
			Subsystem: "httpserver",
			Name:      "roomui_served",
			Help:      "number of times the room ui has been served",
		}),
	}

	if config.Server.Advanced.EnableInstrumentation {
		prometheus.MustRegister(
			prometheus.NewCounterFunc(prometheus.CounterOpts{
				Namespace: "uniview",
				Subsystem: "httpserver",
				Name:      "requests",
				Help:      "the total number of received HTTP requests",
			}, func() float64 { return float64(srv.requestID.Load()) }),
			srv.lobbyPagesServed,
			srv.roomPagesServed,
			prometheus.NewCounterFunc(prometheus.CounterOpts{
				Namespace: "uniview",
				Subsystem: "grpcserver",
				Name:      "sessions",
				Help:      "the total number of initiated sessions",
			}, func() float64 { return float64(m.SessionCounter()) }),
			prometheus.NewCounterFunc(prometheus.CounterOpts{
				Namespace: "uniview",
				Subsystem: "grpcserver",
				Name:      "active_sessions",
				Help:      "the total number of connected sessions",
			}, func() float64 { return float64(m.ActiveSessions()) }),
			prometheus.NewCounterFunc(prometheus.CounterOpts{
				Namespace: "uniview",
				Subsystem: "grpcserver",
				Name:      "open_rooms",
				Help:      "the total number of allocated rooms",
			}, func() float64 { return float64(m.GetRoomCount()) }),
		)
	}

	var templateError bool
	for _, tmpl := range templateList {
		parsed, err := template.ParseFS(templates.Templates, tmpl)
		if err != nil {
			slog.Error("failed to parse template", "path", tmpl, "error", err)
			templateError = true
			continue
		}
		srv.templates = append(srv.templates, parsed)
	}
	if templateError {
		return nil, fmt.Errorf("failed to parse templates")
	}

	return srv, nil
}
