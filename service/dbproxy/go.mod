module dbproxy

go 1.14

replace (
	google.golang.org/grpc => google.golang.org/grpc v1.26.0
	netdisk/common => ../../common
	netdisk/dbclient => ../../dbclient
)

require (
	github.com/micro/cli v0.2.0
	github.com/micro/go-micro v1.18.0
	github.com/micro/go-plugins/registry/kubernetes v0.0.0-20200119172437-4fe21aa238fd
)

require (
	github.com/coreos/etcd v3.3.18+incompatible // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/json-iterator/go v1.1.9 // indirect
	github.com/lucas-clemente/quic-go v0.14.1 // indirect
	github.com/miekg/dns v1.1.43 // indirect
	github.com/nats-io/nats-server/v2 v2.1.6 // indirect
	github.com/onsi/ginkgo/v2 v2.9.5 // indirect
	github.com/tmc/grpc-websocket-proxy v0.0.0-20200122045848-3419fae592fc // indirect
	go.etcd.io/bbolt v1.3.5 // indirect
	go.uber.org/multierr v1.5.0 // indirect
	go.uber.org/zap v1.13.0 // indirect
	golang.org/x/crypto v0.4.0 // indirect
	google.golang.org/genproto v0.0.0-20230306155012-7f2fa6fef1f4 // indirect
	google.golang.org/grpc v1.54.1 // indirect
	netdisk/common v0.0.0-00010101000000-000000000000
	netdisk/dbclient v0.0.0-00010101000000-000000000000
)
