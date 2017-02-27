package nomad

const (
	// Names of environment variables used to supply the coordinates
	// of a Nomad deployment
	EnvNomadAddress = "NOMAD_ADDR"
	EnvNomadRegion  = "NOMAD_REGION"
)

// NomadConfig provides the settings that has the coordinates of a 
// Nomad server or a Nomad cluster deployment.
type NomadConfig struct {

  Address string
}

// nomadClientUtil is the concrete implementation for nomad.NomadClient
// interface
type nomadClientUtil struct {

	// The region to send API requests
	region string

	caCert     string
	caPath     string
	clientCert string
	clientKey  string
	insecure   bool
}

// Client is used to initialize and return a new API client capable
// of calling Nomad. It uses env vars.
//
// TODO
// Make use of nomad.NomadConfig also
func (m *nomadClientUtil) HttpClient() (*api.Client, error) {
  return nil, nil
}
