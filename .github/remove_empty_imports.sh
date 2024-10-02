#!/bin/bash

# remove all blank lines in go 'imports' statements,

if [ $# != 1 ] ; then
  echo "usage: $0 <filename>"
  exit 1
fi

if [[ "$(uname -o)" == "Darwin" ]]; then
  sed -i '' '
    /^import/,/)/ {
      /^$/ d
    }
  ' $1
else
  sed -i'' '
    /^import/,/)/ {
      /^$/ d
    }
  ' $1
fi

