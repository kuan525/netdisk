package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/kuan525/netdisk/client/account"
	userProto "github.com/kuan525/netdisk/client/account/proto"
	cmn "github.com/kuan525/netdisk/common"
	cfg "github.com/kuan525/netdisk/config"
	"github.com/kuan525/netdisk/util"

	"log"
	"net/http"
)

// SignupHandler 响应注册页面
func SignupHandler(c *gin.Context) {
	c.Redirect(http.StatusOK, "/static/view/signup.html")
}

func DoSignupHandler(c *gin.Context) {
	username := c.Request.FormValue("username")
	passwd := c.Request.FormValue("password")

	userClient := account.NewAccountClient()
	defer userClient.Conn.Close()

	resp, err := userClient.Client.Signup(context.TODO(), &userProto.ReqSignup{
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

	userClient := account.NewAccountClient()
	defer userClient.Conn.Close()

	rpcResp, err := userClient.Client.Signin(context.TODO(), &userProto.ReqSignin{
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
	userClient := account.NewAccountClient()
	defer userClient.Conn.Close()

	// 1. 解析请求参数
	username := c.Request.FormValue("username")
	resp, err := userClient.Client.UserInfo(context.TODO(), &userProto.ReqUserInfo{
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
