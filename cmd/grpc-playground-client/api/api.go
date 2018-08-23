package api

import (
	"context"
	"io"
	"time"

	"github.com/SKF/go-utility/log"
	"github.com/lonnblad/grpc-playground/playapi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Client struct {
	api  playapi.PlaygroundClient
	conn *grpc.ClientConn
}

func (c *Client) SetConnection(conn *grpc.ClientConn) {
	c.api = playapi.NewPlaygroundClient(conn)
	c.conn = conn
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) DeepPing() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err = c.api.DeepPing(ctx, &playapi.Void{})
	return
}

func (c *Client) InvokePanic() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err = c.api.InvokePanic(ctx, &playapi.Void{})
	return
}

func (c *Client) GetReceiveStream(ctx context.Context, dc chan<- string) error {
	defer close(dc)
	for {
		stream, err := c.api.GetReceiveStream(ctx, &playapi.Void{})
		if err != nil {
			log.WithError(err).Error("Get receive stream failed")
			return err
		}

		for {
			strObj, err := stream.Recv()
			if err == nil {
				dc <- strObj.Value
				continue
			}
			if err == io.EOF {
				return nil
			}

			status := status.Convert(err)
			if status.Code() == codes.Unavailable ||
				status.Code() == codes.ResourceExhausted {
				log.WithError(err).Error("Will try to reconnect")
				break
			}

			return err
		}
	}
}

func (c *Client) GetSendStream(dc <-chan string) (err error) {
	ctx := context.Background()
	stream, err := c.api.GetSendStream(ctx)
	if err != nil {
		return
	}

	for {
		msg := <-dc
		strObj := &playapi.StringObject{Value: msg}
		err = stream.Send(strObj)
		if err != nil {
			return
		}
	}
}
