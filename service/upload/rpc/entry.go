package rpc

import (
	"context"

	cfg "upload/config"
	upProto "upload/proto"
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
