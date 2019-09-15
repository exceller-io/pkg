GO        = go
GOFMT     = gofmt
GOLINT    = golint

all: test

format:
	test -z "$($(GOFMT) -l $(find . -type f -name '*.go' -not -path "./vendor/*"))" || { echo "Run \"gofmt -s -w\" on your Golang code"; exit 1; }

lint:
	$(GOLINT) $($(GO) list ./...)

vet:
	$(GO) vet $($(GO) list ./...)

test: format lint vet
	$(GO) test $($(GO) list ./...) -coverprofile=cover.out

cover: test
	$(GO) tool cover -html=cover.out -o cover.html

