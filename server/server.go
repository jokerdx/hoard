package server

import (
	"fmt"
	"net"
	"strings"

	"code.monax.io/platform/hoard/core"
	"code.monax.io/platform/hoard/core/storage"
	"github.com/go-kit/kit/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	listenURL  string
	grpcServer *grpc.Server
	logger     log.Logger
}

func New(listenURL string, logger log.Logger) *server {
	return &server{
		listenURL: listenURL,
		logger:    logger,
	}
}

func (serv *server) Serve() error {
	netProtocol, localAddress, err := SplitListenURL(serv.listenURL)
	if err != nil {
		return fmt.Errorf("Failed to split listen URL '%s': %v", serv.listenURL, err)
	}
	listener, err := net.Listen(netProtocol, localAddress)
	if err != nil {
		return fmt.Errorf("Failed to create listener: %v", err)
	}
	serv.grpcServer = grpc.NewServer()
	hoardServer := core.NewHoardServer(core.NewHoard(storage.NewMemoryStore(),
		serv.logger))

	core.RegisterCleartextServer(serv.grpcServer, hoardServer)
	core.RegisterEncryptionServer(serv.grpcServer, hoardServer)
	core.RegisterStorageServer(serv.grpcServer, hoardServer)
	// Register reflection service on gRPC server.
	reflection.Register(serv.grpcServer)
	err = serv.grpcServer.Serve(listener)
	if err != nil {
		return fmt.Errorf("Failed to start GRPC Server: %v", err)
	}
	return nil
}

func (serv *server) Stop() {
	serv.grpcServer.Stop()
}

func SplitListenURL(listenOn string) (string, string, error) {
	// net.Listen does not want a parsed url.URL so it seems to make more sense
	// just to do a dumb split here to support the various networks
	listenParts := strings.Split(listenOn, "://")
	if len(listenParts) != 2 {
		return "", "", fmt.Errorf("Expected a Go net.Listen URL of the form "+
			"'<net>://<laddr>', but got: '%s'", listenOn)
	}
	if listenParts[0] == "" {
		return "", "", fmt.Errorf("Expected the URL scheme to be present, "+
			"but got '%s'", listenOn)
	}
	if listenParts[1] == "" {
		return "", "", fmt.Errorf("Expected the URL host to be present, "+
			"but got '%s'", listenOn)
	}
	return listenParts[0], listenParts[1], nil
}