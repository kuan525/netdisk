package route

import (
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/kuan525/netdisk/asset"
	"github.com/kuan525/netdisk/service/apigw/handler"
	assetfs "github.com/moxiaomomo/go-bindata-assetfs"
	"net/http"
	"strings"
)

type binaryFileSystem struct {
	fs http.FileSystem
}

func (b *binaryFileSystem) Open(name string) (http.File, error) {
	return b.fs.Open(name)
}

func (b *binaryFileSystem) Exists(prefix string, filepath string) bool {
	if p := strings.TrimPrefix(filepath, prefix); len(p) < len(filepath) {
		if _, err := b.fs.Open(p); err != nil {
			return false
		}
		return true
	}
	return false
}

func BinaryFileSystem(root string) *binaryFileSystem {
	fs := &assetfs.AssetFS{
		Asset:     asset.Asset,
		AssetDir:  asset.AssetDir,
		AssetInfo: asset.AssetInfo,
		Prefix:    root,
	}
	return &binaryFileSystem{
		fs,
	}
}

// Router 网关api路由
func Router() *gin.Engine {
	router := gin.Default()

	// router.Static("/static/", "./static")
	// 将静态文件打包到bin文件
	router.Use(static.Serve("/static/", BinaryFileSystem("static")))

	// 注册
	router.GET("/user/signup", handler.SigninHandler)
	router.POST("/user/signup", handler.DoSignupHandler)
	// 登陆
	router.GET("/user/signin", handler.SigninHandler)
	router.POST("/user/signin", handler.DoSigninHandler)

	router.Use(handler.Authorize())

	// 用户查询
	router.POST("/user/info", handler.UserInfoHandler)

	// 用不文件查询
	router.POST("/file/query", handler.FileQueryHandler)
	// 用户文件修改（重命名）
	router.POST("/file/update", handler.FileMetaUpdateHandler)

	return router
}
