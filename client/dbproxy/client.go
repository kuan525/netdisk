package dbproxy

import (
	"context"
	"encoding/json"
	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/hashicorp/consul/api"
	dbProto "github.com/kuan525/netdisk/client/dbproxy/proto"
	"github.com/kuan525/netdisk/common/orm"
	"github.com/mitchellh/mapstructure"
	ggrpc "google.golang.org/grpc"
	"log"
)

// DbClient 保留conn，后续手动关闭
type DbClient struct {
	Conn   *ggrpc.ClientConn // 上述grpc库底层使用的是google的
	Client dbProto.DBProxyServiceClient
}

// FileMeta : 文件元信息结构
type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

// NewDbProxyClient 新建一个DbProxyClient
func NewDbProxyClient() *DbClient {
	consulCli, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		panic(err)
	}
	r := consul.New(consulCli)

	// new grpc client
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("discovery:///go.micro.service.dbproxy"),
		grpc.WithDiscovery(r),
	)
	if err != nil {
		log.Panic(err)
	}
	//defer conn.Close()
	client := dbProto.NewDBProxyServiceClient(conn)

	return &DbClient{
		Client: client,
		Conn:   conn,
	}
}

func (c *DbClient) TableFileToFileMeta(tfile orm.TableFile) FileMeta {
	return FileMeta{
		FileSha1: tfile.FileHash,
		FileName: tfile.FileName.String,
		FileSize: tfile.FileSize.Int64,
		Location: tfile.FileAddr.String,
	}
}

// execAction : 向dbproxy请求执行action
func (c *DbClient) execAction(funcName string, paramJson []byte) (*dbProto.RespExec, error) {
	return c.Client.ExecuteAction(context.TODO(), &dbProto.ReqExec{
		Action: []*dbProto.SingleAction{
			&dbProto.SingleAction{
				Name:   funcName,
				Params: paramJson,
			},
		},
	})
}

// parseBody : 转换rpc返回的结果
func (c *DbClient) parseBody(resp *dbProto.RespExec) *orm.ExecResult {
	if resp == nil || resp.Data == nil {
		return nil
	}
	var resList []orm.ExecResult
	_ = json.Unmarshal(resp.Data, &resList)
	// TODO:
	if len(resList) > 0 {
		return &resList[0]
	}
	return nil
}

func (c *DbClient) ToTableUser(src interface{}) orm.TableUser {
	user := orm.TableUser{}
	mapstructure.Decode(src, &user)
	return user
}

func (c *DbClient) ToTableFile(src interface{}) orm.TableFile {
	file := orm.TableFile{}
	mapstructure.Decode(src, &file)
	return file
}

func (c *DbClient) ToTableFiles(src interface{}) []orm.TableFile {
	var file []orm.TableFile
	mapstructure.Decode(src, &file)
	return file
}

func (c *DbClient) ToTableUserFile(src interface{}) orm.TableUserFile {
	ufile := orm.TableUserFile{}
	mapstructure.Decode(src, &ufile)
	return ufile
}

func (c *DbClient) ToTableUserFiles(src interface{}) []orm.TableUserFile {
	var ufile []orm.TableUserFile
	mapstructure.Decode(src, &ufile)
	return ufile
}

func (c *DbClient) GetFileMeta(filehash string) (*orm.ExecResult, error) {
	uInfo, _ := json.Marshal([]interface{}{filehash})
	res, err := c.execAction("/file/GetFileMeta", uInfo)
	return c.parseBody(res), err
}

func (c *DbClient) GetFileMetaList(limitCnt int) (*orm.ExecResult, error) {
	uInfo, _ := json.Marshal([]interface{}{limitCnt})
	res, err := c.execAction("/file/GetFileMetaList", uInfo)
	return c.parseBody(res), err
}

// OnFileUploadFinished : 新增/更新文件元信息到mysql中
func (c *DbClient) OnFileUploadFinished(fmeta FileMeta) (*orm.ExecResult, error) {
	uInfo, _ := json.Marshal([]interface{}{fmeta.FileSha1, fmeta.FileName, fmeta.FileSize, fmeta.Location})
	res, err := c.execAction("/file/OnFileUploadFinished", uInfo)
	return c.parseBody(res), err
}

func (c *DbClient) UpdateFileLocation(filehash, location string) (*orm.ExecResult, error) {
	uInfo, _ := json.Marshal([]interface{}{filehash, location})
	res, err := c.execAction("/file/UpdateFileLocation", uInfo)
	return c.parseBody(res), err
}

func (c *DbClient) UserSignup(username, encPasswd string) (*orm.ExecResult, error) {
	uInfo, _ := json.Marshal([]interface{}{username, encPasswd})
	res, err := c.execAction("/user/UserSignup", uInfo)
	return c.parseBody(res), err
}

func (c *DbClient) UserSignin(username, encPasswd string) (*orm.ExecResult, error) {
	uInfo, _ := json.Marshal([]interface{}{username, encPasswd})
	res, err := c.execAction("/user/UserSignin", uInfo)
	return c.parseBody(res), err
}

func (c *DbClient) GetUserInfo(username string) (*orm.ExecResult, error) {
	uInfo, _ := json.Marshal([]interface{}{username})
	res, err := c.execAction("/user/GetUserInfo", uInfo)
	return c.parseBody(res), err
}

func (c *DbClient) UserExist(username string) (*orm.ExecResult, error) {
	uInfo, _ := json.Marshal([]interface{}{username})
	res, err := c.execAction("/user/UserExist", uInfo)
	return c.parseBody(res), err
}

func (c *DbClient) UpdateToken(username, token string) (*orm.ExecResult, error) {
	uInfo, _ := json.Marshal([]interface{}{username, token})
	res, err := c.execAction("/user/UpdateToken", uInfo)
	return c.parseBody(res), err
}

func (c *DbClient) QueryUserFileMeta(username, filehash string) (*orm.ExecResult, error) {
	uInfo, _ := json.Marshal([]interface{}{username, filehash})
	res, err := c.execAction("/ufile/QueryUserFileMeta", uInfo)
	return c.parseBody(res), err
}

func (c *DbClient) QueryUserFileMetas(username string, limit int) (*orm.ExecResult, error) {
	uInfo, _ := json.Marshal([]interface{}{username, limit})
	res, err := c.execAction("/ufile/QueryUserFileMetas", uInfo)
	return c.parseBody(res), err
}

// OnUserFileUploadFinished : 新增/更新文件元信息到mysql中
func (c *DbClient) OnUserFileUploadFinished(username string, fmeta FileMeta) (*orm.ExecResult, error) {
	uInfo, _ := json.Marshal([]interface{}{username, fmeta.FileSha1,
		fmeta.FileName, fmeta.FileSize})
	res, err := c.execAction("/ufile/OnUserFileUploadFinished", uInfo)
	return c.parseBody(res), err
}

func (c *DbClient) RenameFileName(username, filehash, filename string) (*orm.ExecResult, error) {
	uInfo, _ := json.Marshal([]interface{}{username, filehash, filename})
	res, err := c.execAction("/ufile/RenameFileName", uInfo)
	return c.parseBody(res), err
}
