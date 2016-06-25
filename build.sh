#!/usr/bin/env bash

VERSION="$1"
if [ -z "$VERSION" ]; then
  echo "version is required"
  exit 1
fi

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$DIR"

grunt --env=production && \
gb generate && \
gb build \
  -ldflags "-X 'main.buildTime=$(date)' -X 'main.buildUser=$(whoami)' -X 'main.buildHash=$(git rev-parse HEAD)' -X 'main.buildVersion=${VERSION}'" \
  all
