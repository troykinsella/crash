#!/usr/bin/env bash

OUT=tmp
mkdir -p $OUT
rm -rf $OUT/docs $OUT/crash_docs

process_file() {
  local parent="$1"
  local file="$2"
  #echo "process_file -> PARENT: $parent, FILE: $file"

  # filter
  [[ "$file" == *.xcf ]] && return 0

  mkdir -p $OUT/$parent

  if [[ "$file" == *md ]]; then
    gucci $parent/$file > $OUT/$parent/$file || exit 2
  else
    cp $parent/$file $OUT/$parent/$file
  fi

  return 0
}

process_dir() {
  local parent=$1
  local dir=$2
  #echo "process_dir -> PARENT: $parent, DIR: $dir"

  for d in $(ls $parent/$dir); do
    local path=$parent/$dir/$d
    test -f $path && process_file $parent/$dir $d
    test -d $path && process_dir $parent/$dir $d
  done

  return 0
}

process_dir . docs || exit 1
cp mkdocs.yml $OUT
(cd $OUT && mkdocs build --clean -d crash_docs)
(cd $OUT && tar -zcf ../crash_docs.tar.gz crash_docs/*)
