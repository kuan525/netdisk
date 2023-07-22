package upload

import (
	"context"
	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/hashicorp/consul/api"
	uploadProto "github.com/kuan525/netdisk/client/upload/proto"
	ggrpc "google.golang.org/grpc"
	"log"
)

// downloadClient 保留conn，后续手动关闭
type uploadClient struct {
	Conn   *ggrpc.ClientConn // 上述grpc库底层使用的是google的
	Client uploadProto.UploadServiceClient
}

func NewAccountClient() *uploadClient {
	consulCli, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		panic(err)
	}
	r := consul.New(consulCli)

	// new grpc client
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("discovery:///go.micro.service.upload"),
		grpc.WithDiscovery(r),
	)
	if err != nil {
		log.Panic(err)
	}
	//defer conn.Close()
	client := uploadProto.NewUploadServiceClient(conn)

	return &uploadClient{
		Client: client,
		Conn:   conn,
	}
}
