package http

import (
	"fmt"
	"html/template"
	"net/http"
	templates "webinterface"
)

type Server struct {
	grpc      http.Handler
	templates []*template.Template
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
		grpc: g,
	}

	for _, tmpl := range templateList {
		parsed, err := template.ParseFS(templates.Templates, tmpl)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %q: %w", tmpl, err)
		}
		srv.templates = append(srv.templates, parsed)
	}

	return srv, nil
}
