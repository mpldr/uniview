GO?=$(shell which go)
BUILD_OPTS?=-trimpath -v

GOSRC!=find * -type f \( -name '*.go' -and -not -name '*_test.go' \)
GOSRC+=go.mod go.sum

all: uniview univiewd

uniview: $(GOSRC)
	$(GO) build $(BUILD_OPTS) -o $@

univiewd: uniview
	ln -f $< $@

graph/model/models_gen.go: graph/schema.graphqls .gqlgen.yml
	$(GO) run github.com/99designs/gqlgen --config .gqlgen.yml generate
