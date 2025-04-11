#!/bin/sh

# Install cp-agent
cp /scripts/bin/cn106xx.cfg /host/usr/bin/cn106xx.cfg
cp /scripts/bin/octep_cp_agent /host/usr/bin/octep_cp_agent
cp -r /scripts/libconfig-bin/* /host/usr/
cp /scripts/cp-agent/setup-dpu.sh /host/etc/init.d/setup-dpu.sh
cp /scripts/cp-agent/setup-dpu.service /host/etc/systemd/system/setup-dpu.service
chroot /host chmod 777 /etc/init.d/setup-dpu.sh

# Start cp-agent as a service
chroot /host systemctl enable setup-dpu
chroot /host systemctl start setup-dpu

export NETCONFPATH=/etc/cni/net.d/
export CNI_PATH=/opt/cni/bin/
export MARVELL_PLUGIN_CONFIG=/marvell-config.yaml
export CNI_OFFLOAD_AGENT_CONFIG=/cni-offload-config.yaml
/cniOffloadAgent

