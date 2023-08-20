GO?=$(shell which go)
BUILD_OPTS?=-trimpath -v

VERSION?=$(shell git describe --always --dirty || echo 0.1.0)
GO_LDFLAGS=-X main.Version=$(VERSION)

GOSRC!=find * -type f \( -name '*.go' -and -not -name '*_test.go' \)
GOSRC+=go.mod go.sum

all: uniview univiewd

uniview: $(GOSRC) protocol/uniview.pb.go protocol/uniview_grpc.pb.go
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
