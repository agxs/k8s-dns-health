#!/bin/bash

addresses=${ADDRESSES:-kubernetes.default,google.com}
failureAction=${FAILURE_ACTION:-teams}
kubeconfig=${KUBECONFIG:-$HOME/.kube/config}
label=${LABEL:-k8s-app=kube-dns}
namespace=${NAMESPACE:-kube-system}
port=${PORT:-53}
server=${SERVER:-192.168.0.1}
testType=${TEST_TYPE:-server}

exec /dns-health \
  -addresses ${addresses} \
  -failureAction ${failureAction} \
  -kubeconfig ${kubeconfig} \
  -label ${label} \
  -namespace ${namespace} \
  -port ${port} \
  -server ${server} \
  -testType ${testType}
