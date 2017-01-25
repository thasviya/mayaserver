package server

import (
	"io"
	"log"
	"sync"
)

// MayaServer is a long running stateless daemon that runs
// at openebs maya master(s)
type MayaServer struct {
	config    *MayaConfig
	logger    *log.Logger
	logOutput io.Writer

	shutdown     bool
	shutdownCh   chan struct{}
	shutdownLock sync.Mutex
}

// NewMayaServer is used to create a new maya server
// with the given configuration
func NewMayaServer(config *MayaConfig, logOutput io.Writer) (*MayaServer, error) {
	ms := &MayaServer{
		config:     config,
		logger:     log.New(logOutput, "", log.LstdFlags|log.Lmicroseconds),
		logOutput:  logOutput,
		shutdownCh: make(chan struct{}),
	}

	return ms, nil
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
