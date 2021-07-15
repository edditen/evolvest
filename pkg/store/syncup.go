package store

import (
	"context"
	"fmt"
	"github.com/EdgarTeng/etlog"
	"github.com/EdgarTeng/evolvest/api/pb/evolvest"
	"github.com/EdgarTeng/evolvest/pkg/common"
	"github.com/EdgarTeng/evolvest/pkg/common/config"
	"github.com/EdgarTeng/evolvest/pkg/common/utils"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"os"
	"path"
	"strings"
	"time"
)

type Syncer struct {
	conf    *config.Config
	reqC    chan *common.TxRequest
	writer  *os.File
	clients []*EvolvestClient
	Store   Store
}

func NewSyncer(conf *config.Config, store Store) *Syncer {
	return &Syncer{
		conf:    conf,
		Store:   store,
		reqC:    make(chan *common.TxRequest, 1000),
		clients: make([]*EvolvestClient, 0),
	}
}

func (s *Syncer) Init() error {
	dataDir := s.conf.DataDir
	if err := os.MkdirAll(dataDir, os.ModePerm); err != nil {
		return errors.Wrap(err, "init syncUp error")
	}
	filename := path.Join(dataDir, common.FileTx)
	f, err := os.OpenFile(filename,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return errors.Wrap(err, "open tx file failed")
	}
	s.writer = f
	return nil
}

func (s *Syncer) Run() error {
	s.Process()

	// clients
	servAddrs := os.Getenv(common.EnvAddrs)
	etlog.Log.WithField(common.EnvAddrs, servAddrs).Info("env")
	if servAddrs != "" {
		addrs := strings.Split(servAddrs, ",")
		for _, addr := range addrs {
			client := NewEvolvestClient(addr)
			client.StartClient()
			s.clients = append(s.clients, client)
		}
	}

	return nil

}

func (s *Syncer) Shutdown() {
}

func (s *Syncer) Submit(req *common.TxRequest) {
	s.reqC <- req
}

func (s *Syncer) Process() {
	go func() {
		for {
			req := <-s.reqC
			s.setToStore(req)
			s.appendTxFile(req)
			if req.Flag == common.FlagReq {
				go s.pushToRemote(req)
			}
		}
	}()
}

func (s *Syncer) setToStore(req *common.TxRequest) {
	switch req.Action {
	case common.SET:
		s.Store.Set(req.Key, DataItem{
			Val: req.Val,
			Ver: req.TxId,
		})
	case common.DEL:
		s.Store.Del(req.Key, req.TxId)
	}
}

func (s *Syncer) appendTxFile(req *common.TxRequest) {
	text := fmt.Sprintf("%d %s %s %s %s\n",
		req.TxId, req.Flag, req.Action, req.Key, utils.Base64Encode(req.Val))
	if _, err := s.writer.WriteString(text); err != nil {
		etlog.Log.WithError(err).
			WithField("text", text).
			Warn("append text to tx file failed")
	}
}

func (s *Syncer) pushToRemote(req *common.TxRequest) {
	text := fmt.Sprintf("%d %s %s %s %s",
		req.TxId, common.FlagSync, req.Action, req.Key, utils.Base64Encode(req.Val))
	for _, cli := range s.clients {
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

func (ec *EvolvestClient) StartClient() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	etlog.Log.WithField("addr", ec.addr).Info("connecting")
	conn, err := grpc.DialContext(ctx, ec.addr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	ec.client = evolvest.NewEvolvestServiceClient(conn)
	ec.Process()
	return nil
}

func (ec *EvolvestClient) Push(pushText string) {
	ec.pushChan <- pushText
}

func (ec *EvolvestClient) Pull() ([]byte, error) {
	resp, err := CallGrpcWithTimeout(func(ctx context.Context) (interface{}, error) {
		return ec.client.Pull(ctx, &evolvest.PullRequest{})
	})
	if err != nil {
		return nil, err
	}

	pullResp, ok := resp.(*evolvest.PullResponse)
	if !ok {
		return nil, fmt.Errorf("type convert error")
	}
	return pullResp.Values, nil

}

func (ec *EvolvestClient) Process() {
	go func() {
		for {
			items := ec.aggr(ec.pushChan, 20, 50)
			if len(items) == 0 {
				time.Sleep(time.Second)
				continue
			}
			req := &evolvest.PushRequest{
				TxCmds: items,
			}

			resp, err := CallGrpcWithTimeout(func(ctx context.Context) (interface{}, error) {
				return ec.client.Push(ctx, req)
			})

			log := etlog.Log.WithField("commands", req.TxCmds)
			if err != nil {
				log = log.WithError(err)
				if retry(ec.pushChan, req.TxCmds) {
					log.Warn("push tx request to remote failed, reaches max retry times, abandon")
				} else {
					log.Info("push tx request to remote failed, retry")
				}
				continue

			}
			resetRetryCount()
			log.WithField("remote_addr", ec.addr).
				WithField("response", resp).
				Debug("push to remote success")
		}
	}()
}

func (ec *EvolvestClient) aggr(ch <-chan string, maxCount int, maxWaitMillis int64) []string {
	items := make([]string, 0)
	timeout := time.Duration(maxWaitMillis) * time.Millisecond
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	for i := 0; i < maxCount; i++ {
		select {
		case item := <-ch:
			items = append(items, item)
		case <-ctx.Done():
			break
		}
	}
	return items

}

var (
	sleepSecs   = 1
	maxInterval = 1024
)

func retry(ch chan<- string, items []string) (reachLimit bool) {
	if sleepSecs > maxInterval {
		return true
	}

	go func() {
		time.Sleep(time.Duration(sleepSecs) * time.Second)
		for _, item := range items {
			ch <- item
		}

		sleepSecs *= 2

	}()
	return false
}

func resetRetryCount() {
	sleepSecs = 1
}

func CallGrpcWithTimeout(fn func(ctx context.Context) (interface{}, error)) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return fn(ctx)
}
