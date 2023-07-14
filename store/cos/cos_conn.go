package cos

import (
	"context"
	"github.com/tencentyun/cos-go-sdk-v5"
	"log"
	"net/http"
	"net/url"
	"netdisk/config"
	"time"
)

var cosCli *cos.Client

// NewClient 创建cos client对象
func NewClient() *cos.Client {
	u, _ := url.Parse(config.BucketURL)
	su, _ := url.Parse(config.ServiceURL)
	b := &cos.BaseURL{
		BucketURL:  u,
		ServiceURL: su,
	}

	cosCli := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  config.COSSecretID,
			SecretKey: config.COSSecretKey,
		},
	})
	return cosCli
}

// DownloadURL : 临时授权下载url
func DownloadURL(objName string) string {
	ctx := context.Background()
	presignedURL, err := NewClient().Object.GetPresignedURL(
		ctx, http.MethodPut, objName, config.COSSecretID, config.COSSecretKey, time.Hour, nil)
	if err != nil {
		log.Printf(err.Error(), "获取临时授权下载url失败")
	}
	return presignedURL.String()
}

// BuildLifeCycleRule 针对指定bucket设置生命周期规则
func BuildLifeCycleRule(bucketName string) {
	lc := &cos.BucketPutLifecycleOptions{
		Rules: []cos.BucketLifecycleRule{
			{
				ID:     "1",
				Filter: &cos.BucketLifecycleFilter{Prefix: bucketName},
				Status: "Enabled",
				Transition: []cos.BucketLifecycleTransition{
					{
						Days:         30,
						StorageClass: "Standard",
					},
				},
			},
		},
	}
	_, err := NewClient().Bucket.PutLifecycle(context.Background(), lc)
	if err != nil {
		log.Printf(err.Error(), "设置生命周期失败")
	}
}
