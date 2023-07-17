package config

import (
	cmn "github.com/kuan525/netdisk/common"
)

const (
	// TempLocalRootDir 本地存储的路径
	TempLocalRootDir = "/data/netdisk"
	// TempPartRootDir 分块文件在本地临时存储地址的路径
	TempPartRootDir = "/data/netdisk_part/"
	// CephRootDir Ceph的存储路径
	CephRootDir = "/ceph"
	// COSRootDir COS的存储路径prefix
	COSRootDir = "cos/"
	// CurrentStoreType 设置当前文件的存储类型
	CurrentStoreType = cmn.StoreLocal
)
