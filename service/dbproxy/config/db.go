package config

import "fmt"

var (
	// MySQLSource : 要连接的数据库源
	// 其中test：test是用户名密码
	// 127.0.0.1:3306 是ip及端口
	// netdisk 是数据库名
	// charset=utf8 指定里数据以utf8字符编码进行传输
	MySQLSource = "root:123456@tcp(127.0.0.1:3306)/netdisk?charset=utf8"
)

func UpdateDBHost(host string) {
	MySQLSource = fmt.Sprintf("root:123456@tcp(%s)/netdisk?charset=utf8", host)
}
