#!/usr/bin/env bash

mkdocs build --clean
tar -zcf crash_docs.tar.gz site/*
