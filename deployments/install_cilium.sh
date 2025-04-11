#!/bin/sh
cd build/cilium/install/kubernetes/cilium/;

helm install cilium --namespace kube-system --set kubeProxyReplacement=true --set bgp.enabled=true   --set bgp.announce.loadbalancerIP=true   --set bgp.announce.podCIDR=true --set debug.enabled=true --set k8sServiceHost=auto --set k8sServicePort=6443 --set socketLB.hostNamespaceOnly=true --set image.repository=$1/cilium/cilium --set image.tag=latest .
