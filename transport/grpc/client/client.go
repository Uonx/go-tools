package client

import (
	"context"
	"net"

	"google.golang.org/grpc"
)

func GrpcDial(network, address string) (*grpc.ClientConn, error) {
	if network != "unix" {
		return grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	}
	unixAddress, err := net.ResolveUnixAddr(network, address)
	if err != nil {
		return nil, err
	}
	return grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
		return net.DialUnix(network, nil, unixAddress)
	}))
}
