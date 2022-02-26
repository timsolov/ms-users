package grpc_client

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// BlockConnect starts client blocking connection to gRPC server.
// The gRPC server should be run before.
// Transport: insecure.
func BlockConnect(ctx context.Context, dialAddr string) (conn *grpc.ClientConn, err error) {
	// Create a client connection to the gRPC Server we just started.
	// This is where the gRPC-Gateway proxies the requests.
	return grpc.DialContext(
		ctx,
		"dns:///"+dialAddr,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
}
