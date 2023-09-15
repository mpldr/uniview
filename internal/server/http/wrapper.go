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
)

type Server struct {
	grpc      http.Handler
	templates []*template.Template
	requestID *atomic.Uint64
}

const (
	TemplateIndex = iota
	TemplateRoom
)

var templateList = []string{
	"dist/index.html",
	"dist/room.html",
}

func NewServer(g http.Handler) (*Server, error) {
	srv := &Server{
		grpc:      g,
		requestID: &atomic.Uint64{},
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
