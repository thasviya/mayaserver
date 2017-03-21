#!/usr/bin/env bash
set -e

if [ -z "${CTLNAME}" ]; 
then
    CTLNAME="m-apiserver"
fi

# Create a temp dir and clean it up on exit
TEMPDIR=`mktemp -d -t m-apiserver-test.XXX`
trap "rm -rf $TEMPDIR" EXIT HUP INT QUIT TERM

# Build the Maya binary for the tests
echo "--> Building ${CTLNAME} ..."
go build -o $TEMPDIR/m-apiserver || exit 1

# Run the tests
echo "--> Running tests"
GOBIN="`which go`"
sudo -E PATH=$TEMPDIR:$PATH  -E GOPATH=$GOPATH \
    $GOBIN test ${GOTEST_FLAGS:--cover -timeout=900s} $($GOBIN list ./... | grep -v /vendor/)

