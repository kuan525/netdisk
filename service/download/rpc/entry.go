package rpc

import (
	"context"

	dlProto "github.com/kuan525/netdisk/client/download/proto"
	cfg "github.com/kuan525/netdisk/config"
)

type Download struct {
	dlProto.UnimplementedDownloadServiceServer
}

// DownloadEntry 获取下载入口
func (u *Download) DownloadEntry(ctx context.Context, req *dlProto.ReqEntry) (res *dlProto.RespEntry, err error) {
	res.Entry = cfg.DownloadEntry
	return
}
