package rpc

import (
	"context"
	dlProto "github.com/kuan525/netdisk/proto/download"
	cfg "github.com/kuan525/netdisk/service/download/config"
)

type Download struct{}

// DownloadEntry 获取下载入口
func (u *Download) DownloadEntry(
	ctx context.Context,
	req *dlProto.ReqEntry,
	res *dlProto.RespEntry) {

	res.Entry = cfg.DownloadEntry
	return
}
