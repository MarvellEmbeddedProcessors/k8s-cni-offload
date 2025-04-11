package debug

//Imports Protobuf as pb
import (
	pb "github.com/MarvellEmbeddedProcessors/k8s-cni-offload/api/pb/cniOffload"
	"github.com/MarvellEmbeddedProcessors/k8s-cni-offload/pkg/logging"
)

func HostIntfCni_InitFn() error {
	logging.Info("Debug Plugin: HostIntfCni_InitFn() called")
	return nil
}

// Add hot plugs a device and adds that port into bridge DPDK dataplane
// It takes pb.CNIAddRequest
func (d *DebugPlugin) Add(in *pb.CNIAddRequest) (*pb.CNIAddResponse, error) {
	logging.Info("Debug Plugin: Add request successful:", in)
	return &pb.CNIAddResponse{}, nil
}

// Del deletes a bridge port
// It takes pb.CNIDelRequest
func (d *DebugPlugin) Del(in *pb.CNIDelRequest) (*pb.CNIDelResponse, error) {
	logging.Info("Debug Plugin: Del request successful:", in)
	return &pb.CNIDelResponse{}, nil
}

// Check gets the information of a port
// It takes pb.CNICheckRequest
func (d *DebugPlugin) Check(in *pb.CNICheckRequest) (*pb.CNICheckResponse, error) {
	logging.Info("Debug Plugin: Check request successful:", in)
	return &pb.CNICheckResponse{}, nil
}
