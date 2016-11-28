VERSION=$(shell cat version)
BUILD_TIME=$(shell date)
BUILD_USER=$(shell whoami)
BUILD_HASH=$(shell git rev-parse HEAD)
ARCH=amd64
OS=linux darwin

LDFLAGS := -ldflags "-X 'main.buildTime=$(BUILD_TIME)' -X 'main.buildUser=$(BUILD_USER)' -X 'main.buildHash=$(BUILD_HASH)' -X 'main.buildVersion=$(VERSION)'"

all: deps assets generate test build

clean:
	rm -rf bin/* pkg/*

deps:
	go get -u github.com/constabulary/gb/...
	go get -u github.com/jteeuwen/go-bindata/...

assets: deps
	npm install
	grunt --env=production

generate: assets
	gb generate goalbum

test: generate
	gb test all

build: test
	gb build ${LDFLAGS} all

distclean:
	@mkdir -p dist
	rm -rf dist/*

dist: generate test distclean
	for arch in ${ARCH}; do \
		for os in ${OS}; do \
			env GOOS=$${os} GOARCH=$${arch} gb build ${LDFLAGS} all; \
			mv bin/goalbum-$${os}-$${arch} dist/goalbum-${VERSION}-$${os}-$${arch}; \
		done; \
	done

sign: dist
	$(eval key := $(shell git config --get user.signingkey))
	for file in dist/*; do \
		gpg2 --armor --local-user ${key} --detach-sign $${file}; \
	done

package: sign
	for arch in ${ARCH}; do \
		for os in ${OS}; do \
			tar czf dist/goalbum-${VERSION}-$${os}-$${arch}.tar.gz -C dist goalbum-${VERSION}-$${os}-$${arch} goalbum-${VERSION}-$${os}-$${arch}.asc; \
		done; \
	done

.PHONY: all clean deps test assets generate build distclean dist sign package
