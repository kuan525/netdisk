package download

import (
	"context"
	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/hashicorp/consul/api"
	downloadProto "github.com/kuan525/netdisk/client/download/proto"
	ggrpc "google.golang.org/grpc"
	"log"
)

// downloadClient 保留conn，后续手动关闭
type downloadClient struct {
	Conn   *ggrpc.ClientConn // 上述grpc库底层使用的是google的
	Client downloadProto.DownloadServiceClient
}

func NewAccountClient() *downloadClient {
	consulCli, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		panic(err)
	}
	r := consul.New(consulCli)

	// new grpc client
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("discovery:///go.micro.service.download"),
		grpc.WithDiscovery(r),
	)
	if err != nil {
		log.Panic(err)
	}
	//defer conn.Close()
	client := downloadProto.NewDownloadServiceClient(conn)

	return &downloadClient{
		Client: client,
		Conn:   conn,
	}
}
