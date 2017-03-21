package nomad

import (
	"fmt"
	"io"
	"os"

	"github.com/golang/glog"
	"github.com/hashicorp/nomad/api"
	"github.com/openebs/mayaserver/lib/api/v1"
	v1nomad "github.com/openebs/mayaserver/lib/api/v1/nomad"
	gcfg "gopkg.in/gcfg.v1"
)

const (
	// Names of environment variables used to supply the coordinates
	// of a Nomad deployment
	EnvNomadAddress = "NOMAD_ADDR"
	EnvNomadRegion  = "NOMAD_REGION"
)

// NomadConfig provides the settings that has the coordinates of a
// Nomad server or a Nomad cluster deployment.
//
// In addition, it provides some of the networking options
// that can be consumed by a container spawned via Nomad orchestrator.
//
// A NomadConfig file has .INI extension.
// Below is a sample:
//
// [datacenter "dc1"]
// address = http://10.0.0.1:4646
//
// [datacenter "dc2"]
// address = http://20.0.0.2:4646
// cn-type = host
// cn-network-cidr = 172.28.128.1/24
// cn-interface = enp0s8
//
// NOTE:
//    This is as per gcfg lib's conventions
type NomadConfig struct {
	Datacenter map[string]*struct {
		// Address of Nomad cluster
		Address string

		// Whether it is a host based networking or something else
		// Required while placing a container inside Nomad
		CNType string `gcfg:"cn-type"`

		// The Network address in CIDR notation. Available IP addresses will
		// be considered from this network range.
		CNNetworkCIDR string `gcfg:"cn-network-cidr"`

		// Networking interface that is available in the Nomad cluster
		CNInterface string `gcfg:"cn-interface"`
	}
}

// NomadUtilInterface is an abstraction over Hashicorp's
// Nomad properties & communication utilities.
type NomadUtilInterface interface {

	// Name of nomad utility
	Name() string

	// This is a builder for NomadClients interface. Will return
	// false if not supported.
	NomadClients() (NomadClients, bool)

	// This is a builder for NomadNetworks interface. Will return
	// false if not supported.
	NomadNetworks() (NomadNetworks, bool)
}

// NomadClients is an abstraction over various connection modes (http, rpc)
// to Nomad. Http client is currently supported.
//
// NOTE:
//    This abstraction makes use of Nomad's api package. Nomad's api
// package provides:
//
// 1. http client abstraction &
// 2. structures that can send http requests to Nomad's APIs.
type NomadClients interface {
	// Http returns the http client that is capable to communicate
	// with Nomad
	Http() (*api.Client, error)
}

// NomadNetworks is a blueprint to expose various networking options
// available in a Nomad cluster.
type NomadNetworks interface {
	// CN exposes various networking values that is supported at a
	// particular datacenter where Nomad is running
	CN(dc string) (map[v1.ContainerNetworkingLbl]string, error)
}

// nomadUtil is the concrete implementation for
//
// 1. nomad.NomadClients interface
// 2. nomad.NomadNetworks interface
type nomadUtil struct {

	// The region to send API requests to
	// TODO
	// This will be set during this instance creation time
	region string

	// The datacenter to send API requests to
	// TODO
	// This will be set during this instance creation time
	datacenter string

	// Nomad server / cluster coordinates
	// This will be set based on the region
	nomadConf *NomadConfig

	caCert     string
	caPath     string
	clientCert string
	clientKey  string
	insecure   bool
}

// newNomadUtil provides a new instance of nomadUtil
//
// TODO
// region may be passed as an argument
// & hence NomadConfig should be instantiated based on the region
// at this place
func newNomadUtil(nConfig *NomadConfig) (*nomadUtil, error) {
	return &nomadUtil{
		nomadConf: nConfig,
	}, nil
}

// This is a plain nomad utility & hence the name
func (m *nomadUtil) Name() string {
	return "nomadutil"
}

// nomadUtil implements NomadClients interface. Hence it returns
// self
func (m *nomadUtil) NomadClients() (NomadClients, bool) {
	return m, true
}

// nomadUtil implements NomadNetworks interface. Hence it returns
// self
func (m *nomadUtil) NomadNetworks() (NomadNetworks, bool) {
	return m, true
}

// CN provides the container networking data in key-value pairs.
// These networking data are supposed to be available in the target Nomad
// cluster. These pairs are provided based on datacenter.
func (m *nomadUtil) CN(dcName string) (map[v1.ContainerNetworkingLbl]string, error) {

	err := m.validateConf(dcName)
	if err != nil {
		return nil, err
	}

	// build the cn map
	cn := map[v1.ContainerNetworkingLbl]string{
		v1.CNTypeLbl:            m.getCNType(dcName),
		v1.CNNetworkCIDRAddrLbl: m.getCNNetworkCIDR(dcName),
		v1.CNInterfaceLbl:       m.getCNInterface(dcName),
	}

	return cn, nil
}

func (m *nomadUtil) validateConf(dcName string) error {

	if dcName == "" {
		return fmt.Errorf("DC name is empty")
	}

	if m.nomadConf == nil {
		return fmt.Errorf("Nil nomad config provided")
	}

	if m.nomadConf.Datacenter == nil {
		return fmt.Errorf("DC not available in nomad config")
	}

	if m.nomadConf.Datacenter[dcName] == nil {
		return fmt.Errorf("No details available for dc '%s'", dcName)
	}

	return nil
}

// getCNInterface extracts the network type from conf or returns the default value
func (m *nomadUtil) getCNType(dcName string) string {

	if m.nomadConf.Datacenter[dcName].CNType == "" {
		return v1nomad.DefaultNomadCNType
	}

	return m.nomadConf.Datacenter[dcName].CNType
}

// getCNInterface extracts the network CIDR from conf or returns the default value
func (m *nomadUtil) getCNNetworkCIDR(dcName string) string {

	if m.nomadConf.Datacenter[dcName].CNNetworkCIDR == "" {
		return v1nomad.DefaultNomadCNNetworkCIDR
	}

	return m.nomadConf.Datacenter[dcName].CNNetworkCIDR
}

// getCNInterface extracts the interface from conf or returns the default value
func (m *nomadUtil) getCNInterface(dcName string) string {

	if m.nomadConf.Datacenter[dcName].CNInterface == "" {
		return v1nomad.DefaultNomadCNInterface
	}

	return m.nomadConf.Datacenter[dcName].CNInterface
}

// Client is used to initialize and return a new API client capable
// of calling Nomad APIs.
// TODO
// datacenter may be passed as a parameter
func (m *nomadUtil) Http() (*api.Client, error) {
	// Nomad API client config
	apiCConf := api.DefaultConfig()

	// Set from environment variable
	val, found := os.LookupEnv(EnvNomadAddress)

	if !found {
		glog.V(2).Infof("Env variable '%s' is not set", EnvNomadAddress)
	}

	if val != "" {
		glog.V(2).Infof("Nomad address is set to '%s' via env var", val)
		apiCConf.Address = val
	}

	// Override from conf structure
	if m.nomadConf != nil && m.nomadConf.Datacenter != nil {
		// TODO
		// Derive the datacenter at runtime
		// Remove the region & datacenter properties from Mayaconfig
		glog.V(2).Infof("Nomad address is set to: '%s' via conf", m.nomadConf.Datacenter["dc1"].Address)
		apiCConf.Address = m.nomadConf.Datacenter["dc1"].Address
	}

	if apiCConf.Address == "" {
		return nil, fmt.Errorf("Nomad address is not set")
	}

	glog.V(2).Infof("Nomad will be reached at: '%s'", apiCConf.Address)

	if v := os.Getenv(EnvNomadRegion); v != "" {
		apiCConf.Region = v
	}

	if m.region != "" {
		apiCConf.Region = m.region
	}

	// If we need custom TLS configuration, then set it
	if m.caCert != "" || m.caPath != "" || m.clientCert != "" || m.clientKey != "" || m.insecure {
		t := &api.TLSConfig{
			CACert:     m.caCert,
			CAPath:     m.caPath,
			ClientCert: m.clientCert,
			ClientKey:  m.clientKey,
			Insecure:   m.insecure,
		}
		apiCConf.TLSConfig = t
	}

	// This has the http address & authentication details
	// required to invoke Nomad APIs
	return api.NewClient(apiCConf)
}

// readNomadConfig reads an instance of NomadConfig from config reader.
func readNomadConfig(config io.Reader) (*NomadConfig, error) {
	var nCfg NomadConfig
	var err error

	if config != nil {
		err = gcfg.ReadInto(&nCfg, config)
		if err != nil {
			return nil, err
		}
	}

	// TODO
	// validations w.r.t config

	return &nCfg, nil
}
