package rpc

import (
	"context"
	"github.com/EdgarTeng/evolvest/api/pb/evolvest"
	"github.com/EdgarTeng/evolvest/pkg/common/logger"
	"github.com/EdgarTeng/evolvest/pkg/store"
	"google.golang.org/grpc"
	"net"
)

var evolvestServer *EvolvestServer

func init() {
	evolvestServer = NewEvolvestServer()
}

func GetEvolvestServer() *EvolvestServer {
	return evolvestServer
}

type EvolvestServer struct {
	store store.Store
}

func NewEvolvestServer() *EvolvestServer {
	return &EvolvestServer{
		store: store.GetStore(),
	}
}

func StartServer(port string) error {

	lis, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}

	srv := grpc.NewServer()
	evolvest.RegisterEvolvestServiceServer(srv, GetEvolvestServer())

	go func() {
		logger.Fatal("%v", srv.Serve(lis))
	}()
	return nil
}

func (e *EvolvestServer) Get(ctx context.Context, request *evolvest.GetRequest) (*evolvest.GetResponse, error) {
	logger.Debug("request get params, context: %v, request: %v", ctx, request)
	val, err := e.store.Get(request.GetKey())
	logger.Debug("request get result, value: %v, err: %v", val, err)
	if err != nil {
		return nil, err
	}
	return &evolvest.GetResponse{
		Key: request.GetKey(),
		Val: val,
	}, nil
}

func (e *EvolvestServer) Set(ctx context.Context, request *evolvest.SetRequest) (*evolvest.SetResponse, error) {
	logger.Debug("request set, context params: %v, request: %v", ctx, request)
	oldVal, exists := e.store.Set(request.GetKey(), request.GetVal())
	logger.Debug("request set result, oldVal: %v, exists: %v", oldVal, exists)
	if exists {
		return &evolvest.SetResponse{
			Key:      request.GetKey(),
			ExistVal: true,
			OldVal:   oldVal,
			NewVal:   request.GetVal(),
		}, nil
	}
	return &evolvest.SetResponse{
		Key:      request.GetKey(),
		ExistVal: false,
		NewVal:   request.GetVal(),
	}, nil
}

func (e *EvolvestServer) Del(ctx context.Context, request *evolvest.DelRequest) (*evolvest.DelResponse, error) {
	logger.Debug("request del, context: %v, request: %v", ctx, request)
	oldVal, err := e.store.Del(request.GetKey())
	logger.Debug("request del result, oldVal: %v, err: %v", oldVal, err)
	if err != nil {
		return nil, err
	}
	return &evolvest.DelResponse{
		Key: request.GetKey(),
		Val: oldVal,
	}, nil
}
