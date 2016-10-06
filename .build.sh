#!/usr/bin/env bash

set -xe

make_binaries() {
    gox -arch="amd64" \
        -os="darwin linux windows" \
        -output="crash_{{.OS}}_{{.Arch}}" \
        github.com/troykinsella/crash/cmd
}

gen_docs() {
    mkdocs build --clean
    tar -zcf crash_docs.tar.gz site/*

    ls -altR # TEMP
}

# Main

# Is a tag build?
if [ "$TRAVIS_PULL_REQUEST" == "false" ] && [ -n "$TRAVIS_TAG" ]; then
    make_binaries
    gen_docs
    sha256sum crash_* > sha256sum.txt
fi
