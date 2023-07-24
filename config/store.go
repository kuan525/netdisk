package config

import (
	cmn "github.com/kuan525/netdisk/common"
)

const (
	// TempLocalRootDir 本地存储的路径
	TempLocalRootDir = "/Users/kuan525/TODO/netdisk_storage/"
	// TempPartRootDir 分块文件在本地临时存储地址的路径
	TempPartRootDir = "/Users/kuan525/TODO/netdisk_part_storage/"
	// CephRootDir Ceph的存储路径
	CephRootDir = "/Users/kuan525/TODO/ceph/"
	// COSRootDir COS的存储路径prefix
	COSRootDir = "cos/"
	// CurrentStoreType 设置当前文件的存储类型
	CurrentStoreType = cmn.StoreLocal
)
