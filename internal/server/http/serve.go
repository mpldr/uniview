package http

import (
	"mime"
	"net/http"
	"path"
	"strings"
	templates "webinterface"

	"git.sr.ht/~poldi1405/glog"
)

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
		glog.Debug("handling as gRPC request")
		s.grpc.ServeHTTP(w, r)
		return
	}

	glog.Tracef("path: %q", r.URL.Path)
	switch r.URL.Path {
	case "/":
		glog.Debug("serving index")
		s.templates[TemplateIndex].Execute(w, nil)
	default:
		data, err := templates.Templates.ReadFile(path.Join("dist", r.URL.Path))
		if err == nil {
			glog.Debug("static asset found")
			w.Header().Add("Content-Type", mime.TypeByExtension(path.Ext(r.URL.Path)))
			w.Write(data)
			return
		}

		glog.Debug("serving room interface")
		s.templates[TemplateRoom].Execute(w, nil)
	}
}
