package handler

import (
	"context"
	"encoding/json"
	"github.com/kuan525/netdisk/common"
	dbcli "github.com/kuan525/netdisk/dbclient"
	accProto "github.com/kuan525/netdisk/proto/account"
)

// UserFiles 获取用户文件列表
func (u *User) UserFiles(ctx context.Context, req *accProto.ReqUserFile, res *accProto.RespUserFile) error {
	dbResp, err := dbcli.QueryUserFileMetas(req.Username, int(req.Limit))
	if err != nil || !dbResp.Suc {
		res.Code = common.StatusServerError
		return err
	}

	userFiles := dbcli.ToTableUserFile(dbResp.Data)
	data, err := json.Marshal(userFiles)
	if err != nil {
		res.Code = common.StatusServerError
		return nil
	}

	res.FileData = data
	return nil
}

// UserFileRename 用户文件重命名
func (u *User) UserFileRename(
	ctx context.Context,
	req *accProto.ReqUserFileRename,
	res *accProto.RespUserFileRename) error {

	dbResp, err := dbcli.RenameFileName(req.Username, req.Filehash, req.NewFileName)
	if err != nil || !dbResp.Suc {
		res.Code = common.StatusServerError
		return err
	}

	userFiles := dbcli.ToTableUserFiles(dbResp.Data)
	data, err := json.Marshal(userFiles)
	if err != nil {
		res.Code = common.StatusServerError
		return nil
	}

	res.FileData = data
	return nil
}
