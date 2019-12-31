#!/bin/bash
PULL_WORKSPACE=$1
PULL_BRANCH=$2
DIST_PATH=$3

echo "Pulling Source Code"
cd ${PULL_WORKSPACE}/${PULL_BRANCH}
sudo git fetch -q
sudo git checkout -q --force origin/${PULL_BRANCH}

go mod download
go build  -ldflags "-s -w" -o bin/main
md5sum bin/main

cp bin/main ${DIST_PATH}/main
echo "Finished"