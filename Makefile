.SILENT :
.PHONY : docker-squash clean fmt

TAG:=`git describe --abbrev=0 --tags`
LDFLAGS:=-X main.buildVersion=$(TAG) -s -w

all: docker-squash

docker-squash:
	echo "Building docker-squash"
	go install -ldflags "$(LDFLAGS)"

dist-clean:
	rm -rf dist
	rm -f docker-squash-linux-*.tar.gz
	rm -f docker-squash-darwin-*.tar.gz

dist: dist-clean
	mkdir -p dist/linux/amd64 && GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/linux/amd64/docker-squash
	mkdir -p dist/darwin/amd64 && GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/darwin/amd64/docker-squash

release: dist
	tar -cvzf docker-squash-linux-amd64-$(TAG).tar.gz -C dist/linux/amd64 docker-squash
	tar -cvzf docker-squash-darwin-amd64-$(TAG).tar.gz -C dist/darwin/amd64 docker-squash
