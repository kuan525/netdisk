package common

type StoreType int

// 文件存储位置
const (
	_ StoreType = iota
	// StoreLocal 本地节点
	StoreLocal
	// StoreCeph ceph集群
	StoreCeph
	// StoreCOS 腾讯COS
	StoreCOS
	// StoreMix 混合（ceph和COS）
	StoreMix
	// StoreAll 所有类型的存储都能存一份数据
	StoreAll
)
