#!/bin/bash

set -x

# Total memory need by octep_cp_agent: 256Mb
# Hugepage sizes supported: 32Mb, 512Mb, 1Gb
setup_hugepages() {
    if ! mount | grep -q "^none on /dev/huge type hugetlbfs " ; then
        mkdir -p /dev/huge
        mount -t hugetlbfs none /dev/huge
        echo 8 > /proc/sys/vm/nr_hugepages
    fi
}

setup_sdk12_24_11_cpagent() {
	modprobe pcie_marvell_cnxk_ep
	exec /usr/bin/octep_cp_agent /usr/bin/cn106xx.cfg 2>&1 > /tmp/octep-cp-log.txt &
}

setup_sdk12_25_03_cpagent() {
    setup_hugepages
    modprobe vfio-pci
    rvu_pf_1=$(lspci -d 177d:a0ef -n | awk 'NR==1{print $1}')
    rvu_pf_2=$(lspci -d 177d:a0ef -n | awk 'NR==2{print $1}')
    if [ -z $rvu_pf_2 ]
    then
        echo "Second RVU PF device not found"
        return 1
    fi

    echo "Second RVU PF device: $rvu_pf_2"

    rvu1drv=$(lspci -ks $rvu_pf_1 | grep "driver in use" | awk 'NR==1{print $5}')
    echo $rvu_pf_1 > /sys/bus/pci/drivers/$rvu1drv/unbind
    echo $rvu_pf_1 > /sys/bus/pci/devices/$rvu_pf_1/driver/unbind

    rvu2drv=$(lspci -ks $rvu_pf_2 | grep "driver in use" | awk 'NR==1{print $5}')
    echo $rvu_pf_2 > /sys/bus/pci/drivers/$rvu2drv/unbind
    echo $rvu_pf_2 > /sys/bus/pci/devices/$rvu_pf_2/driver/unbind

    pem=$(lspci -d 177d:a06c -n | awk 'NR==1{print $1}')
    if [ -z $pem ]
    then
        echo "PEM device not found"
        return 1
    fi

    echo "PEM device: $pem"

    pemdrv=$(lspci -ks $pem | grep "driver in use" | awk 'NR==1{print $5}')

    echo $pem > /sys/bus/pci/drivers/$pemdrv/unbind
    echo $pem > /sys/bus/pci/devices/$pem/driver/unbind

    echo vfio-pci > /sys/bus/pci/devices/$rvu_pf_1/driver_override && echo $rvu_pf_1 > /sys/bus/pci/drivers_probe
    echo vfio-pci > /sys/bus/pci/devices/$rvu_pf_2/driver_override && echo $rvu_pf_2 > /sys/bus/pci/drivers_probe
    echo vfio-pci > /sys/bus/pci/devices/$pem/driver_override && echo $pem > /sys/bus/pci/drivers_probe

    exec /usr/bin/octep_cp_agent /usr/bin/cn106xx.cfg -- --sdp_rvu_pf $rvu_pf_1,$rvu_pf_2 --pem_dev $pem &> /tmp/octep-cp-log.txt &

    return 0

}

setup_sdk12_25_01_cpagent() {
    setup_hugepages
    modprobe vfio-pci
    dpi=$(lspci -d 177d:a080 -n | awk 'NR==1{print $1}')
    if [ -z $dpi ]
    then
        echo "DPI device not found"
        return 1
    fi

    echo "DPI device: $dpi"

    dpidrv=$(lspci -ks $dpi | grep "driver in use" | awk 'NR==1{print $5}')

    pem=$(lspci -d 177d:a06c -n | awk 'NR==1{print $1}')
    if [ -z $pem ]
    then
        echo "PEM device not found"
        return 1
    fi

    echo "PEM device: $pem"

    pemdrv=$(lspci -ks $pem | grep "driver in use" | awk 'NR==1{print $5}')

    if [ $? -ne 0 ]; then
        echo $dpi > /sys/bus/pci/drivers/$dpidrv/unbind
    fi

    echo $dpi > /sys/bus/pci/devices/$dpi/driver/unbind

    echo vfio-pci > /sys/bus/pci/devices/$dpi/driver_override && echo $dpi > /sys/bus/pci/drivers_probe

    if [ $? -ne 0 ]; then
        echo $pem > /sys/bus/pci/drivers/$pemdrv/unbind
    fi

    echo $pem > /sys/bus/pci/devices/$pem/driver/unbind

    echo vfio-pci > /sys/bus/pci/devices/$pem/driver_override && echo $pem > /sys/bus/pci/drivers_probe

    exec /usr/bin/octep_cp_agent /usr/bin/cn106xx.cfg -- --dpi_dev $dpi --pem_dev $pem &> /tmp/octep-cp-log.txt &

    return 0
}

run() {
	swapoff -a

	# Refer to pcie_ep_octeon_target to compile and run octep_cp_agent
	setup_sdk12_25_01_cpagent

	# Marvell vendor ID: 177d, SDP VF device ID: a0f7
	# Based on the card type, may need to change the device id type
	vendorID=177d
	deviceID=a0f7
	devstr=$vendorID:$deviceID

	pf=$(lspci -d $devstr -n | awk 'NR==1{print $1}')

	echo $pf

	path="/sys/bus/pci/devices/$pf/net/"
	echo $path
	ifname=$(ls $path)
	echo $ifname
	ifconfig $ifname 192.168.1.1/24 up

	# Set VFs to Network Manager unmanaged
	list=$(lspci -d $devstr -n | awk '{print $1}' |  tail -n +2)
	for f in $list; do
	        path="/sys/bus/pci/devices/$f/net/"
	        ifname=$(ls $path)
	        echo $ifname
	        nmcli dev set $ifname managed no
	        i=$((i+1))
	done
}

run
