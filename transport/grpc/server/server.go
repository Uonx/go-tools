package server

import (
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

func StartGrpcServer(network, address string, registrys func(server *grpc.Server) func(), interceptors ...grpc.UnaryServerInterceptor) {
	var kaep = keepalive.EnforcementPolicy{
		MinTime:             5 * time.Second, // If a client pings more than once every 5 seconds, terminate the connection
		PermitWithoutStream: true,            // Allow pings even when there are no active streams
	}

	var kasp = keepalive.ServerParameters{
		Time:    5 * time.Second, // Ping the client if it is idle for 5 seconds to ensure the connection is still active
		Timeout: 1 * time.Second, // Wait 1 second for the ping ack before assuming the connection is dead
	}
	server := grpc.NewServer(
		grpc.MaxConcurrentStreams(8192),
		grpc.InitialConnWindowSize(1*1024*1024*1024),
		grpc.KeepaliveEnforcementPolicy(kaep),
		grpc.KeepaliveParams(kasp),
		grpc.ChainUnaryInterceptor(interceptors...),
	)

	var (
		err    error
		listen net.Listener
	)
	if network == "unix" {
		serverAddress, err := net.ResolveUnixAddr(network, address)
		if err != nil {
			panic(err)
		}
		if listen, err = net.ListenUnix("unix", serverAddress); err != nil {
			panic(err)
		}
	} else {
		listen, err = net.Listen(network, address)
	}
	if err != nil {
		panic(err)
	}
	registrys(server)()
	log.Printf("server listening at %s \n", listen.Addr().String())
	if err := server.Serve(listen); err != nil {
		panic(err)
	}
}
