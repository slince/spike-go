#!/usr/bin/env bash

echo "Build Spike Start..."

Version=`git describe --tags`
BuildTime=`date +%FT%T%z`
GoVersion=`go version`

LDFLAGS="-w -s\
 -X 'github.com/slince/spike/pkg/build.Version=${Version}'\
 -X 'github.com/slince/spike/pkg/build.BuildTime=${BuildTime}'\
 -X 'github.com/slince/spike/pkg/build.GoVersion=${GoVersion}'\
"

echo "LDFLAGS=${LDFLAGS}"

function build() {
  echo "build $1 $2"
  CGO_ENABLED=0 GOOS="$1" GOARCH="$2" go build -trimpath -ldflags "$LDFLAGS" -o "dist/${1}_${2}"/ ./cmd/spike
  CGO_ENABLED=0 GOOS="$1" GOARCH="$2" go build -trimpath -ldflags "$LDFLAGS" -o "dist/${1}_${2}"/ ./cmd/spiked
}

os=(linux darwin windows)
arch=(386 amd64 arm)

for i in "${os[@]}" ; do
    for j in "${arch[@]}" ; do
        build "$i" "$j"
    done
done

echo "Build end.."