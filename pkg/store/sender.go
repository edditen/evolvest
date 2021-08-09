package store

import (
	"context"
	"fmt"
	"github.com/edditen/etlog"
	"github.com/edditen/evolvest/api/pb/evolvest"
	"github.com/edditen/evolvest/pkg/common"
	"github.com/edditen/evolvest/pkg/common/config"
	"github.com/edditen/evolvest/pkg/common/utils"
	"github.com/edditen/evolvest/pkg/runnable"
	"google.golang.org/grpc"
	"os"
	"strings"
	"time"
)

type Sender interface {
	runnable.Runnable
	Send(req *common.TxRequest) error
}

type TxSender struct {
	cfg      *config.Config
	clients  []*EvolvestClient
	shutdown chan interface{}
}

func NewTxSender(cfg *config.Config) *TxSender {
	return &TxSender{
		cfg:     cfg,
		clients: make([]*EvolvestClient, 0),
	}
}

func (ts *TxSender) Init() error {
	servAddrs := os.Getenv(common.EnvAddrs)
	etlog.Log.WithField(common.EnvAddrs, servAddrs).Info("env")
	if servAddrs != "" {
		addrs := strings.Split(servAddrs, ",")
		for _, addr := range addrs {
			client := NewEvolvestClient(addr)
			client.StartClient()
			ts.clients = append(ts.clients, client)
		}
	}
	return nil
}

func (ts *TxSender) Send(req *common.TxRequest) error {
	text := fmt.Sprintf("%d %s %s %s %s",
		req.TxId, common.FlagSync, req.Action, req.Key, utils.Base64Encode(req.Val))
	for _, cli := range ts.clients {
		if cli != nil {
			cli.Push(text)
		}
	}
	return nil
}

func (ts *TxSender) Run() error {
	return nil
}

func (ts *TxSender) Shutdown() {
}

type EvolvestClient struct {
	addr        string
	client      evolvest.EvolvestServiceClient
	reqC        chan string
	sleepSecs   int
	maxInterval int
}

func NewEvolvestClient(addr string) *EvolvestClient {
	return &EvolvestClient{
		addr:        addr,
		reqC:        make(chan string, 100),
		sleepSecs:   1,
		maxInterval: 1024,
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
	ec.reqC <- pushText
}

func (ec *EvolvestClient) Pull() ([]byte, error) {
	resp, err := ec.CallGrpcWithTimeout(func(ctx context.Context) (interface{}, error) {
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
			items := ec.aggr(ec.reqC, 20, 50)
			if len(items) == 0 {
				time.Sleep(time.Second)
				continue
			}
			req := &evolvest.PushRequest{
				TxCmds: items,
			}

			resp, err := ec.CallGrpcWithTimeout(func(ctx context.Context) (interface{}, error) {
				return ec.client.Push(ctx, req)
			})

			log := etlog.Log.WithField("commands", req.TxCmds)
			if err != nil {
				log = log.WithError(err)
				if ec.retry(ec.reqC, req.TxCmds) {
					log.Warn("push tx request to remote failed, reaches max retry times, abandon")
				} else {
					log.Info("push tx request to remote failed, retry")
				}
				continue

			}
			ec.resetRetryCount()
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

func (ec *EvolvestClient) retry(ch chan<- string, items []string) (reachLimit bool) {
	if ec.sleepSecs > ec.maxInterval {
		return true
	}

	go func() {
		time.Sleep(time.Duration(ec.sleepSecs) * time.Second)
		for _, item := range items {
			ch <- item
		}

		ec.sleepSecs *= 2

	}()
	return false
}

func (ec *EvolvestClient) resetRetryCount() {
	ec.sleepSecs = 1
}

func (ec *EvolvestClient) CallGrpcWithTimeout(fn func(ctx context.Context) (interface{}, error)) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return fn(ctx)
}
