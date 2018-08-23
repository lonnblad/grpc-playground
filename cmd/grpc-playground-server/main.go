package main

import (
	"net"
	"runtime"

	"github.com/SKF/go-utility/grpc-interceptor/requestid"
	"github.com/SKF/go-utility/log"
	"github.com/lonnblad/grpc-playground/cmd/grpc-playground-server/server"
	"github.com/lonnblad/grpc-playground/playapi"

	grpc_middleware "github.com/lonnblad/go-grpc-middleware"
	grpc_recovery "github.com/lonnblad/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/lonnblad/go-grpc-middleware/tags"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
)

func listen(grpcServer playapi.PlaygroundServer, port string, opts ...grpc.ServerOption) {
	server := grpc.NewServer(opts...)

	listenPort := ":" + port
	lis, err := net.Listen("tcp", listenPort)
	if err != nil {
		log.WithError(err).Fatal("Failed to listen")
	}

	playapi.RegisterPlaygroundServer(server, grpcServer)
	reflection.Register(server)

	log.Info("grpc server is started")
	if err := server.Serve(lis); err != nil {
		log.WithError(err).Fatal("Failed to serve")
	}
}

func getRecoveryOption() grpc_recovery.Option {
	recoveryFunc := func(p interface{}) error {
		buf := make([]byte, 10000)
		written := runtime.Stack(buf, false)
		log.
			WithField("panic", p).
			WithField("stackTrace", string(buf[:written])).
			Error("Internal Server Error")
		return grpc.Errorf(codes.Internal, "Internal Server Error")
	}
	return grpc_recovery.WithRecoveryHandler(recoveryFunc)
}

func createServerOptions(logName string) []grpc.ServerOption {
	recoveryOption := getRecoveryOption()

	serverOpts := []grpc.ServerOption{
		grpc_middleware.WithUnaryServerChain(
			grpc_ctxtags.UnaryServerInterceptor(),
			requestid.UnaryServerInterceptor(logName),
			grpc_recovery.UnaryServerInterceptor(recoveryOption),
		),
		grpc_middleware.WithStreamServerChain(
			grpc_ctxtags.StreamServerInterceptor(),
			requestid.StreamServerInterceptor(logName),
			grpc_recovery.StreamServerInterceptor(recoveryOption),
		),
	}

	serverOpts = append(serverOpts,
		grpc.StatsHandler(statsHandler{}),
	)

	return serverOpts
}

func main() {
	port := "50051"
	logName := "Bepa"

	grpcServer := server.Create()
	listen(grpcServer, port, createServerOptions(logName)...)
}
