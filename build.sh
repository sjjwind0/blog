#!/bin/bash
current_go_path=$GOPATH
current_path=$(cd `dirname $0`; pwd)
GOPATH=$GOPATH":"$current_path
echo $GOPATH
rm -rf ./src/pkg
go build ./src/main.go
GOPATH=$current_go_path

