package main

import (
	"context"

	"github.com/SKF/go-utility/log"
	"google.golang.org/grpc/stats"
)

type statsHandler struct{}

func (h statsHandler) TagRPC(ctx context.Context, info *stats.RPCTagInfo) context.Context {
	if info != nil {
		log.
			WithField("failFast", info.FailFast).
			WithField("fullMethodName", info.FullMethodName).
			Print("TagRPC")
	}
	return ctx
}

func (h statsHandler) HandleRPC(_ context.Context, rpc stats.RPCStats) {}

func (h statsHandler) TagConn(ctx context.Context, info *stats.ConnTagInfo) context.Context {
	if info != nil {
		log.
			WithField("localAddr", info.LocalAddr.String()).
			WithField("remoteAddr", info.RemoteAddr.String()).
			Print("TagConn")
	}
	return ctx
}

func (h statsHandler) HandleConn(_ context.Context, _ stats.ConnStats) {}
