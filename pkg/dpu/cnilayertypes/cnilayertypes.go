package cnilayertypes

// TODO: Add all the field that is needed for the offload agent only required fields are added
type OffloadCNIAddRequest struct {
	PodInfo             *PodDetails
	IntfInfo            *IntfDetails
	MacAddr             string
	IpamAllocatedbyHost bool
	IpInfo              *IpDetails
	OffloadFlags        *OffloadFlags
	Vlan                string
}
type OffloadCNIAddResponse struct {
	IpInfo *IpDetails
	Status string
	Vlan   string
}

type OffloadCNIDelRequest struct {
	IntfInfo *IntfDetails
	PodInfo  *PodDetails
}

/*
OffloadCNIDelResponse Response type for Del
*/
type OffloadCNIDelResponse struct {
	IpInfo *IpDetails
	Status string
}

/*
OffloadCNICheckResponse Response type for Check
*/
type OffloadCNICheckResponse struct {
	IpInfo *IpDetails
	Status string
}

/*
OffloadCNICheckRequest Fields are added based on what is needed to cover the scenarios
defined for CniAddRequest type
*/
type OffloadCNICheckRequest struct {
	IntfInfo *IntfDetails
	PodInfo  *PodDetails
	IpInfo   *IpDetails
}

type PodDetails struct {
	PodNamespace string
	PodName      string
	ContainerId  string
}

type IpDetails struct {
	Version string
	Address string
	Gateway string
	Dns     string
}
type IntfDetails struct {
	Name   string
	Vfid   uint32
	Pfid   uint32
	NumPfs uint32
	Vlan   uint32
}
type OffloadFlags struct {
	CheckSumOffload bool
}
