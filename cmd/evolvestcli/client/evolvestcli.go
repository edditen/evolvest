package client

import (
	"context"
	"github.com/EdgarTeng/evolvest/api/pb/evolvest"
	"google.golang.org/grpc"
	"log"
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

func StartClient(port string) {
	conn, err := grpc.Dial(port, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalln(err)
	}
	//evolvestClient.conn = conn
	evolvestClient.client = evolvest.NewEvolvestServiceClient(conn)
}

func (e *EvolvestClient) Get(ctx context.Context, key string) (val string, err error) {
	req := &evolvest.GetRequest{Key: key}

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	resp, err := e.client.Get(ctx, req)
	log.Printf("response: %v\n", resp)
	if err != nil {
		return "", err
	}
	return resp.GetVal(), nil

}

func (e *EvolvestClient) Set(ctx context.Context, key, val string) (err error) {
	req := &evolvest.SetRequest{Key: key}

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	resp, err := e.client.Set(ctx, req)
	log.Printf("response: %v\n", resp)
	if err != nil {
		return err
	}
	return nil
}

func (e *EvolvestClient) Del(ctx context.Context, key string) (err error) {
	req := &evolvest.DelRequest{Key: key}

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	resp, err := e.client.Del(ctx, req)
	log.Printf("response: %v\n", resp)
	if err != nil {
		return err
	}
	return nil
}
