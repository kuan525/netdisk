package config

const (
	// TracerAgentHost tracing agent地址
	TracerAgentHost = "localhost:6831"

	// UploadEntry 配置上传入口地址
	UploadEntry = "localhost:28080"
	// UploadServiceHost 上传服务监听地址
	UploadServiceHost = "localhost:28080"

	// UploadLBHost 上传服务LB地址
	UploadLBHost = "http://upload.netdisk.com"
	// DownloadLBHost 下载服务LB地址
	DownloadLBHost = "http://download.netdisk.com"
	
	// DownloadEntry 配置下载入口地址
	DownloadEntry = "localhost:38080"
	// DownloadServiceHost 下载服务监听的地址
	DownloadServiceHost = "localhost:38080"

	// Apigw 监听地址
	// Apigw = "localhost:8080"
)
