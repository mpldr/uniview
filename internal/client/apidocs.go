package client

// SPDX-FileCopyrightText: © Moritz Poldrack & AUTHORS
// SPDX-License-Identifier: AGPL-3.0-or-later

import (
	_ "embed"
	"net/http"
)

type apiDocWrapper struct {
	next http.Handler
}

func NewDocWrapper(next http.Handler) *apiDocWrapper {
	return &apiDocWrapper{
		next: next,
	}
}

func (a *apiDocWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		w.Header().Add("Content-Type", "text/html")
		w.Write(docPage)
		return
	}
	a.next.ServeHTTP(w, r)
}

//go:embed index.html
var docPage []byte
