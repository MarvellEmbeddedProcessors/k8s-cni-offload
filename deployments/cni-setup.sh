#!/bin/sh

# Deletes CNI config files, but back up and restore is not supported yet
rm -fr /etc/cni/net.d/*
rm -fr /opt/cni/bin/offload-cni
cp /offload-cni /opt/cni/bin/offload-cni
cp /20-offload-cni.conf /etc/cni/net.d/
while true; do sleep 3600; done;
