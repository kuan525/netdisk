package rpc

import (
	"context"
	cfg "download/config"
	dlProto "github.com/kuan525/netdisk/proto/download"
)

type Download struct{}

// DownloadEntry 获取下载入口
func (u *Download) DownloadEntry(
	ctx context.Context,
	req *dlProto.RespEntry,
	res *dlProto.RespEntry) {

	res.Entry = cfg.DownloadEntry
	return
}
