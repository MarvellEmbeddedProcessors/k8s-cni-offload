package marvell

//Imports Protobuf as pb
import (
	"errors"
	"fmt"
	"os"

	pb "github.com/MarvellEmbeddedProcessors/k8s-cni-offload/api/pb/cniOffload"
	cnilayer "github.com/MarvellEmbeddedProcessors/k8s-cni-offload/pkg/dpu/cnilayer"
	cnilayertypes "github.com/MarvellEmbeddedProcessors/k8s-cni-offload/pkg/dpu/cnilayertypes"
	"github.com/MarvellEmbeddedProcessors/k8s-cni-offload/pkg/logging"

	"gopkg.in/yaml.v2"
)

type MarvellConfig struct {
	DpdkPath      string `yaml:"dpdkPath"`
	CiliumCNIPath string `yaml:"ciliumCNIPath"`
	CniLayer      string `yaml:"cniLayer"`
	LogLevel      string `yaml:"logLevel,omitempty"`
	LogFile       string `yaml:"logFile,omitempty"`
}

type MarvellPlugin struct {
	Name     string
	cnilayer *cnilayer.CniLayer
}

func NewMarvellPlugin() *MarvellPlugin {
	return &MarvellPlugin{Name: "Marvell"}
}

func readConfig() (MarvellConfig, error) {
	config := MarvellConfig{}

	// read config file
	config_path := os.Getenv("MARVELL_PLUGIN_CONFIG")
	if config_path == "" {
		return config, errors.New("config file not found")
	}

	configFile, err := os.ReadFile(config_path)
	if err != nil {
		return config, fmt.Errorf("error reading config file: %w", err)
	}

	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		return config, fmt.Errorf("error unmarshaling config file: %w", err)
	}

	return config, nil
}

// GetVendorPluginName returns the name of the plugin
func (m *MarvellPlugin) GetVendorPluginName() string {
	return "Marvell"
}

// InitFn initializes the Marvell plugin and choose the dataPath between DPDK and CiliumCNI
func (m *MarvellPlugin) InitFn() error {
	var err error
	config, err := readConfig()
	if err != nil {
		logging.Error("Marvell Plugin: readConfig() failed", err)
		return err
	}

	logging.Info("Marvell Plugin: readConfig()", config)
	m.cnilayer = cnilayer.NewCniLayer(config.CniLayer)
	logging.Info("Marvell Plugin: InitFn() complete")
	return nil
}

// Add adds the interface to the dataplane using either DPDK or CiliumCNI
func (m *MarvellPlugin) Add(in *pb.CNIAddRequest) (*pb.CNIAddResponse, error) {
	var out pb.CNIAddResponse = pb.CNIAddResponse{}

	logging.Info("Marvell Plugin: Add request: ", in)
	AddReq := &cnilayertypes.OffloadCNIAddRequest{}
	if in.PodInfo != nil {
		AddReq.PodInfo = &cnilayertypes.PodDetails{
			PodNamespace: in.PodInfo.PodNamespace,
			PodName:      in.PodInfo.PodName,
			ContainerId:  in.PodInfo.ContainerId,
		}
	}

	if in.IntfInfo != nil {
		AddReq.IntfInfo = &cnilayertypes.IntfDetails{
			Name:   in.IntfInfo.Name,
			Vfid:   in.IntfInfo.Vfid,
			Pfid:   in.IntfInfo.Pfid,
			NumPfs: in.IntfInfo.NumPfs,
		}
	}

	if in.MacAddr != "" {
		AddReq.MacAddr = in.MacAddr
	}

	if in.IpamAllocatedbyHost != nil {
		AddReq.IpamAllocatedbyHost = in.IpamAllocatedbyHost.IpamAllocatedbyHost
	}

	if in.OffloadFlags != nil {
		AddReq.OffloadFlags = &cnilayertypes.OffloadFlags{
			CheckSumOffload: in.OffloadFlags.CheckSumOffload,
		}
	}

	if in.IpInfo != nil {
		AddReq.IpInfo = &cnilayertypes.IpDetails{
			Version: in.IpInfo.Version,
			Address: in.IpInfo.Address,
			Gateway: in.IpInfo.Gateway,
			Dns:     in.IpInfo.Dns,
		}
	}

	AddResp, err := m.cnilayer.Add(AddReq)
	if err != nil {
		logging.Info("Marvell Plugin: Add request failed", AddResp, err)
		return &pb.CNIAddResponse{}, err
	}

	if AddResp.IpInfo != nil {
		out.IpInfo = &pb.IpDetails{
			Version: AddResp.IpInfo.Version,
			Address: AddResp.IpInfo.Address,
			Gateway: AddResp.IpInfo.Gateway,
			Dns:     AddResp.IpInfo.Dns,
		}
	}

	out.Status = AddResp.Status
	logging.Info("Marvell Plugin: Add request successful", out.Status, err)
	return &out, err
}

// Check checks the interface in the dataplane using either DPDK or CiliumCNI
func (m *MarvellPlugin) Check(in *pb.CNICheckRequest) (*pb.CNICheckResponse, error) {
	logging.Info("Marvell Plugin: Check request:", in)
	CheckReq := &cnilayertypes.OffloadCNICheckRequest{
		PodInfo: &cnilayertypes.PodDetails{
			PodNamespace: in.PodInfo.PodNamespace,
			PodName:      in.PodInfo.PodName,
			ContainerId:  in.PodInfo.ContainerId,
		},
		IntfInfo: &cnilayertypes.IntfDetails{
			Name:   in.IntfInfo.Name,
			Vfid:   in.IntfInfo.Vfid,
			Pfid:   in.IntfInfo.Pfid,
			NumPfs: in.IntfInfo.NumPfs,
		},
	}

	CheckResp, err := m.cnilayer.Check(CheckReq)
	if err != nil {
		logging.Info("Marvell Plugin: Check request failed", CheckResp, err)
		return &pb.CNICheckResponse{}, err
	}

	out := pb.CNICheckResponse{Status: CheckResp.Status}
	logging.Info("Marvell Plugin: Check request successful", out.Status, err)
	return &out, err
}

// Del deletes the interface in the dataplane using either DPDK or CiliumCNI
func (m *MarvellPlugin) Del(in *pb.CNIDelRequest) (*pb.CNIDelResponse, error) {
	logging.Info("Marvell Plugin: Del request:", in)
	DelReq := &cnilayertypes.OffloadCNIDelRequest{
		PodInfo: &cnilayertypes.PodDetails{
			PodNamespace: in.PodInfo.PodNamespace,
			PodName:      in.PodInfo.PodName,
			ContainerId:  in.PodInfo.ContainerId,
		},
		IntfInfo: &cnilayertypes.IntfDetails{
			Name:   in.IntfInfo.Name,
			Vfid:   in.IntfInfo.Vfid,
			Pfid:   in.IntfInfo.Pfid,
			NumPfs: in.IntfInfo.NumPfs,
		},
	}

	DelResp, err := m.cnilayer.Del(DelReq)
	if err != nil {
		logging.Info("Marvell Plugin: Del request failed", DelResp, err)
		return &pb.CNIDelResponse{}, err
	}

	out := pb.CNIDelResponse{Status: DelResp.Status}
	logging.Info("Marvell Plugin: Del request successful", out.Status, err)
	return &out, err
}
