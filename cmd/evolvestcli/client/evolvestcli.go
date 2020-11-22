package client

import (
	"context"
	"fmt"
	"github.com/EdgarTeng/evolvest/api/pb/evolvest"
	"google.golang.org/grpc"
	"time"
)

var evolvestClient *EvolvestClient

func init() {
	evolvestClient = NewEvolvestClient()
}

func GetEvolvestClient() *EvolvestClient {
	return evolvestClient
}

type EvolvestClient struct {
	//conn   *grpc.ClientConn
	client evolvest.EvolvestServiceClient
}

func NewEvolvestClient() *EvolvestClient {
	return &EvolvestClient{}
}

func StartClient(addr string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	fmt.Printf("connecting to %s\n", addr)
	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	evolvestClient.client = evolvest.NewEvolvestServiceClient(conn)
	return nil
}

func (e *EvolvestClient) Get(ctx context.Context, key string) (val string, err error) {
	req := &evolvest.GetRequest{Key: key}

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	resp, err := e.client.Get(ctx, req)
	if err != nil {
		return "", err
	}
	return string(resp.GetVal()), nil

}

func (e *EvolvestClient) Set(ctx context.Context, key, val string) (err error) {
	req := &evolvest.SetRequest{Key: key, Val: []byte(val)}

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	_, err = e.client.Set(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (e *EvolvestClient) Del(ctx context.Context, key string) (err error) {
	req := &evolvest.DelRequest{Key: key}

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	_, err = e.client.Del(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (e *EvolvestClient) Sync(ctx context.Context) (values string, err error) {
	req := &evolvest.SyncRequest{}

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	resp, err := e.client.Sync(ctx, req)
	if err != nil {
		return "", err
	}
	return string(resp.Values), nil
}

func (e *EvolvestClient) Push(ctx context.Context, txCmd string) (ok string, err error) {
	req := &evolvest.PushRequest{
		TxCmds: []string{txCmd},
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	_, err = e.client.Push(ctx, req)
	if err != nil {
		return "", err
	}
	return "ok", nil
}
