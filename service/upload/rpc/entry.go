package rpc

import (
	"context"

	upProto "github.com/kuan525/netdisk/proto/upload"
	cfg "github.com/kuan525/netdisk/service/upload/config"
)

type Upload struct{}

// UploadEntry 获取上传入口
func (u *Upload) UploadEntry(
	ctx context.Context,
	req *upProto.ReqEntry,
	res *upProto.RespEntry) error {

	res.Entry = cfg.UploadEntry
	return nil
}
