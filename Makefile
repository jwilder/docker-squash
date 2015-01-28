.SILENT :
.PHONY : docker-squash clean fmt

SHELL := /bin/bash
TAG := `git describe --always --dirty --tags`
LDFLAGS := -X main.buildVersion $(TAG)

all: docker-squash

docker-squash:
	echo "Building docker-squash"
	go install -ldflags "$(LDFLAGS)"

dist-clean:
	rm -rf dist
	rm -f docker-squash-*.tar.gz

dist: dist-clean
	for os in linux darwin ; do \
	  mkdir -p dist/$$os/amd64 && \
	  GOOS=$$os GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/$$os/amd64/docker-squash ; \
	  done

release: dist
	glock sync github.com/jwilder/docker-squash
	for os in linux darwin ; do \
	  tar -cvzf docker-squash-$${os}-amd64-$(TAG).tar.gz -C dist/$$os/amd64 docker-squash ; \
	  done

