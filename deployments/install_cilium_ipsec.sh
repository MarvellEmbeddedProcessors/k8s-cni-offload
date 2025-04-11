#!/bin/sh
cd build/cilium/install/kubernetes/cilium/;

# Install cilium-ipsec-keys
kubectl create -n kube-system secret generic cilium-ipsec-keys \
   --from-literal=keys="6 rfc4106(gcm(aes)) $(echo $(dd if=/dev/urandom count=20 bs=1 2> /dev/null | xxd -p -c 64)) 128"

# Install cilium with IPSec
helm install cilium --namespace kube-system --set k8sServiceHost=auto --set debug.enabled=true --set socketLB.hostNamespaceOnly=true --set encryption.enabled=true --set encryption.type=ipsec --set image.repository=$1/cilium/cilium --set image.tag=latest .
