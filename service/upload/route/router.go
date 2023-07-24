package route

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"upload/api"
)

// Router 路由表配置
func Router() *gin.Engine {
	// gin framework 包括Logger, Recovery
	router := gin.Default()

	// 处理静态资源 TODO netdisk底下的相对路径
	router.Static("/static/", "./static")

	// // 加入中间件，用于校验token的拦截器(将会从account微服务中验证)
	// router.Use(handler.HTTPInterceptor())

	// 使用gin插件支持跨域请求
	router.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:  []string{"Origin", "Range", "x-requested-with", "content-type"},
		ExposeHeaders: []string{"Content-Length", "Accept-Ranges", "Content-Range", "Content-Disposition"},
		// AllowCredentials: true,
	}))

	// 文件上传相关接口
	router.POST("/file/upload", api.DoUploadHandler)

	// 秒传接口
	router.POST("/file/fastupload", api.TryFastUploadHandler)

	// 分块上传接口
	router.POST("/file/mpupload/init", api.InitialMultipartUploadHandler)
	router.POST("/file/mpupload/uppart", api.UploadPartHandler)
	router.POST("/file/mpupload/complete", api.CompleteUploadHandler)

	return router
}
