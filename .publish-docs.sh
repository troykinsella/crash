#!/usr/bin/env bash

set -e

if [ "$TRAVIS_PULL_REQUEST" == "false" ] && [ "$TRAVIS_BRANCH" == "master" ]; then

  echo "Publishing generated documentation..."

  SITE_DIR="$(pwd)/site"
  test -d $SITE_DIR || { echo "site dir doesn't exist: $SITE_DIR"; exit 1; }

  # Clone gh-pages
  cd
  git config --global user.email "travis@travis-ci.org"
  git config --global user.name "travis-ci"
  git clone --quiet --branch=gh-pages https://${GH_TOKEN}@github.com/troykinsella/crash gh-pages > /dev/null

  # Update gh-pages
  cd gh-pages
  rm -rf * > /dev/null
  git add --all
  cp -R $SITE_DIR/* .

  # Commit and push changes
  git add -f .
  git commit -m "Generated docs for master build $TRAVIS_BUILD_NUMBER"
  git push -fq origin gh-pages > /dev/null

  echo "Successfully published generated documentation"
fi
