package mq

import cmn "github.com/kuan525/netdisk/common"

// TransferData 即将要写到rabbitmq的数据的结构体
type TransferData struct {
	FileHash      string
	CurLocation   string
	DestLocation  string
	DestStoreType cmn.StoreType
}
