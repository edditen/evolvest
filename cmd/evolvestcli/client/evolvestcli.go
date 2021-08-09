package client

import (
	"context"
	"fmt"
	"github.com/edditen/evolvest/api/pb/evolvest"
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

func (e *EvolvestClient) Keys(ctx context.Context, pattern string) (keys string, err error) {
	req := &evolvest.KeysRequest{Pattern: pattern}

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	resp, err := e.client.Keys(ctx, req)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v", resp.GetKeys()), nil

}

func (e *EvolvestClient) Pull(ctx context.Context) (values string, err error) {
	req := &evolvest.PullRequest{}

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	resp, err := e.client.Pull(ctx, req)
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
