package conn

import (
	"time"

	"github.com/SKF/go-utility/log"
	grpc_retry "github.com/lonnblad/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type Client interface {
	SetConnection(*grpc.ClientConn)
	Close()
}

func Connect(client Client, host, port string, opts ...grpc.DialOption) (err error) {
	retryOpts := []grpc_retry.CallOption{
		grpc_retry.WithBackoff(
			grpc_retry.BackoffLinear(time.Second),
		),
		grpc_retry.WithMax(10),
	}

	kp := keepalive.ClientParameters{
		Time:                1 * time.Second,
		Timeout:             100 * time.Millisecond,
		PermitWithoutStream: false,
	}

	opts = append(
		opts,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithStreamInterceptor(grpc_retry.StreamClientInterceptor(retryOpts...)),
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(retryOpts...)),
		grpc.WithKeepaliveParams(kp),
	)

	conn, err := grpc.Dial(host+":"+port, opts...)
	if err != nil {
		log.
			WithField("error", err).
			WithField("host", host).
			WithField("port", port).
			Fatal("failed to dial server")
		return
	}

	log.Info("Dialed successfully")
	client.SetConnection(conn)

	return
}
