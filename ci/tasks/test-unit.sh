#!/bin/bash -eux

mkdir -p $GOPATH/src/github.com/cunnie/
cp -Rp u2date/ /go/src/github.com/cunnie/u2date
cd $GOPATH/src/github.com/cunnie/u2date
ginkgo -v -r .
