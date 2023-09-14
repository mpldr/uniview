package client

// SPDX-FileCopyrightText: © Moritz Poldrack & AUTHORS
// SPDX-License-Identifier: AGPL-3.0-or-later

import (
	_ "embed"
	"fmt"
	"net/http"
	"net/url"
)

type apiDocWrapper struct {
	next       http.Handler
	remoteHost string
	secure     string
}

func NewDocWrapper(next http.Handler, connectTo *url.URL) *apiDocWrapper {
	secure := "s"
	if connectTo.Query().Has("insecure") {
		secure = ""
	}
	return &apiDocWrapper{
		next:       next,
		remoteHost: connectTo.Host,
		secure:     secure,
	}
}

func (a *apiDocWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		w.Header().Add("Content-Type", "text/html")
		w.Write(docPage)
		return
	}
	w.Header().Add("Access-Control-Allow-Origin", fmt.Sprintf("http%s://%s", a.secure, a.remoteHost))
	w.Header().Add("Access-Control-Allow-Methods", "*")
	a.next.ServeHTTP(w, r)
}

//go:embed index.html
var docPage []byte
