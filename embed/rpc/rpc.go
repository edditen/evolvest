package rpc

import (
	"context"
	"encoding/json"
	"github.com/EdgarTeng/evolvest/api/pb/evolvest"
	"github.com/EdgarTeng/evolvest/pkg/common"
	"github.com/EdgarTeng/evolvest/pkg/common/logger"
	"github.com/EdgarTeng/evolvest/pkg/common/utils"
	"github.com/EdgarTeng/evolvest/pkg/store"
	"google.golang.org/grpc"
	"net"
	"regexp"
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

func (e *EvolvestServer) Keys(ctx context.Context, request *evolvest.KeysRequest) (*evolvest.KeysResponse, error) {
	log := logger.WithField("ctx", ctx).WithField("params", request)
	pattern := request.GetPattern()
	r, err := regexp.Compile(pattern)
	if err != nil {
		log.WithError(err).Warn("request keys")
		return nil, err
	}

	allKeys, err := e.store.Keys()
	log.WithField("keys", allKeys).WithError(err).Debug("request keys")
	if err != nil {
		log.WithError(err).Warn("request keys")
		return nil, err
	}

	keys := make([]string, 0)
	for _, key := range allKeys {
		if r.MatchString(key) {
			keys = append(keys, key)
		}
	}

	return &evolvest.KeysResponse{
		Keys: keys,
	}, nil
}

func (e *EvolvestServer) Pull(ctx context.Context, request *evolvest.PullRequest) (*evolvest.PullResponse, error) {
	log := logger.WithField("ctx", ctx).WithField("params", request)
	keys, err := e.store.Keys()

	if err != nil {
		log.WithError(err).Warn("get keys error")
		return nil, err
	}

	values, err := valuesByKeys(keys, e.store)
	if err != nil {
		log.WithError(err).Warn("get values error")
		return nil, err
	}
	data, err := json.Marshal(values)
	if err != nil {
		log.WithError(err).Warn("convert to json error")
		return nil, err
	}

	log.WithField("values", values).Debug("Pull request")

	return &evolvest.PullResponse{
		Values: data,
	}, nil
}

func valuesByKeys(keys []string, s store.Store) (vals map[string]store.DataItem, err error) {
	vals = make(map[string]store.DataItem, 0)
	for _, key := range keys {
		val, err := s.Get(key)
		if err != nil {
			return nil, err
		}
		vals[key] = val
	}
	return vals, nil
}

func (e *EvolvestServer) Push(ctx context.Context, request *evolvest.PushRequest) (*evolvest.PushResponse, error) {
	logger.WithField("ctx", ctx).WithField("params", request).
		Debug("request push")
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
	texts := strings.Fields(strings.TrimSpace(cmdText))
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
