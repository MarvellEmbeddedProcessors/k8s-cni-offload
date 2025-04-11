package utils

import (
	"errors"
	"strings"

	"github.com/MarvellEmbeddedProcessors/k8s-cni-offload/pkg/logging"
	ghw "github.com/jaypipes/ghw"
)

const (
	VendorID = "177d"
	deviceID = "a0f7"
)

func Mapped_VF(pf_count int, pfid int, vfid int) (string, error) {
	/*
		Cavium, Inc. vendor id: 177d
		Ethernet Controller device id: a0f7
	*/
	logging.Info("Utils", "Mapped_VF: pf_count", pf_count, "pfid", pfid,
		"vfid", vfid)

	targetVendorID := VendorID
	targetDeviceID := deviceID
	pci, err := ghw.PCI()
	if err != nil {
		return "", err
	}

	devices := pci.Devices
	var list []string
	for _, device := range devices {
		if device.Vendor != nil && device.Product != nil {
			if device.Vendor.ID == targetVendorID &&
				device.Product.ID == targetDeviceID {
				list = append(list, device.Address)
			}
		}
	}
	dpu_vfid := pf_count*vfid + pfid
	size := len(list) - 1
	if dpu_vfid >= size {
		return "", errors.New("mapped VF out of bounds")
	}

	vf_pci := strings.Split(list[dpu_vfid+1], " ")
	logging.Info("Utils", "vf_pci", vf_pci)
	return vf_pci[0], nil
}

func LinknameByPci(addrs string) string {
	var err error

	logging.Info("Utils", "LinknameByPci: addrs", addrs)
	net, err := ghw.Network()
	if err != nil {
		return ""
	}

	for _, nic := range net.NICs {
		if nic.PCIAddress != nil && strings.Contains(*nic.PCIAddress, addrs) {
			logging.Info("Utils", "Name", nic.Name, "Address", *nic.PCIAddress)
			return nic.Name
		}
	}

	return ""
}
