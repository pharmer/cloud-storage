#!/usr/bin/env bash

pushd $GOPATH/src/github.com/pharmer/cloud-storage/hack/gendocs
go run main.go
popd
