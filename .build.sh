#!/usr/bin/env bash

set -xe

compile() {
  go build -o crash -v github.com/troykinsella/crash/cmd
  ln -s crash $GOPATH/bin/crash
}

cross_compile() {
  gox -arch="amd64" \
    -os="darwin linux windows" \
    -output="crash_{{.OS}}_{{.Arch}}" \
    github.com/troykinsella/crash/cmd
  ln -s crash_linux_amd64 $GOPATH/bin/crash
}

# Is a tag build?
if [ "$TRAVIS_PULL_REQUEST" == "false" ] && [ -n "$TRAVIS_TAG" ]; then
  cross_compile
else
  compile
fi
