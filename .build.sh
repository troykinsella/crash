#!/usr/bin/env bash

set -xe

make_binaries() {
    gox -arch="amd64" \
        -os="darwin linux windows" \
        -output="crash_{{.OS}}_{{.Arch}}" \
        github.com/troykinsella/crash/cmd
    ln -s crash_linux_amd64 $GOPATH/bin/crash
}

# Is a tag build?
if [ "$TRAVIS_PULL_REQUEST" == "false" ] && [ -n "$TRAVIS_TAG" ]; then
    make_binaries
    sha256sum crash_* > sha256sum.txt
fi
