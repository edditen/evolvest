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
	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure(), grpc.WithBlock())
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
	//log.Printf("get response: %v\n", resp)
	if err != nil {
		return "", err
	}
	return resp.GetVal(), nil

}

func (e *EvolvestClient) Set(ctx context.Context, key, val string) (err error) {
	req := &evolvest.SetRequest{Key: key, Val: val}

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	_, err = e.client.Set(ctx, req)
	//log.Printf("set response: %v\n", resp)
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
	//log.Printf("del response: %v\n", resp)
	if err != nil {
		return err
	}
	return nil
}
