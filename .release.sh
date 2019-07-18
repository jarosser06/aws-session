#!/bin/bash

release=$(git describe --always --tags)
release_dir="./releases"

systems="linux windows darwin"
for sys in $systems
do
  name=aws-session-${sys}.zip
  zip -j ${release_dir}/${name} ./bin/${sys}/aws-session* README.md
done
