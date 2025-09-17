package configs

import (
	"context"
	authv1 "music-player/api/proto/auth/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClients struct {
	AuthClient authv1.AuthServiceClient
	authConn   *grpc.ClientConn
}

func NewGRPCClients(ctx context.Context, authEndpoint string) (*GRPCClients, error) {
	authConn, err := grpc.DialContext(ctx, authEndpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, err
	}

	return &GRPCClients{
		AuthClient: authv1.NewAuthServiceClient(authConn),
		authConn:   authConn,
	}, nil
}

// Close closes all gRPC connections
func (c *GRPCClients) Close() error {
	if c.authConn != nil {
		return c.authConn.Close()
	}
	return nil
}
