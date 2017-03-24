package v1

type ContainerNetworkingLbl string

const (
	CNTypeLbl            ContainerNetworkingLbl = "cn.openebs.io/type"
	CNNetworkCIDRAddrLbl ContainerNetworkingLbl = "cn.openebs.io/network-cidr-addr"
	CNSubnetLbl          ContainerNetworkingLbl = "cn.openebs.io/subnet"
	CNInterfaceLbl       ContainerNetworkingLbl = "cn.openebs.io/interface"
)

type ContainerStorageLbl string

const (
	CSPersistenceLocationLbl ContainerStorageLbl = "cs.openebs.io/persistence-location"
)

type RequestsLbl string

const (
	RegionLbl     RequestsLbl = "requests.openebs.io/region"
	DatacenterLbl RequestsLbl = "requests.openebs.io/dc"
)

const (
	VolumePluginNamePrefix string = "name.plugin.volume.openebs.io/"
)

const (
	DefaultOrchestratorConfigPath string = "/etc/mayaserver/orchprovider/"
)
