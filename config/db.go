package config

const (
	host   = "localhost"                       // 数据库地址
	port   = "3306"                            // 数据库端口
	user   = "root"                            // 数据库用户名
	pwd    = "123456"                          // 数据库密码
	dbname = "netdisk"                         // 数据库名
	format = "?charset=utf8mb4&parseTime=true" // 编码格式
)

const MySQLSource = user + ":" + pwd + "@tcp(" + host + ":" + port + ")/" + dbname + format
