package main

import (
	"context"
	pb "github.com/MarvellEmbeddedProcessors/k8s-cni-offload/api/pb/cniOffload"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	dest_ip := "192.168.1.1:8000"
	conn, err := grpc.Dial(dest_ip, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	defer conn.Close()
	client := pb.NewCNIActionClient(conn)
	sampleAddReq := getSampleAddReq()
	CNIAddRespOut, err := client.CmdAdd(context.Background(), sampleAddReq)
	if err != nil {
		log.Fatalf("\nError in CNIAddRequest:%v\n", err)
	}
	log.Printf("The Output of CNIAddRequst is :%v\n", CNIAddRespOut)

	sampleDelReq := getSampleDelReq()
	CNIDelRespOut, err := client.CmdDel(context.Background(), sampleDelReq)
	if err != nil {
		log.Fatalf("\nError in CNIDelRequest:%v\n", err)
	}
	log.Printf("The Output of CNIDelRequst is :%v\n", CNIDelRespOut)

}

func getSampleAddReq() *pb.CNIAddRequest {
	sampleAddReq := &pb.CNIAddRequest{
		PodInfo: &pb.PodDetails{
			// Should not fill PodNamespace and PodName w/o actual pod
			ContainerId: "cnitool-77383ca0a0715733ca6a",
		},
		IntfInfo: &pb.IntfDetails{
			Vfid:   3,
			Pfid:   0,
			NumPfs: 1,
		},
		MacAddr: "00:00:00:01:01:02",
		IpamAllocatedbyHost: &pb.IPAMallocatedbyHost{
			IpamAllocatedbyHost: false,
		},
		IpInfo: &pb.IpDetails{
			Version: "4",
			Address: "10.0.0.1",
			Gateway: "10.0.0.0",
			Dns:     "",
		},
		OffloadFlags: &pb.OffloadFlags{
			CheckSumOffload: false,
		},
	}
	return sampleAddReq
}

func getSampleDelReq() *pb.CNIDelRequest {
	sampleDelReq := &pb.CNIDelRequest{
		IntfInfo: &pb.IntfDetails{
			Vfid:   3,
			Pfid:   0,
			NumPfs: 1,
		},
		PodInfo: &pb.PodDetails{
			// Should not fill PodNamespace and PodName w/o actual pod
			ContainerId: "cnitool-77383ca0a0715733ca6a",
		},
	}
	return sampleDelReq
}

func getSampleCheckReq() *pb.CNICheckRequest {
	sampleCheckReq := &pb.CNICheckRequest{
		IntfInfo: &pb.IntfDetails{
			Vfid:   0,
			Pfid:   0,
			NumPfs: 1,
		},
		PodInfo: &pb.PodDetails{
			// Should not fill PodNamespace and PodName w/o actual pod
			ContainerId: "",
		},
	}
	return sampleCheckReq
}
