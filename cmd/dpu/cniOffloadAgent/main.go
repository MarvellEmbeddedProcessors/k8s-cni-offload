/*
Copyright 2024 Marvell.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"errors"
	"log"
	"net"
	"os"

	pb "github.com/MarvellEmbeddedProcessors/k8s-cni-offload/api/pb/cniOffload"
	vendorplugin "github.com/MarvellEmbeddedProcessors/k8s-cni-offload/pkg/dpu/vendorplugin"
	"github.com/MarvellEmbeddedProcessors/k8s-cni-offload/pkg/logging"

	"github.com/takama/daemon"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"
)

type Config struct {
	ServerPort       string `yaml:"serverPort"`
	VendorPluginName string `yaml:"vendorPluginName"`
}

type cniOffloadServer struct {
	pb.UnimplementedCNIActionServer
	grpcServer *grpc.Server
}

var vendorPlugin *vendorplugin.VendorPlugin

// readConfig reads the config file
func readConfig() (Config, error) {
	config_path := os.Getenv("CNI_OFFLOAD_AGENT_CONFIG")
	if config_path == "" {
		return Config{}, errors.New("CNI_OFFLOAD_AGENT_CONFIG is not set")
	}

	configFile, err := os.ReadFile(config_path)
	if err != nil {
		return Config{}, err
	}

	var config Config
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		return Config{}, err
	}

	return config, err
}

// CmdAdd is the gRPC handler for the Add command
func (s *cniOffloadServer) CmdAdd(ctx context.Context, in *pb.CNIAddRequest) (*pb.CNIAddResponse, error) {
	return vendorPlugin.Plugin.Add(in)
}

// CmdDel is the gRPC handler for the Del command
func (s *cniOffloadServer) CmdDel(ctx context.Context, in *pb.CNIDelRequest) (*pb.CNIDelResponse, error) {
	return vendorPlugin.Plugin.Del(in)
}

// CmdCheck is the gRPC handler for the Check command
func (s *cniOffloadServer) CmdCheck(ctx context.Context, in *pb.CNICheckRequest) (*pb.CNICheckResponse, error) {
	return vendorPlugin.Plugin.Check(in)
}

// Listen creates a listener on the specified port
func (s *cniOffloadServer) Listen(config Config) (net.Listener, error) {
	listener, err := net.Listen("tcp", config.ServerPort)
	if err != nil {
		panic(err)
	}
	s.grpcServer = grpc.NewServer()
	// reflection.Register(s)
	pb.RegisterCNIActionServer(s.grpcServer, &cniOffloadServer{})
	logging.Info("cniOffloadAgent", "CNIOffloadAgent: server listening on Port: ", config.ServerPort)
	return listener, nil
}

// Serve starts the gRPC server
func (s *cniOffloadServer) Serve(listener net.Listener) error {
	err := s.grpcServer.Serve(listener)
	if err != nil {
		logging.Error("CNIOffloadAgent: failed to serve: ", err)
	}
	return err
}

// NewOffloadServer creates a new cniOffloadServer
func NewOffloadServer() *cniOffloadServer {
	return &cniOffloadServer{}
}

func main() {
	// read config file
	config, err := readConfig()
	if err != nil {
		logging.Error("CNIOffloadAgent: readConfig() failed:", err)
		return
	}

	logging.Info("cniOffloadAgent", "readConfig", config)

	if config.ServerPort == "" {
		logging.Error("Missing fields in config file - serverPort not set")
		return
	}

	if config.VendorPluginName == "" {
		logging.Error("Missing fields in config file - vendorPluginName not set")
		return
	}

	cniOffloadServer := NewOffloadServer()

	// Get vendor plugin
	vendorPlugin, err = vendorplugin.NewVendorPlugin(config.VendorPluginName)
	if err != nil {
		logging.Error("CNIOffloadAgent: NewVendorPlugin() failed: ", err)
		return
	}

	// Init vendor plugin
	err = vendorPlugin.InitFn()
	if err != nil {
		logging.Error("CNIOffloadAgent: failed to init cni layer: ", err)
		return
	}

	listener, err := cniOffloadServer.Listen(config)
	if err != nil {
		logging.Error("CNIOffloadAgent: failed to listen: ", err)
		return
	}
	logging.Info("cniOffloadAgent", "Listen", config)

	err = cniOffloadServer.Serve(listener)
	if err != nil {
		logging.Error("CNIOffloadAgent: failed to serve: ", err)
		return
	}
	logging.Info("cniOffloadAgent", "serve", "")

	service, err := daemon.New("name", "offload-agent", daemon.SystemDaemon)
	if err != nil {
		logging.Error("CNIOffloadAgent: failed to create daemon: ", err)
		return
	}
	status, err := service.Install()
	if err != nil {
		log.Fatal(status, "\nError: ", err)

	}
	os.Exit(1)
}
