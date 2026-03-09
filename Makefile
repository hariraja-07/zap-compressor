.PHONY: build test clean install release

VERSION := 1.0.0

build:
	@mkdir -p dist
	./build.sh

test:
	go test -v ./...

clean:
	rm -rf dist

install:
	go install ./cmd/zap

release: clean build
	@echo "Release files in ./dist/"
	@ls -lh dist/

help:
	@echo "ZAP Build System"
	@echo ""
	@echo "Usage:"
	@echo "  make build     - Build for all platforms"
	@echo "  make test     - Run tests"
	@echo "  make clean    - Clean build artifacts"
	@echo "  make install  - Install to \$\$GOPATH/bin"
	@echo "  make release  - Create release build"