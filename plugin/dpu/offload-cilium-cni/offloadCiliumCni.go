package offloadCiliumCni

import (
	"context"
	"encoding/json"
	"github.com/MarvellEmbeddedProcessors/k8s-cni-offload/plugin/dpu/vendorplugins/marvell/utils"
	"os"
	"path/filepath"

	cnilayertypes "github.com/MarvellEmbeddedProcessors/k8s-cni-offload/pkg/dpu/cnilayertypes"
	"github.com/MarvellEmbeddedProcessors/k8s-cni-offload/pkg/logging"

	"github.com/containernetworking/cni/libcni"
	cniTypesV1 "github.com/containernetworking/cni/pkg/types/100"
	"github.com/containernetworking/plugins/pkg/ns"
)

const (
	EnvCNIPath        = "CNI_PATH"
	EnvNetDir         = "NETCONFPATH"
	EnvCapabilityArgs = "CAP_ARGS"
	EnvCNIArgs        = "CNI_ARGS"
	EnvCNIIfname      = "CNI_IFNAME"

	DefaultNetDir = "/etc/cni/net.d"

	CmdAdd   = "add"
	CmdCheck = "check"
	CmdDel   = "del"
)

type OffloadCiliumCni struct {
	Name string
}

func NewOffloadCiliumCni() *OffloadCiliumCni {
	return &OffloadCiliumCni{Name: "Offload Cilium CNI"}
}

type CNIDetails_Cilium struct {
	PodName      string `json:"PodName,omitempty"`
	PodNamespace string `json:"PodNamespace,omitempty"`
	MacAddr      string `json:"MacAddr,omitempty"`
	ContainerID  string `json:"ContainerID,omitempty"`
	VfId         uint32 `json:"VfId,omitempty"`
	PfId         uint32 `json:"PfId,omitempty"`
	PfCount      uint32 `json:"PfCount,omitempty"`
}

// parseArgs parses the CNI_ARGS environment variable into a slice of key-value pairs.
func parseArgs(args string) ([][2]string, error) {
	var cniArgs [][2]string
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(args), &data); err != nil {
		return nil, err
	}
	K8S_POD_NAME, _ := data["K8S_POD_NAME"].(string)
	K8S_POD_NAMESPACE, _ := data["K8S_POD_NAMESPACE"].(string)
	cniArgs = [][2]string{
		{"K8S_POD_NAME", K8S_POD_NAME},
		{"K8S_POD_NAMESPACE", K8S_POD_NAMESPACE},
	}
	return cniArgs, nil
}

// cniInvoke invokes the CNI plugin with the given command and CNI details.
func cniInvoke(cmdAPI string, in *CNIDetails_Cilium) (*cniTypesV1.Result, error) {
	netdir := os.Getenv(EnvNetDir)
	if netdir == "" {
		netdir = DefaultNetDir
	}
	netconf, err := libcni.LoadConfList(netdir, "cep")
	if err != nil {
		logging.Error("offloadCiliumCni: Error while loading the CNI configuration:", err)
		return nil, err
	}

	VfPCI, err := utils.Mapped_VF(int(in.PfCount), int(in.PfId), int(in.VfId))
	if err != nil {
		logging.Error("offloadCiliumCni: Error while getting the VF PCI:", err)
		return nil, err
	}
	ifName := utils.LinknameByPci(VfPCI)
	logging.Info("offloadCiliumCni: The VF PCI is:", VfPCI, "ifName", ifName)
	netns, err := ns.GetCurrentNS()
	if err != nil {
		logging.Error("offloadCiliumCni: Error while getting the current netns:", err)
		return nil, err
	}
	if cmdAPI == "add" {
		netconfData := map[string]interface{}{
			"MacAddr": in.MacAddr,
			"VfPCI":   VfPCI,
		}
		netconf.Plugins[0], err = libcni.InjectConf(netconf.Plugins[0], netconfData)
	}

	var capabilityArgs map[string]interface{}
	capabilityArgsValue := os.Getenv(EnvCapabilityArgs)
	if len(capabilityArgsValue) > 0 {
		if err := json.Unmarshal([]byte(capabilityArgsValue), &capabilityArgs); err != nil {
			logging.Error("offloadCiliumCni: Error while unmarshalling the capability args:", err)
			return nil, err
		}
	}
	var cniArgs [][2]string
	cniArgsdata := map[string]interface{}{
		"K8S_POD_NAME":      in.PodName,
		"K8S_POD_NAMESPACE": in.PodNamespace,
	}
	jsoncniArgsData, err := json.Marshal(cniArgsdata)
	if err != nil {
		logging.Error("offloadCiliumCni: Error while marshalling the cniArgsData:", err)
		return nil, err
	}
	args := string(jsoncniArgsData)
	cniArgs, err = parseArgs(args)
	if err != nil {
		logging.Error("offloadCiliumCni: Error while running parseArgs(args):", err)
		return nil, err
	}
	logging.Info("offloadCiliumCni: The cniArgs Value are:", cniArgs)
	if err != nil {
		logging.Error("offloadCiliumCni: Error while running parseArgs(args):", err)
		return nil, err
	}
	cninet := libcni.NewCNIConfig(filepath.SplitList(os.Getenv(EnvCNIPath)), nil)

	rt := &libcni.RuntimeConf{
		ContainerID:    in.ContainerID,
		NetNS:          netns.Path(),
		IfName:         ifName,
		Args:           cniArgs,
		CapabilityArgs: capabilityArgs,
	}

	switch cmdAPI {
	case CmdAdd:
		result, err := cninet.AddNetworkList(context.TODO(), netconf, rt)
		if err != nil {
			logging.Error("offloadCiliumCni: Error while calling AddNetworkList():", err)
			return nil, err
		}
		res := cniTypesV1.Result{}
		data, err := json.Marshal(result)
		if err != nil {
			logging.Error("offloadCiliumCni: Error while marshalling the result:", err)
			return nil, err
		}
		if err := json.Unmarshal(data, &res); err != nil {
			logging.Error("offloadCiliumCni: Error while unmarshalling the data:", err)
			return nil, err
		}

		logging.Info("offloadCiliumCni: The Add response is:", res)
		return &res, nil
	case CmdCheck:
		err := cninet.CheckNetworkList(context.TODO(), netconf, rt)
		if err != nil {
			logging.Error("offloadCiliumCni: Error while calling CheckNetworkList():", err)
			return nil, err
		}
		return nil, err
	case CmdDel:
		err := cninet.DelNetworkList(context.TODO(), netconf, rt)
		if err != nil {
			logging.Error("offloadCiliumCni: Error while calling DelNetworkList():", err)
			return nil, err
		}
		return nil, err
	}
	return nil, nil
}

// DelCiliumCNI deletes the CNI interface for CiliumCNI.
func (occ *OffloadCiliumCni) Del(in *cnilayertypes.OffloadCNIDelRequest) (*cnilayertypes.OffloadCNIDelResponse, error) {
	logging.Info("offloadCiliumCni: The Delete Interface call triggered", in)
	CNIDeldata := CNIDetails_Cilium{
		ContainerID:  in.PodInfo.ContainerId,
		PodName:      in.PodInfo.PodName,
		PodNamespace: in.PodInfo.PodNamespace,
		VfId:         in.IntfInfo.Vfid,
		PfId:         in.IntfInfo.Pfid,
		PfCount:      in.IntfInfo.NumPfs,
	}
	_, err := cniInvoke("del", &CNIDeldata)
	if err != nil {
		logging.Error("offloadCiliumCni: There is error in Del Call of CNI:", err)
		return nil, err
	}
	return &cnilayertypes.OffloadCNIDelResponse{}, nil
}

// CheckCiliumCNI checks the CNI interface for CiliumCNI.
func (occ *OffloadCiliumCni) Check(in *cnilayertypes.OffloadCNICheckRequest) (*cnilayertypes.OffloadCNICheckResponse, error) {
	logging.Info("offloadCiliumCni: The Check Interface call triggered", in)
	CNICheckData := CNIDetails_Cilium{
		ContainerID:  in.PodInfo.ContainerId,
		PodName:      in.PodInfo.PodName,
		PodNamespace: in.PodInfo.PodNamespace,
		VfId:         in.IntfInfo.Vfid,
		PfId:         in.IntfInfo.Pfid,
		PfCount:      in.IntfInfo.NumPfs,
	}
	_, err := cniInvoke("check", &CNICheckData)
	if err != nil {
		logging.Error("offloadCiliumCni: There is error in Check Call of CNI:", err)
		return nil, err
	}
	return &cnilayertypes.OffloadCNICheckResponse{}, nil
}

// AddCiliumCNI adds the CNI interface for CiliumCNI.
func (occ *OffloadCiliumCni) Add(in *cnilayertypes.OffloadCNIAddRequest) (*cnilayertypes.OffloadCNIAddResponse, error) {
	logging.Info("offloadCiliumCni: The Add Interface call triggered", in)
	data := CNIDetails_Cilium{
		ContainerID:  in.PodInfo.ContainerId,
		PodName:      in.PodInfo.PodName,
		PodNamespace: in.PodInfo.PodNamespace,
		MacAddr:      in.MacAddr,
		VfId:         in.IntfInfo.Vfid,
		PfId:         in.IntfInfo.Pfid,
		PfCount:      in.IntfInfo.NumPfs,
	}
	resOut, err := cniInvoke("add", &data)
	if err != nil {
		logging.Error("offloadCiliumCni: There is error in Add Call of CNI:", err)
		return nil, err
	}
	OutIps := resOut.IPs[0].Address.IP
	Outgateway := resOut.IPs[0].Gateway
	OutIp_str := OutIps.String()
	Outgateway_str := Outgateway.String()
	OffloadCNIAddResp := &cnilayertypes.OffloadCNIAddResponse{
		IpInfo: &cnilayertypes.IpDetails{
			Version: "4",
			Address: OutIp_str,
			Gateway: Outgateway_str,
			Dns:     "",
		},
	}
	logging.Info("offloadCiliumCni: AddCiliumCNI(): OffloadCNIAddRespOut", OffloadCNIAddResp)
	return OffloadCNIAddResp, nil
}

func (occ *OffloadCiliumCni) GetCniName() string {
	return "Offload Cilium CNI"
}
