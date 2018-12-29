#!/usr/bin/env bash

function cleanup {
  rm ../artifacts/container-image/cluster-apiserver
}
trap cleanup EXIT

pushd ../
cp -v cluster-apiserver ./artifacts/container-image/cluster-apiserver
docker build -t cluster-apiserver:latest ./artifacts/container-image
docker save cluster-apiserver:latest | gzip > ./artifacts/cluster-apiserver.tar.gz
docker rmi cluster-apiserver:latest
popd
