package configs

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// GRPCServer wraps the gRPC server
type GRPCServer struct {
	server   *grpc.Server
	listener net.Listener
}

// NewGRPCServer creates a new gRPC server instance
func NewGRPCServer(port string) (*GRPCServer, error) {
	addr := fmt.Sprintf(":%s", port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on %s: %w", addr, err)
	}

	opts := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(4 * 1024 * 1024),
		grpc.MaxSendMsgSize(4 * 1024 * 1024),
	}

	server := grpc.NewServer(opts...)

	reflection.Register(server)

	return &GRPCServer{
		server:   server,
		listener: lis,
	}, nil
}

// GetServer returns the underlying gRPC server
func (s *GRPCServer) GetServer() *grpc.Server {
	return s.server
}

// Start starts the gRPC server
func (s *GRPCServer) Start() error {
	return s.server.Serve(s.listener)
}

// Stop gracefully stops the gRPC server
func (s *GRPCServer) Stop() {
	s.server.GracefulStop()
}

// GetAddr returns the server address
func (s *GRPCServer) GetAddr() net.Addr {
	return s.listener.Addr()
}
