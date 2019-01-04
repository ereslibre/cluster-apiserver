#!/usr/bin/env bash

DIR=$(dirname "${BASH_SOURCE[0]}")
ARTIFACTS_PATH="${DIR}/../artifacts/local-server"
CERTS_PATH="${ARTIFACTS_PATH}/certs"

mkdir -p $CERTS_PATH

if [ ! -f ${CERTS_PATH}/ca.crt ]; then
  openssl req -nodes -new -x509 -keyout ${CERTS_PATH}/ca.key -out ${CERTS_PATH}/ca.crt
fi

if [ ! -f ${CERTS_PATH}/client.crt ]; then
  openssl req -out ${CERTS_PATH}/client.csr -new -newkey rsa:4096 -nodes -keyout ${CERTS_PATH}/client.key -subj "/CN=development/O=system:masters"
  openssl x509 -req -days 3650 -in ${CERTS_PATH}/client.csr -CA ${CERTS_PATH}/ca.crt -CAkey ${CERTS_PATH}/ca.key -set_serial 01 -out ${CERTS_PATH}/client.crt
fi

docker run -d --name=etcd --net=host --rm k8s.gcr.io/etcd:3.2.24 etcd --advertise-client-urls=http://0.0.0.0:2379 &> /dev/null

${DIR}/../cluster-apiserver --secure-port 6443 --etcd-servers http://127.0.0.1:2379 --client-ca-file ${CERTS_PATH}/ca.crt --kubeconfig ${ARTIFACTS_PATH}/kubeconfig --authentication-kubeconfig ${ARTIFACTS_PATH}/kubeconfig --authorization-kubeconfig ${ARTIFACTS_PATH}/kubeconfig --authentication-tolerate-lookup-failure
