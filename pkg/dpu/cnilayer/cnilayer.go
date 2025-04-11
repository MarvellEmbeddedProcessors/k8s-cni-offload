/*
This package abstracts CNI layers.
*/
package cnilayer

import (
	"strings"

	"github.com/MarvellEmbeddedProcessors/k8s-cni-offload/pkg/dpu/cnilayertypes"
	"github.com/MarvellEmbeddedProcessors/k8s-cni-offload/pkg/logging"
	offloadciliumcni "github.com/MarvellEmbeddedProcessors/k8s-cni-offload/plugin/dpu/offload-cilium-cni"
	offloaddebugcni "github.com/MarvellEmbeddedProcessors/k8s-cni-offload/plugin/dpu/offload-debug-cni"
)

type cniLayer interface {
	GetCniName() string
	Add(*cnilayertypes.OffloadCNIAddRequest) (*cnilayertypes.OffloadCNIAddResponse, error)
	Del(*cnilayertypes.OffloadCNIDelRequest) (*cnilayertypes.OffloadCNIDelResponse, error)
	Check(*cnilayertypes.OffloadCNICheckRequest) (*cnilayertypes.OffloadCNICheckResponse, error)
}

type CniLayer struct {
	layer cniLayer
}

// NewCniLayer Function returns CniLayer
func NewCniLayer(cniLayerNameArg string) *CniLayer {
	cniLayer := []cniLayer{offloadciliumcni.NewOffloadCiliumCni(),
		offloaddebugcni.NewOffloadDebugCni()}

	for _, layer := range cniLayer {
		if strings.Contains(cniLayerNameArg, layer.GetCniName()) {
			logging.Info("cnilayer: NewCniLayer() cni layer:", layer.GetCniName())
			return &CniLayer{layer}
		}
	}

	logging.Error("cnilayer: NewCniLayer(): cni layer not found:", cniLayerNameArg)
	return &CniLayer{}
}

// Add Function Calls CNI Add command
func (cniLayer *CniLayer) Add(in *cnilayertypes.OffloadCNIAddRequest) (*cnilayertypes.OffloadCNIAddResponse, error) {
	logging.Info("cnilayer: Add request:", in)
	resp, err := cniLayer.layer.Add(in)
	if err != nil {
		logging.Info("cnilayer: Add request failed:", err)
		return nil, err
	}

	logging.Info("cnilayer: Add request successful", resp, err)
	return resp, err
}

// Del Function Calls CNI Del command
func (cniLayer *CniLayer) Del(in *cnilayertypes.OffloadCNIDelRequest) (*cnilayertypes.OffloadCNIDelResponse, error) {
	logging.Info("cnilayer: Del request:", in)
	resp, err := cniLayer.layer.Del(in)
	if err != nil {
		logging.Info("cnilayer: Del request failed:", err)
		return nil, err
	}

	logging.Info("cnilayer: Del request successful", resp, err)
	return resp, err
}

// Check Function Calls CNI Check command
func (cniLayer *CniLayer) Check(in *cnilayertypes.OffloadCNICheckRequest) (*cnilayertypes.OffloadCNICheckResponse, error) {
	logging.Info("cnilayer: Check request:", in)
	resp, err := cniLayer.layer.Check(in)
	if err != nil {
		logging.Info("cnilayer: Check request failed:", err)
		return nil, err
	}

	logging.Info("cnilayer: Check request successful", resp, err)
	return resp, err
}
