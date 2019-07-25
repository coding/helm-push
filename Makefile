VERSION := $(shell sed -n -e 's/version:[ "]*\([^"]*\).*/\1/p' plugin.yaml)
DIST := $(CURDIR)/_dist
LDFLAGS := "-X main.version=${VERSION}"

all: build

build:
	CGO_ENABLED=0 go build -mod=vendor -ldflags $(LDFLAGS) -o helm-push ./main.go

dist:
	mkdir -p $(DIST)
	GOOS=linux GOARCH=amd64 go build -mod=vendor -ldflags $(LDFLAGS) -o helm-push ./main.go
	tar -zcf $(DIST)/helm-push_$(VERSION)_linux-amd64.tgz helm-push
	GOOS=darwin GOARCH=amd64 go build -mod=vendor -ldflags $(LDFLAGS) -o helm-push ./main.go
	tar -zcf $(DIST)/helm-push_$(VERSION)_darwin-amd64.tgz helm-push
	GOOS=windows GOARCH=amd64 go build -mod=vendor -ldflags $(LDFLAGS) -o helm-push.exe ./main.go
	tar -zcf $(DIST)/helm-push_$(VERSION)_windows.amd64.tgz helm-push.exe

bootstrap:
	go mod download
	go mod vendor