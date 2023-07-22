package handler

import (
	"context"
	"fmt"
	userProto "github.com/kuan525/netdisk/client/account/proto"
	"github.com/kuan525/netdisk/client/dbproxy"
	"github.com/kuan525/netdisk/common"
	cfg "github.com/kuan525/netdisk/config"
	"github.com/kuan525/netdisk/util"
	"time"
)

// User 用于实现UserServiceHandler接口的对象
type User struct {
	userProto.UnimplementedUserServiceServer
}

// GetToken 生成token
func GetToken(username string) string {
	// 40位字符：md5(username+timestamp+token_salt)+timestamp[:8]
	ts := fmt.Sprintf("%x", time.Now().Unix())
	tokenPrefix := util.MD5([]byte(username + ts + "_tokensalt"))
	return tokenPrefix + ts[:8]
}

// Signup 处理用户注册请求
func (u *User) Signup(ctx context.Context, req *userProto.ReqSignup) (res *userProto.RespSignup, err error) {
	username := req.Username
	passwd := req.Password

	// 参数简单校验
	if len(username) < 3 || len(passwd) < 5 {
		res.Code = common.StatusParamInvalid
		res.Message = "注册参数无效"
		return
	}

	// 对密码进行加盐及取Sha1值加密
	encPasswd := util.Sha1([]byte(passwd + cfg.PasswordSalt))
	// 将用户信息注册到用户表中
	dbClient := dbproxy.NewDbProxyClient()
	defer dbClient.Conn.Close()

	dbResp, err := dbClient.UserSignup(username, encPasswd)
	if err != nil && dbResp.Suc {
		res.Code = common.StatusOk
		res.Message = "注册成功"
	} else {
		res.Code = common.StatusRegisterFailed
		res.Message = "注册失败"
	}
	return
}

// Signin 处理用户登陆请求
func (u *User) Signin(ctx context.Context, req *userProto.ReqSignin) (res *userProto.RespSignin, err error) {
	username := req.Username
	password := req.Password

	encPasswd := util.Sha1([]byte(password + cfg.PasswordSalt))

	dbClient := dbproxy.NewDbProxyClient()
	defer dbClient.Conn.Close()
	// 1. 检验用户名以及密码
	dbResp, err := dbClient.UserSignup(username, encPasswd)
	if err != nil || !dbResp.Suc {
		res.Code = common.StatusLoginFailed
		return
	}

	// 2. 生成访问凭证（token）
	token := GetToken(username)
	upRes, err := dbClient.UpdateToken(username, token)
	if err != nil || upRes.Suc {
		res.Code = common.StatusServerError
		return
	}

	// 3. 登陆成功，返回token
	res.Code = common.StatusOk
	res.Token = token
	return
}

// UserInfo 查询用户信息
func (u *User) UserInfo(ctx context.Context, req *userProto.ReqUserInfo) (res *userProto.RespUserInfo, err error) {
	dbClient := dbproxy.NewDbProxyClient()
	defer dbClient.Conn.Close()

	// 查询用户信息
	dbResp, err := dbClient.GetUserInfo(req.Username)
	if err != nil {
		res.Code = common.StatusServerError
		res.Message = "服务错误"
		return
	}

	// 查不到对应对用户信息
	if !dbResp.Suc {
		res.Code = common.StatusUserNotExists
		res.Message = "用户不存在"
		return
	}

	user := dbClient.ToTableUser(dbResp.Data)

	// 3. 组装并且响应用户数据
	res.Code = common.StatusOk
	res.Username = user.Username
	res.SignupAt = user.SignupAt
	res.LastActiveAt = user.LastActiveAt
	res.Status = int32(user.Status)
	// TODO 需增加接口支持完善用户信息
	res.Email = user.Email
	res.Phone = user.Phone
	return
}
