package ceph

import (
	"gopkg.in/amz.v1/aws"
	"gopkg.in/amz.v1/s3"
	"netdisk/config"
)

var cephConn *s3.S3

// GetCephConnection 获取ceph连接
func GetCephConnection() *s3.S3 {
	if cephConn != nil {
		return cephConn
	}

	// 1.初始化ceph的一些信息
	auth := aws.Auth{
		AccessKey: config.CephAccessKey,
		SecretKey: config.CephSecretkey,
	}

	curRegion := aws.Region{
		Name:                 "default",
		EC2Endpoint:          config.CephGWEndpoint,
		S3Endpoint:           config.CephGWEndpoint,
		S3BucketEndpoint:     "",
		S3LocationConstraint: false,
		S3LowercaseBucket:    false,
		Sign:                 aws.SignV2,
	}

	// 2.创建S3类型的连接
	return s3.New(auth, curRegion)
}

// GetCephBucket 获取指定的bucket对象
func GetCephBucket(bucket string) *s3.Bucket {
	conn := GetCephConnection()
	return conn.Bucket(bucket)
}

// putObject 上传文件到ceph集群
func putObject(bucket string, path string, data []byte) error {
	return GetCephBucket(bucket).Put(path, data, "octet-stream", s3.PublicRead)
}
