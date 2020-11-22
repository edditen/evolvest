package rpc

import (
	"context"
	"github.com/EdgarTeng/evolvest/api/pb/evolvest"
	"github.com/EdgarTeng/evolvest/pkg/common"
	"github.com/EdgarTeng/evolvest/pkg/common/logger"
	"github.com/EdgarTeng/evolvest/pkg/common/utils"
	"github.com/EdgarTeng/evolvest/pkg/store"
	"google.golang.org/grpc"
	"net"
	"strconv"
	"strings"
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

	logger.Fatal("%v", srv.Serve(lis))
	return nil
}

func (e *EvolvestServer) Get(ctx context.Context, request *evolvest.GetRequest) (*evolvest.GetResponse, error) {
	log := logger.WithField("params", request).
		WithField("ctx", ctx)
	val, err := e.store.Get(request.GetKey())
	log.WithField("val", val).WithError(err).Debug("request get")
	if err != nil {
		return nil, err
	}
	return &evolvest.GetResponse{
		Key: request.GetKey(),
		Val: val.Val,
	}, nil
}

func (e *EvolvestServer) Set(ctx context.Context, request *evolvest.SetRequest) (*evolvest.SetResponse, error) {
	log := logger.WithField("params", request).
		WithField("ctx", ctx)
	oldVal, exists := e.store.Set(request.GetKey(), store.DataItem{
		Val: request.GetVal(),
		Ver: utils.CurrentMillis(),
	})
	log.WithField("oldVal", oldVal).
		WithField("exists", exists).
		Debug("request set")
	if exists {
		return &evolvest.SetResponse{
			Key:      request.GetKey(),
			ExistVal: true,
			OldVal:   oldVal.Val,
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
	log := logger.WithField("params", request).
		WithField("ctx", ctx)
	oldVal, err := e.store.Del(request.GetKey(), utils.GenerateId())
	log.WithField("oldVal", oldVal).
		WithError(err).
		Debug("request del")
	if err != nil {
		return nil, err
	}
	return &evolvest.DelResponse{
		Key: request.GetKey(),
		Val: oldVal.Val,
	}, nil
}

func (e *EvolvestServer) Sync(ctx context.Context, request *evolvest.SyncRequest) (*evolvest.SyncResponse, error) {
	log := logger.WithField("params", request).
		WithField("ctx", ctx)
	values, err := e.store.Serialize()
	log.WithField("values", values).
		WithError(err).
		Debug("request sync")
	if err != nil {
		return nil, err
	}
	return &evolvest.SyncResponse{
		Values: values,
	}, nil
}

func (e *EvolvestServer) Push(ctx context.Context, request *evolvest.PushRequest) (*evolvest.PushResponse, error) {
	logger.WithField("params", request).
		WithField("ctx", ctx).Debug("access")
	for _, req := range request.TxCmds {
		if txReq := parseCmd(req); txReq != nil {
			store.Submit(txReq)
		}
	}
	return &evolvest.PushResponse{
		Ok: true,
	}, nil
}

func parseCmd(cmdText string) *common.TxRequest {
	log := logger.WithField("cmdText", cmdText)
	texts := strings.Split(cmdText, " ")
	if len(texts) < 4 {
		log.Warn("parse cmd error, missing required")
		return nil
	}

	txReq := &common.TxRequest{}

	id, err := strconv.Atoi(texts[0])
	if err != nil {
		log.Warn("parse cmd error, txid is wrong format")
		return nil
	}

	txReq.TxId = int64(id)
	txReq.Flag = common.FlagSync
	txReq.Key = texts[3]

	switch texts[2] {
	case common.DEL:
		txReq.Action = common.DEL
	case common.SET:
		txReq.Action = common.SET
		if len(texts) == 4 {
			txReq.Val = []byte{}
		} else if len(texts) == 5 {
			if val, err := utils.Base64Decode(texts[4]); err != nil {
				log.Warn("parse cmd error, value is wrong format")
				return nil
			} else {
				txReq.Val = val
			}
		} else {
			log.Warn("parse cmd error, more than one values")
			return nil
		}
	default:
		log.Warn("parse cmd error, cmd not support")
		return nil
	}
	log.WithField("req", txReq).Debug("parsed request")
	return txReq
}
