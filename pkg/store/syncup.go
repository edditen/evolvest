package store

import (
	"context"
	"fmt"
	"github.com/EdgarTeng/evolvest/api/pb/evolvest"
	"github.com/EdgarTeng/evolvest/pkg/common"
	"github.com/EdgarTeng/evolvest/pkg/common/config"
	"github.com/EdgarTeng/evolvest/pkg/common/logger"
	"github.com/EdgarTeng/evolvest/pkg/common/utils"
	"google.golang.org/grpc"
	"os"
	"path"
	"strings"
	"time"
)

var (
	reqChan         chan *common.TxRequest
	fileWriter      *os.File
	evolvestClients []*EvolvestClient
)

func init() {
	reqChan = make(chan *common.TxRequest, 100)
	evolvestClients = make([]*EvolvestClient, 0)
}

func InitSyncUp() {
	dataDir := config.Config().DataDir
	if err := os.MkdirAll(dataDir, os.ModePerm); err != nil {
		logger.Warn("mkdir error, %v", err)
	}

	filename := path.Join(dataDir, common.FileTx)
	f, err := os.OpenFile(filename,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger.WithError(err).Fatal("open tx file failed")
	}
	fileWriter = f
	Process()

	// clients
	servAddrs := os.Getenv(common.EnvAddrs)
	logger.WithField(common.EnvAddrs, servAddrs).Info("env")
	if servAddrs != "" {
		addrs := strings.Split(servAddrs, ",")
		for _, addr := range addrs {
			client := NewEvolvestClient(addr)
			client.StartClient()
			evolvestClients = append(evolvestClients, client)
		}
	}

}

func Submit(req *common.TxRequest) {
	reqChan <- req
}

func Process() {
	go func() {
		for {
			req := <-reqChan
			setToStore(req)
			appendTxFile(req)
			if req.Flag == common.FlagReq {
				pushToRemote(req)
			}
		}
	}()
}

func setToStore(req *common.TxRequest) {
	switch req.Action {
	case common.SET:
		GetStore().Set(req.Key, DataItem{
			Val: req.Val,
			Ver: req.TxId,
		})
	case common.DEL:
		GetStore().Del(req.Key, req.TxId)
	}
}

func appendTxFile(req *common.TxRequest) {
	text := fmt.Sprintf("%d %s %s %s %s\n",
		req.TxId, req.Flag, req.Action, req.Key, utils.Base64Encode(req.Val))
	if _, err := fileWriter.WriteString(text); err != nil {
		logger.WithError(err).
			WithField("text", text).
			Warn("append text to tx file failed")
	}
}

func pushToRemote(req *common.TxRequest) {
	text := fmt.Sprintf("%d %s %s %s %s",
		req.TxId, common.FlagSync, req.Action, req.Key, utils.Base64Encode(req.Val))
	for _, cli := range evolvestClients {
		if cli != nil {
			cli.Push(text)
		}
	}
}

type EvolvestClient struct {
	addr     string
	client   evolvest.EvolvestServiceClient
	pushChan chan string
}

func NewEvolvestClient(addr string) *EvolvestClient {
	return &EvolvestClient{
		addr:     addr,
		pushChan: make(chan string, 100),
	}
}

func (e *EvolvestClient) StartClient() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	logger.Info("connecting to %s", e.addr)
	conn, err := grpc.DialContext(ctx, e.addr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	e.client = evolvest.NewEvolvestServiceClient(conn)
	e.Process()
	return nil
}

func (e *EvolvestClient) Push(pushText string) {
	e.pushChan <- pushText
}

func (e *EvolvestClient) Process() {
	go func() {
		for {
			pushText := <-e.pushChan
			req := &evolvest.PushRequest{
				TxCmds: []string{pushText},
			}

			resp, err := CallGrpcWithTimeout(func(ctx context.Context) (interface{}, error) {
				return e.client.Push(ctx, req)
			})
			if err != nil {
				logger.WithError(err).Warn("push tx request to remote failed")
				continue
			}
			logger.WithField("response", resp).Debug("push to remote response")
		}
	}()
}

func CallGrpcWithTimeout(fn func(ctx context.Context) (interface{}, error)) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return fn(ctx)
}
