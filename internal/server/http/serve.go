package http

// SPDX-FileCopyrightText: © Moritz Poldrack & AUTHORS
// SPDX-License-Identifier: AGPL-3.0-or-later

import (
	"log/slog"
	"mime"
	"net/http"
	"path"
	"strings"
	templates "webinterface"

	"git.sr.ht/~mpldr/uniview/internal/config"
	"git.sr.ht/~mpldr/uniview/internal/conman"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := s.requestID.Add(1)
	ip := r.RemoteAddr
	if forwardedFor := r.Header.Get("X-Forwarded-For"); forwardedFor != "" {
		ip = forwardedFor
	}
	log := slog.With("connection_id", id, "from", ip)
	log.Debug("incoming request", "headers", r.Header, "method", r.Method, "protocol", r.Proto)

	if r.ProtoMajor == 2 && strings.HasPrefix(r.Header.Get("Content-Type"), "application/grpc") {
		ctx := conman.SetLogger(r.Context(), log)
		log.Debug("handling as gRPC request")
		s.grpc.ServeHTTP(w, r.WithContext(ctx))
		return
	}

	switch r.URL.Path {
	case "/":
		log.Debug("serving index")
		s.lobbyPagesServed.Inc()
		s.templates[TemplateIndex].Execute(w, nil)
	case "/healthcheck":
		w.WriteHeader(http.StatusNoContent)
	case "/metrics":
		if config.Server.Advanced.EnableInstrumentation {
			promhttp.Handler().ServeHTTP(w, r)
			return
		}
		fallthrough
	default:
		data, err := templates.Templates.ReadFile(path.Join("dist", r.URL.Path))
		if err == nil {
			log.Debug("static asset found", "file", r.URL.Path)
			w.Header().Add("Content-Type", mime.TypeByExtension(path.Ext(r.URL.Path)))
			w.Write(data)
			return
		}

		log.Debug("serving room interface")
		s.roomPagesServed.Inc()
		s.templates[TemplateRoom].Execute(w, nil)
	}
}
