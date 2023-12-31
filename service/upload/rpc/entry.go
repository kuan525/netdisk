package rpc

import (
	"context"

	upProto "github.com/kuan525/netdisk/client/upload/proto"
	cfg "github.com/kuan525/netdisk/config"
)

type Upload struct {
	upProto.UnimplementedUploadServiceServer
}

// UploadEntry 获取上传入口
func (u *Upload) UploadEntry(ctx context.Context, req *upProto.ReqEntry) (res *upProto.RespEntry, err error) {
	res = new(upProto.RespEntry)
	res.Entry = cfg.UploadEntry
	return
}
