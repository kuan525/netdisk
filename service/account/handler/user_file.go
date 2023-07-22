package handler

import (
	"context"
	"encoding/json"
	userProto "github.com/kuan525/netdisk/client/account/proto"
	"github.com/kuan525/netdisk/client/dbproxy"
	"github.com/kuan525/netdisk/common"
)

// UserFiles 获取用户文件列表
func (u *User) UserFiles(ctx context.Context, req *userProto.ReqUserFile) (res *userProto.RespUserFile, err error) {
	dbClient := dbproxy.NewDbProxyClient()
	defer dbClient.Conn.Close()

	dbResp, err := dbClient.QueryUserFileMetas(req.Username, int(req.Limit))
	if err != nil || !dbResp.Suc {
		res.Code = common.StatusServerError
		return
	}

	userFiles := dbClient.ToTableUserFile(dbResp.Data)
	data, err := json.Marshal(userFiles)
	if err != nil {
		res.Code = common.StatusServerError
		return
	}

	res.FileData = data
	return
}

// UserFileRename 用户文件重命名
func (u *User) UserFileRename(ctx context.Context, req *userProto.ReqUserFileRename) (res *userProto.RespUserFileRename, err error) {
	dbClient := dbproxy.NewDbProxyClient()
	defer dbClient.Conn.Close()

	dbResp, err := dbClient.RenameFileName(req.Username, req.Filehash, req.NewFileName)
	if err != nil || !dbResp.Suc {
		res.Code = common.StatusServerError
		return
	}

	userFiles := dbClient.ToTableUserFiles(dbResp.Data)
	data, err := json.Marshal(userFiles)
	if err != nil {
		res.Code = common.StatusServerError
		return
	}

	res.FileData = data
	return
}
