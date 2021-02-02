#!/bin/bash

set -e

VERSION="(untracked)"
REPO=$(cd $(dirname $0); pwd)
COMMIT_SHA=$(git rev-parse --short HEAD)
ASSETS="false"
BINARY="false"
RELEASE=""

debugInfo () {
  echo "Repo:           $REPO"
  echo "Build assets:   $ASSETS"
  echo "Build binary:   $BINARY"
  echo "Release:        $RELEASE"
  echo "Version:        $VERSION"
}

buildAssets () {
  cd $REPO
  rm -rf frontend/dist
  rm -f http/rice-box.go

  cd $REPO/frontend

  if [ "$CI" = "true" ]; then
    npm ci
  else
    cnpm install
  fi

  npm run lint
  npm run build
}

buildBinary () {
  if ! [ -x "$(command -v rice)" ]; then
    go get github.com/GeertJohan/go.rice
    go get github.com/GeertJohan/go.rice/rice
  fi

  cd $REPO/http
  rm -rf rice-box.go
  rice embed-go

  cd $REPO
  go build -a -o filebrowser -ldflags "-s -w -X github.com/filebrowser/filebrowser/v2/version.CommitSHA=$COMMIT_SHA -X github.com/filebrowser/filebrowser/v2/version.Version=$VERSION"

  export GOOS=linux
  export GOARCH=arm
  export GOARM=7
  export CGO_ENABLED=0 
  go build -a -o filebrowser_armv7 -ldflags "-s -w -X github.com/filebrowser/filebrowser/v2/version.CommitSHA=$COMMIT_SHA -X github.com/filebrowser/filebrowser/v2/version.Version=$VERSION"
}

release () {
  cd $REPO

  echo "👀 Checking semver format"

  if [ $# -ne 1 ]; then
    echo "❌ This release script requires a single argument corresponding to the semver to be released. See semver.org"
    exit 1
  fi
  
  GREP="grep"
  if [ -x "$(command -v ggrep)" ]; then
    GREP="ggrep"
  fi
  semver=$(echo "$1" | $GREP -P "^v(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)")
  if [ $? -ne 0 ]; then
    echo "❌ Not valid semver format. See semver.org"
    exit 1
  fi
  echo "🧼  Tidying up go modules"
  go mod tidy
  echo "🐑 Creating a new commit for the new release"
  git commit --allow-empty -am "chore: version $semver"
  git tag "$1"
  git push
  git push --tags origin
  echo "📦 Done! $semver released."
}

buildDockerImage() {
    if [ "$VERSION" == "" ]; then
        echo "build docker image failed!Not exist version string".
        exit 1
    fi
    docker build -f Dockerfile.armv7.debian -t filebrowser:${VERSION}_armv7 .
    docker save filebrowser:${VERSION}_armv7 -o filebrowser_${VERSION}_armv7.tar
    docker build -t filebrowser:$VERSION .
    docker save filebrowser:${VERSION} -o filebrowser_${VERSION}.tar
    echo "build docker images succeed!"
}

usage() {
  echo "Usage: $0 [-a] [-c] [-b] [-i] [-r <string>] [-v <string>]" 1>&2;
  exit 1;
}

DEBUG="false"
DOCKER_IMAGE="false"

while getopts "bacr:dv:i" o; do
  case "${o}" in
    b)
      ASSETS="true"
      BINARY="true"
      ;;
    a)
      ASSETS="true"
      ;;
    c)
      BINARY="true"
      ;;
    r)
      RELEASE=${OPTARG}
      VERSION=${OPTARG}
      ;;
    d)
      DEBUG="true"
      ;;
    v)
      VERSION=${OPTARG}
      ;;
    i)
      DOCKER_IMAGE="true"
      ;;
    *)
      usage
      ;;
  esac
done
shift $((OPTIND-1))

if [ "$DEBUG" = "true" ]; then
  debugInfo
fi

if [ "$ASSETS" = "true" ]; then
  buildAssets
fi

if [ "$BINARY" = "true" ]; then
  buildBinary
fi

if [ "$RELEASE" != "" ]; then
  release $RELEASE
fi

if [ "$DOCKER_IMAGE" = "true" ]; then
  buildDockerImage
fi
