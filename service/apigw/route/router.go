package route

import (
	"apigw/handler"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/kuan525/netdisk/asset"
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

// BinaryFileSystem 返回asset封装的静态资源对象
func BinaryFileSystem(root string) *binaryFileSystem {
	fs := &assetfs.AssetFS{
		Asset:     asset.Asset,     // 获取文件内容
		AssetDir:  asset.AssetDir,  // 获取目录列表
		AssetInfo: asset.AssetInfo, // 获取文件信息
		Prefix:    root,            // 文件系统的根路径
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

	// 用户文件查询
	router.POST("/file/query", handler.FileQueryHandler)
	// 用户文件修改（重命名）
	router.POST("/file/update", handler.FileMetaUpdateHandler)

	return router
}
