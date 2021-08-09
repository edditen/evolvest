package rpc

import (
	"context"
	"encoding/json"
	"github.com/edditen/etlog"
	"github.com/edditen/evolvest/api/pb/evolvest"
	"github.com/edditen/evolvest/pkg/common"
	"github.com/edditen/evolvest/pkg/common/config"
	"github.com/edditen/evolvest/pkg/common/utils"
	"github.com/edditen/evolvest/pkg/store"
	"google.golang.org/grpc"
	"log"
	"net"
	"regexp"
	"strconv"
	"strings"
)

type SyncServer struct {
	cfg    *config.Config
	syncer *store.Syncer
}

func NewSyncServer(conf *config.Config, syncer *store.Syncer) *SyncServer {
	return &SyncServer{
		cfg:    conf,
		syncer: syncer,
	}
}

func (es *SyncServer) Init() error {
	return nil
}

func (es *SyncServer) Run() error {
	addr := es.cfg.Host + ":" + es.cfg.SyncPort
	log.Println("listen sync server at", addr)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	srv := grpc.NewServer()
	evolvest.RegisterEvolvestServiceServer(srv, es)

	err = srv.Serve(lis)
	etlog.Log.WithError(err).Fatal("serve failed")
	return nil
}

func (es *SyncServer) Shutdown() {
}

func (es *SyncServer) Keys(ctx context.Context, request *evolvest.KeysRequest) (*evolvest.KeysResponse, error) {
	log := etlog.Log.WithField("ctx", ctx).WithField("params", request)
	pattern := request.GetPattern()
	r, err := regexp.Compile(pattern)
	if err != nil {
		log.WithError(err).Warn("request keys")
		return nil, err
	}

	allKeys, err := es.syncer.Store.Keys()
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

func (es *SyncServer) Pull(ctx context.Context, request *evolvest.PullRequest) (*evolvest.PullResponse, error) {
	log := etlog.Log.WithField("ctx", ctx).WithField("params", request)
	keys, err := es.syncer.Store.Keys()

	if err != nil {
		log.WithError(err).Warn("get keys error")
		return nil, err
	}

	values, err := valuesByKeys(keys, es.syncer.Store)
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

func (es *SyncServer) Push(ctx context.Context, request *evolvest.PushRequest) (*evolvest.PushResponse, error) {
	etlog.Log.WithField("ctx", ctx).WithField("params", request).
		Debug("request push")
	for _, req := range request.TxCmds {
		if txReq := parseCmd(req); txReq != nil {
			es.syncer.Submit(txReq)
		}
	}
	return &evolvest.PushResponse{
		Ok: true,
	}, nil
}

func parseCmd(cmdText string) *common.TxRequest {
	log := etlog.Log.WithField("cmdText", cmdText)
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
