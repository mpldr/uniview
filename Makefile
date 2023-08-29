# SPDX-FileCopyrightText: © nobody
# SPDX-License-Identifier: CC0-1.0

GO?=$(shell which go)
BUILD_OPTS?=-trimpath -v

VERSION?=$(shell git describe --always --dirty || echo 0.2.1)
DISTRIBUTION?=$(shell source /etc/os-release || source /usr/lib/os-release || source /etc/initrd-release && echo -n $$PRETTY_NAME \($$BUILD_ID\) | sed 's/ /_/g')

GO_LDFLAGS:=
GO_LDFLAGS+=-X git.sr.ht/~mpldr/uniview/internal/buildinfo.Version=$(VERSION)
GO_LDFLAGS+=-X git.sr.ht/~mpldr/uniview/internal/buildinfo.BuiltFor=$(DISTRIBUTION)
GO_LDFLAGS+=$(EXTRA_GO_LDFLAGS)

GOSRC!=find * -type f \( -name '*.go' -and -not -name '*_test.go' \)
GOSRC+=go.mod go.sum

DESTDIR?=
PREFIX?=/usr
BINDIR?=$(PREFIX)/bin

all: uniview univiewd

uniview: $(GOSRC) protocol/uniview.pb.go protocol/uniview_grpc.pb.go Makefile
	$(GO) build $(BUILD_OPTS) -ldflags "$(GO_LDFLAGS)" -o $@

univiewd: uniview
	ln -f $< $@

protocol/uniview.pb.go: protocol/uniview.proto tools/protoc-gen-go
	protoc --plugin=tools/protoc-gen-go \
		--go_out=./ \
		--go_opt=paths=source_relative \
		$<

protocol/uniview_grpc.pb.go: protocol/uniview.proto tools/protoc-gen-go-grpc
	protoc --plugin=tools/protoc-gen-go-grpc \
		--go-grpc_out=./ \
		--go-grpc_opt=paths=source_relative \
		$<

tools/protoc-gen-go: go.mod
	$(GO) build -o $@ -v google.golang.org/protobuf/cmd/protoc-gen-go

tools/protoc-gen-go-grpc: go.mod
	$(GO) build -o $@ -v google.golang.org/grpc/cmd/protoc-gen-go-grpc

AUTHORS: .git/index
	git log '--pretty=format:%an%n%(trailers:key=co-authored-by,valueonly)' | sed -e 's/ <.*//' | sort -f | uniq | tail -n+2 > $@

.PHONY: install
install:
	install -Dm755 uniview $(DESTDIR)$(BINDIR)/uniview
	install -Dm755 uniview $(DESTDIR)$(BINDIR)/univiewd
	install -Dm755 contrib/uniview.desktop $(DESTDIR)$(PREFIX)/share/applications/uniview.desktop
	install -Dm755 contrib/icon.svg $(DESTDIR)$(PREFIX)/share/icons/hicolor/scalable/apps/uniview.svg
