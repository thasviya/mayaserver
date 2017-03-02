package orchprovider

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/golang/glog"
)

// Factory is a function that returns a orchprovider.OrchestratorInterface.
// The config parameter provides an io.Reader handler to the factory in
// order to load specific configurations. If no configuration is provided
// the parameter is nil.
type Factory func(config io.Reader) (OrchestratorInterface, error)

// All registered orchestration providers.
var (
	providersMutex sync.Mutex
	providers      = make(map[string]Factory)
)

// RegisterOrchProvider registers a orchprovider.Factory by name.
//
// This is expected to happen during binary startup.
//
// How ?
//    Each implementation of orchestration provider need to call
// RegisterOrchProvider inside their init() function.
func RegisterOrchProvider(name string, factory Factory) {
	providersMutex.Lock()
	defer providersMutex.Unlock()
	if _, found := providers[name]; found {
		glog.Fatalf("Orchestration provider %q was registered twice", name)
	}
	glog.V(1).Infof("Registered orchestration provider %q", name)
	providers[name] = factory
}

// IsOrchProvider returns true if name corresponds to an already
// registered orchestration provider.
func IsOrchProvider(name string) bool {
	providersMutex.Lock()
	defer providersMutex.Unlock()

	_, found := providers[name]
	return found
}

// OrchProviders returns the name of all registered orchestration
// providers in a string slice
func OrchProviders() []string {
	names := []string{}
	providersMutex.Lock()
	defer providersMutex.Unlock()

	for name := range providers {
		names = append(names, name)
	}
	return names
}

// GetOrchProvider creates an instance of the named orchestration provider,
// or nil if the name is unknown.  The error return is only used if the named
// provider was known but failed to initialize. The config parameter specifies
// the io.Reader handler of the configuration file for the orchestration
// provider, or nil for no configuation.
func GetOrchProvider(name string, config io.Reader) (OrchestratorInterface, error) {
	providersMutex.Lock()
	defer providersMutex.Unlock()

	factory, found := providers[name]
	if !found {
		return nil, nil
	}
	return factory(config)
}

// TODO
// Who calls this ?
// This will currently be triggered while starting the binary as a http service ?
//
// InitOrchProvider creates an instance of the named orchestration provider.
func InitOrchProvider(name string, configFilePath string) (OrchestratorInterface, error) {
	var orchestrator OrchestratorInterface
	var err error

	if name == "" {
		glog.Info("No orchestration provider specified.")
		return nil, nil
	}

	if configFilePath != "" {
		var config *os.File
		config, err = os.Open(configFilePath)
		if err != nil {
			glog.Fatalf("Couldn't open orchestration provider configuration %s: %#v",
				configFilePath, err)
		}

		defer config.Close()
		orchestrator, err = GetOrchProvider(name, config)
	} else {
		// Pass explicit nil so plugins can actually check for nil. See
		// "Why is my nil error value not equal to nil?" in golang.org/doc/faq.
		orchestrator, err = GetOrchProvider(name, nil)
	}

	if err != nil {
		return nil, fmt.Errorf("could not init orchestration provider %q: %v", name, err)
	}
	if orchestrator == nil {
		return nil, fmt.Errorf("unknown orchestration provider %q", name)
	}

	return orchestrator, nil
}
