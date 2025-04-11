/*
This package abstracts vendor plugins.
*/
package vendorplugin

import (
	"fmt"
	"strings"

	pb "github.com/MarvellEmbeddedProcessors/k8s-cni-offload/api/pb/cniOffload"
	"github.com/MarvellEmbeddedProcessors/k8s-cni-offload/pkg/logging"
	dbgpl "github.com/MarvellEmbeddedProcessors/k8s-cni-offload/plugin/dpu/vendorplugins/debug"
	mrvlpl "github.com/MarvellEmbeddedProcessors/k8s-cni-offload/plugin/dpu/vendorplugins/marvell"
)

type vendorPlugin interface {
	GetVendorPluginName() string
	InitFn() error
	Add(*pb.CNIAddRequest) (*pb.CNIAddResponse, error)
	Del(*pb.CNIDelRequest) (*pb.CNIDelResponse, error)
	Check(*pb.CNICheckRequest) (*pb.CNICheckResponse, error)
}

type VendorPlugin struct {
	Plugin vendorPlugin
}

// NewVendorPlugin Function returns VendorPlugin
func NewVendorPlugin(vendorPluginName string) (*VendorPlugin, error) {
	logging.Info("vendorplugin: NewVendorPlugin(): vendorPluginName:", vendorPluginName)
	var vendorpl vendorPlugin
	vendorplugins := []vendorPlugin{mrvlpl.NewMarvellPlugin(), dbgpl.NewDebugPlugin()}
	for _, plugin := range vendorplugins {
		if strings.Contains(vendorPluginName, plugin.GetVendorPluginName()) {
			vendorpl = plugin
			break
		}
	}
	if vendorpl == nil {
		return nil, fmt.Errorf("vendorplugin: NewVendorPlugin(): Plugin not found: vendorPluginName %v",
			vendorPluginName)
	}

	logging.Info("vendorplugin: configured vendorPluginName:", vendorPluginName)
	return &VendorPlugin{vendorpl}, nil
}

// InitFn Function initializes the vendor plugin/vendor
func (vendorpl *VendorPlugin) InitFn() error {
	logging.Info("vendorplugin: Plugin Init")
	if vendorpl == nil || vendorpl.Plugin == nil {
		return fmt.Errorf("vendorpl.Plugin Not set")
	}

	err := vendorpl.Plugin.InitFn()
	if err != nil {
		logging.Info("vendorplugin: Plugin Init failed:", err)
		return err
	}

	logging.Info("vendorplugin: Plugin Init successful")
	return nil
}

// Add Function Calls Vendor Add command for the plugin/vendor
func (vendorpl *VendorPlugin) Add(in *pb.CNIAddRequest) (*pb.CNIAddResponse, error) {
	logging.Info("vendorplugin: Add request:", in)
	if vendorpl == nil || vendorpl.Plugin == nil {
		return nil, fmt.Errorf("vendorpl.Plugin Not set")
	}

	AddResp, err := vendorpl.Plugin.Add(in)
	if err != nil {
		logging.Info("vendorplugin: Add request failed:", err)
		return nil, err
	}

	logging.Info("vendorplugin: Add request successful:", AddResp)
	return AddResp, nil
}

// Check Function Calls Vendor Check command for the plugin/vendor
func (vendorpl *VendorPlugin) Check(in *pb.CNICheckRequest) (*pb.CNICheckResponse, error) {
	logging.Info("vendorplugin: Check request:", in)
	if vendorpl == nil || vendorpl.Plugin == nil {
		return nil, fmt.Errorf("vendorpl.Plugin Not set")
	}

	CheckResp, err := vendorpl.Plugin.Check(in)
	if err != nil {
		logging.Info("vendorplugin: Check request failed:", err)
		return nil, err
	}

	logging.Info("vendorplugin: Check request successful:", CheckResp)
	return CheckResp, nil
}

// Del Function Calls Vendor Del command for the plugin/vendor
func (vendorpl *VendorPlugin) Del(in *pb.CNIDelRequest) (*pb.CNIDelResponse, error) {
	logging.Info("vendorplugin: Del request:", in)
	if vendorpl == nil || vendorpl.Plugin == nil {
		return nil, fmt.Errorf("vendorpl.Plugin Not set")
	}

	DelResp, err := vendorpl.Plugin.Del(in)
	if err != nil {
		logging.Info("vendorplugin: Del request failed:", err)
		return nil, err
	}

	logging.Info("vendorplugin: Del request successful:", DelResp)
	return DelResp, nil
}
