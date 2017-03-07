package nomad

import (
	"fmt"
	"io"
	"os"

	"github.com/golang/glog"
	"github.com/hashicorp/nomad/api"
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
// A NomadConfig file has .INI extension.
// Below is a sample:
//
// [datacenter "dc1"]
// address = http://10.0.0.1:4646
//
// [datacenter "dc2"]
// address = http://20.0.0.2:4646
//
// NOTE:
//    This is as per gcfg lib's conventions
type NomadConfig struct {
	Datacenter map[string]*struct {
		Address string
	}
}

// NomadClient is an abstraction over various connection modes (http, rpc)
// to Nomad. Http client is currently supported.
//
// NOTE:
//    This abstraction makes use of Nomad's api package. Nomad's api
// package provides:
//
// 1. http client abstraction &
// 2. structures that can send http requests to Nomad's APIs.
type NomadClient interface {
	Http() (*api.Client, error)
}

// nomadClientUtil is the concrete implementation for nomad.NomadClient
// interface.
type nomadClientUtil struct {

	// The region to send API requests
	region string

	// Nomad server / cluster coordinates
	nomadConf *NomadConfig

	caCert     string
	caPath     string
	clientCert string
	clientKey  string
	insecure   bool
}

// newNomadClientUtil provides a new instance of nomadClientUtil
func newNomadClientUtil(nConfig *NomadConfig) (*nomadClientUtil, error) {
	return &nomadClientUtil{
		nomadConf: nConfig,
	}, nil
}

// Client is used to initialize and return a new API client capable
// of calling Nomad APIs. It uses env vars.
func (m *nomadClientUtil) Http() (*api.Client, error) {
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
