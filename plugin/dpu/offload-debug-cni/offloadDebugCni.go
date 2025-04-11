package offloadDebugCni

//Imports Protobuf as pb
import (
	cnilayertypes "github.com/MarvellEmbeddedProcessors/k8s-cni-offload/pkg/dpu/cnilayertypes"
	"github.com/MarvellEmbeddedProcessors/k8s-cni-offload/pkg/logging"
)

type OffloadDebugCni struct {
	Name string
}

func NewOffloadDebugCni() *OffloadDebugCni {
	return &OffloadDebugCni{Name: "Offload Debug CNI"}
}

// Add hot plugs a device and adds that port into bridge DPDK dataplane
// It takes pb.OffloadCNIAddRequest
func (odc *OffloadDebugCni) Add(in *cnilayertypes.OffloadCNIAddRequest) (*cnilayertypes.OffloadCNIAddResponse, error) {
	logging.Info("offloadDebugCni: Add request successful:", in)
	return &cnilayertypes.OffloadCNIAddResponse{}, nil
}

// Del deletes a bridge port
// It takes pb.OffloadCNIDelRequest
func (odc *OffloadDebugCni) Del(in *cnilayertypes.OffloadCNIDelRequest) (*cnilayertypes.OffloadCNIDelResponse, error) {
	logging.Info("offloadDebugCni: Del request successful:", in)
	return &cnilayertypes.OffloadCNIDelResponse{}, nil
}

// Check gets the information of a port
// It takes pb.OffloadCNICheckRequest
func (odc *OffloadDebugCni) Check(in *cnilayertypes.OffloadCNICheckRequest) (*cnilayertypes.OffloadCNICheckResponse, error) {
	logging.Info("offloadDebugCni: Check request successful:", in)
	return &cnilayertypes.OffloadCNICheckResponse{}, nil
}

func (odc *OffloadDebugCni) GetCniName() string {
	return "Offload Debug CNI"
}
