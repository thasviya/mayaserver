package nomad

import (
	"github.com/openebs/mayaserver/lib/api/v1"
)

const (

	// This is the default Nomad orchestrator's region. This is the assumed
	// region where Nomad will be running/deployed.
	DefaultNomadRegionName string = "global"

	// This is the default Nomad orchestrator's datacenter. This is the assumed
	// datacenter where Nomad will be operational.
	DefaultNomadDCName string = "dc1"

	// This is the default Nomad orchestrator's container network type.
	DefaultNomadCNType string = "host"

	// This is the default Nomad orchestrator's container network CIDR
	DefaultNomadCNNetworkCIDR string = "172.28.128.1/24"

	// This is the default Nomad orchestrator's container network interface
	DefaultNomadCNInterface string = "enp0s8"

	// This is the default Nomad orchestrator config file. This typically
	// points to Nomad config when Nomad is running in default region.
	DefaultNomadConfigFile string = v1.DefaultOrchestratorConfigPath + "nomad_" + DefaultNomadRegionName + ".INI"

	// This is the default variant of Nomad orchestrator. This typically
	// points to Nomad running in default region.
	DefaultNomadPluginName = "nomad_" + DefaultNomadRegionName
)
