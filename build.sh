#!/bin/bash

set -e

VERSION="1.0.0"
DATE=$(date +%Y-%m-%d)

echo "Building ZAP v${VERSION} for all platforms..."

mkdir -p dist

PLATFORMS=(
	"darwin/arm64"
	"darwin/amd64"
	"linux/386"
	"linux/amd64"
	"linux/arm64"
	"windows/amd64"
)

for PLATFORM in "${PLATFORMS[@]}"; do
	GOOS=$(echo $PLATFORM | cut -d'/' -f1)
	GOARCH=$(echo $PLATFORM | cut -d'/' -f2)
	OUTPUT="zap-${VERSION}-${GOOS}-${GOARCH}"
	
	# Add .exe extension for Windows builds
	if [ "$GOOS" = "windows" ]; then
		OUTPUT="${OUTPUT}.exe"
	fi
	
	echo "Building ${OUTPUT}..."
	GOOS=$GOOS GOARCH=$GOARCH go build -ldflags="-s -w -X main.Version=${VERSION}" -o "dist/${OUTPUT}" ./cmd/zap/
done

echo "Creating checksums..."
cd dist
sha256sum zap-* > checksums.txt
cd ..

echo "✅ Build complete! Output in ./dist/"
ls -lh dist/