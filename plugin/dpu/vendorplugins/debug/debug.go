package debug

//Imports Protobuf as pb
import (
	"github.com/MarvellEmbeddedProcessors/k8s-cni-offload/pkg/logging"
)

type DebugPlugin struct {
	Name string
}

func NewDebugPlugin() *DebugPlugin {
	return &DebugPlugin{Name: "Debug Plugin"}
}

func (d *DebugPlugin) GetVendorPluginName() string {
	return "debug"
}

func (d *DebugPlugin) InitFn() error {
	if err := HostIntfCni_InitFn(); err != nil {
		logging.Error("Debug plugin: HostIntfCni_InitFn() failed", err)
		return err
	}

	logging.Info("Debug plugin: Dataplane init successful")
	return nil
}
