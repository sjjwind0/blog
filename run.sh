#!/bin/bash
current_go_path=$GOPATH
current_path=$(cd `dirname $0`; pwd)
GOPATH=$GOPATH":"$current_path
rm -rf ./src/pkg
echo 'begin building ...'
go build ./src/main.go
GOPATH=$current_go_path
if [ $? -eq 0 ]; then
	echo 'run project ...'
	sudo ./main
fi