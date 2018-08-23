package server

import (
	"time"

	"github.com/SKF/go-utility/log"
	"github.com/lonnblad/grpc-playground/playapi"
	"golang.org/x/net/context"
)

type grpcServer struct{}

func Create() playapi.PlaygroundServer {
	return &grpcServer{}
}

func (s *grpcServer) DeepPing(ctx context.Context, _ *playapi.Void) (output *playapi.StringObject, err error) {
	output = &playapi.StringObject{Value: "Pong"}
	return
}

func (s *grpcServer) InvokePanic(_ context.Context, _ *playapi.Void) (*playapi.Void, error) {
	panic("InvokePanic")
}

func (s *grpcServer) GetReceiveStream(_ *playapi.Void, stream playapi.Playground_GetReceiveStreamServer) (err error) {
	for {
		select {
		case <-stream.Context().Done():
			return
		default:
			stream.Send(&playapi.StringObject{Value: "Hey"})
			time.Sleep(time.Second)
		}
	}
}

func (s *grpcServer) GetSendStream(stream playapi.Playground_GetSendStreamServer) (err error) {
	for {
		select {
		case <-stream.Context().Done():
			stream.SendAndClose(&playapi.Void{})
			return nil
		default:
			msg, err := stream.Recv()
			if err != nil {
				stream.SendAndClose(&playapi.Void{})
				return err
			}
			log.WithField("msg", msg.Value).
				Info("Received Message")
		}
	}
}
