package configs

import (
	"context"
	authv1 "music-player/api/proto/auth/v1"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

type GRPCClients struct {
	AuthClient authv1.AuthServiceClient
	authConn   *grpc.ClientConn
}

func NewGRPCClients(ctx context.Context, authEndpoint string) (*GRPCClients, error) {
	const svcCfg = `{
	  "loadBalancingPolicy": "round_robin",
	  "methodConfig": [{
	    "name": [{"service": "auth.v1.AuthService"}],
	    "timeout": "3s"
	  }]
	}`

	authConn, err := grpc.DialContext(
		ctx,
		authEndpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                30 * time.Second,
			Timeout:             5 * time.Second,
			PermitWithoutStream: true,
		}),
		grpc.WithConnectParams(grpc.ConnectParams{
			MinConnectTimeout: 5 * time.Second,
			Backoff: backoff.Config{
				BaseDelay:  100 * time.Millisecond,
				Multiplier: 1.6,
				Jitter:     0.2,
				MaxDelay:   2 * time.Second,
			},
		}),
		grpc.WithDefaultServiceConfig(svcCfg),
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
