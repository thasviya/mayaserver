package server

import (
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/openebs/mayaserver/lib/config"
	"github.com/openebs/mayaserver/lib/orchprovider"
	"github.com/openebs/mayaserver/lib/orchprovider/nomad"
	"github.com/openebs/mayaserver/lib/volume"
	"github.com/openebs/mayaserver/lib/volume/jiva"
)

// MayaServer is a long running stateless daemon that runs
// at openebs maya master(s)
type MayaServer struct {
	config       *config.MayaConfig
	pluginsMutex sync.Mutex
	volPlugins   map[string]volume.VolumeInterface
	logger       *log.Logger
	logOutput    io.Writer

	shutdown     bool
	shutdownCh   chan struct{}
	shutdownLock sync.Mutex
}

// NewMayaServer is used to create a new maya server
// with the given configuration
func NewMayaServer(config *config.MayaConfig, logOutput io.Writer) (*MayaServer, error) {
	ms := &MayaServer{
		config:     config,
		volPlugins: make(map[string]volume.VolumeInterface),
		logger:     log.New(logOutput, "", log.LstdFlags|log.Lmicroseconds),
		logOutput:  logOutput,
		shutdownCh: make(chan struct{}),
	}

	err := ms.BootstrapPlugins()
	if err != nil {
		return nil, err
	}

	return ms, nil
}

// TODO
// Create a Bootstrap interface that facilitates initialization
// Create another Bootstraped interface that provides the initialized instances
// Perhaps at lib/bootstrap
// MayaServer struct will make use of above interfaces & hence specialized
// structs that cater to bootstraping & bootstraped features.
//
// NOTE:
//    The current implementation is tightly coupled & cannot be unit tested.
func (ms *MayaServer) BootstrapPlugins() error {

	// TODO
	// Use MayaConfig
	// Fetch the names of volume plugins to be initialized
	// Iterate over the volumes:
	//  0. Fetch the config file location of orchestrator
	//  1. Initialize volume plugin's orchestrator
	//  2. Build an aspect that points to above orchestrator
	//  3. Fetch the config file location of volume plugin
	//  4. Initialize volume plugin

	ms.pluginsMutex.Lock()
	defer ms.pluginsMutex.Unlock()

	orchestrator, err := orchprovider.InitOrchProvider(nomad.NomadOrchProviderName, "")
	if err != nil {
		return err
	}

	jivaAspect := &jiva.JivaStorNomadAspect{
		Nomad: orchestrator,
	}

	jivaStor, err := volume.InitVolumePlugin(jiva.JivaStorPluginName, "", jivaAspect)
	if err != nil {
		return err
	}

	ms.volPlugins[jiva.JivaStorPluginName] = jivaStor
	return nil
}

// GetVolumePlugin is an accessor that fetches a volume.VolumeInterface instance
// The volume.VolumeInterface should have been initialized earlier.
func (ms *MayaServer) GetVolumePlugin(name string) (volume.VolumeInterface, error) {
	ms.pluginsMutex.Lock()
	defer ms.pluginsMutex.Unlock()

	storage, found := ms.volPlugins[name]
	if !found {
		return nil, fmt.Errorf("Volume plugin '%s' not found", name)
	}

	return storage, nil
}

// Shutdown is used to terminate MayaServer.
func (ms *MayaServer) Shutdown() error {
	ms.shutdownLock.Lock()
	defer ms.shutdownLock.Unlock()

	ms.logger.Println("[INFO] mayaserver: requesting shutdown")

	if ms.shutdown {
		return nil
	}

	ms.logger.Println("[INFO] mayaserver: shutdown complete")
	ms.shutdown = true
	close(ms.shutdownCh)
	return nil
}

// Leave is used gracefully exit.
func (ms *MayaServer) Leave() error {
	// Nothing as of now
	return nil
}
