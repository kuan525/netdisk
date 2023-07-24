package account

import (
	"context"
	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/hashicorp/consul/api"
	userProto "github.com/kuan525/netdisk/client/account/proto"
	ggrpc "google.golang.org/grpc"
	"log"
)

// UserClient 保留conn，后续手动关闭
type UserClient struct {
	Conn   *ggrpc.ClientConn // 上述grpc库底层使用的是google的
	Client userProto.UserServiceClient
}

func NewAccountClient() *UserClient {
	consulCli, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		panic(err)
	}
	r := consul.New(consulCli)

	// new grpc client
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("discovery:///go.kratos.service.user"),
		grpc.WithDiscovery(r),
	)
	if err != nil {
		log.Panic(err)
	}
	//defer conn.Close()
	client := userProto.NewUserServiceClient(conn)

	return &UserClient{
		Client: client,
		Conn:   conn,
	}
}
