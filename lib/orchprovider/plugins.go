package orchprovider

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/golang/glog"
)

// Factory is a function that returns a orchprovider.Interface.
// The config parameter provides an io.Reader handler to the factory in
// order to load specific configurations. If no configuration is provided
// the parameter is nil.
type Factory func(config io.Reader) (Interface, error)

// All registered orchestration providers.
var (
	providersMutex sync.Mutex
	providers      = make(map[string]Factory)
)

// RegisterOrchProvider registers a orchprovider.Factory by name.
// This is expected to happen during binary startup.
func RegisterOrchProvider(name string, orchestrator Factory) {
	providersMutex.Lock()
	defer providersMutex.Unlock()
	if _, found := providers[name]; found {
		glog.Fatalf("Orchestration provider %q was registered twice", name)
	}
	glog.V(1).Infof("Registered orchestration provider %q", name)
	providers[name] = orchestrator
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
func GetOrchProvider(name string, config io.Reader) (Interface, error) {
	providersMutex.Lock()
	defer providersMutex.Unlock()
	f, found := providers[name]
	if !found {
		return nil, nil
	}
	return f(config)
}

// InitOrchProvider creates an instance of the named orchestration provider.
func InitOrchProvider(name string, configFilePath string) (Interface, error) {
	var orchestrator Interface
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
		orchestration, err = GetOrchProvider(name, config)
	} else {
		// Pass explicit nil so plugins can actually check for nil. See
		// "Why is my nil error value not equal to nil?" in golang.org/doc/faq.
		orchestration, err = GetOrchProvider(name, nil)
	}

	if err != nil {
		return nil, fmt.Errorf("could not init orchestration provider %q: %v", name, err)
	}
	if orchestration == nil {
		return nil, fmt.Errorf("unknown orchestration provider %q", name)
	}

	return orchestration, nil
}
