package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	ratelimit2 "github.com/juju/ratelimit"
	cmn "github.com/kuan525/netdisk/common"
	cfg "github.com/kuan525/netdisk/config"
	accProto "github.com/kuan525/netdisk/proto/account"
	dlProto "github.com/kuan525/netdisk/proto/download"
	upProto "github.com/kuan525/netdisk/proto/upload"
	"github.com/kuan525/netdisk/util"
	"github.com/micro/go-micro"
	"github.com/micro/go-plugins/wrapper/breaker/hystrix"
	"github.com/micro/go-plugins/wrapper/ratelimiter/ratelimit"
	"log"
	"net/http"
)

var (
	userCli accProto.UserService
	upCli   upProto.UploadService
	dlCli   dlProto.DownloadService
)

func init() {
	// 配置请求容量以及qps
	bRate := ratelimit2.NewBucketWithRate(100, 100)
	service := micro.NewService(
		micro.Flags(cmn.CustomFlags...),
		micro.WrapClient(ratelimit.NewClientWrapper(bRate, false)), // 加入限流功能，false为不等待(超时返回请求失败)
		micro.WrapClient(hystrix.NewClientWrapper()),               // 加入熔断功能，处理rpc调用失败的情况（cirucuit breaker）
	)

	// 初始化，解析命令行参数
	service.Init()

	cli := service.Client()

	// 初始化 account/upload/download的客户端
	userCli = accProto.NewUserService("go.micro.service.user", cli)
	upCli = upProto.NewUploadService("go.micro.service.upload", cli)
	dlCli = dlProto.NewDownloadService("go.micro.service.download", cli)
}

// SignupHandler 响应注册页面
func SignupHandler(c *gin.Context) {
	c.Redirect(http.StatusOK, "/static/view/signup.html")
}

func DoSignupHandler(c *gin.Context) {
	username := c.Request.FormValue("username")
	passwd := c.Request.FormValue("password")

	resp, err := userCli.Signup(context.TODO(), &accProto.ReqSignup{
		Username: username,
		Password: passwd,
	})

	if err != nil {
		log.Println(err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": resp.Code,
		"msg":  resp.Message,
	})
}

// SigninHandler 响应登陆页面
func SigninHandler(c *gin.Context) {
	c.Redirect(http.StatusFound, "/static/view/signin.html")
}

// DoSigninHandler 处理登陆post请求
func DoSigninHandler(c *gin.Context) {
	username := c.Request.FormValue("username")
	password := c.Request.FormValue("password")

	rpcResp, err := userCli.Signin(context.TODO(), &accProto.ReqSignin{
		Username: username,
		Password: password,
	})

	if err != nil {
		log.Println(err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	if rpcResp.Code != cmn.StatusOk {
		c.JSON(http.StatusOK, gin.H{
			"code": rpcResp.Code,
			"msg":  "登陆失败",
		})
		return
	}

	cliResp := util.RespMsg{
		Code: int(cmn.StatusOk),
		Msg:  "登陆成功",
		Data: struct {
			Location      string
			Username      string
			Token         string
			UploadEntry   string
			DownloadEntry string
		}{
			Location:      "/static/view/home.html",
			Username:      username,
			Token:         rpcResp.Token,
			UploadEntry:   cfg.UploadLBHost,
			DownloadEntry: cfg.DownloadLBHost,
		},
	}
	c.Data(http.StatusOK, "application/json", cliResp.JSONBytes())
}

// UserInfoHandler 查询用户信息
func UserInfoHandler(c *gin.Context) {
	// 1. 解析请求参数
	username := c.Request.FormValue("username")
	resp, err := userCli.UserInfo(context.TODO(), &accProto.ReqUserInfo{
		Username: username,
	})

	if err != nil {
		log.Println(err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	// 2. 组装并响应用户数据
	cliResp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: gin.H{
			"Username":   username,
			"SignupAt":   resp.SignupAt,
			"LastActive": resp.LastActiveAt,
		},
	}
	c.Data(http.StatusOK, "application/json", cliResp.JSONBytes())
}
