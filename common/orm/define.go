package orm

import "database/sql"

// TableFile 文件表结构体
type TableFile struct {
	FileHash string
	// 多了一个判断是否有效的情况
	// 默认就部位空的时候，可以使用string，但是可能为空的时候最好使用sql.NullString
	FileName sql.NullString
	FileSize sql.NullInt64
	FileAddr sql.NullString
}

// TableUser 用户表model
type TableUser struct {
	Username     string
	Email        string
	Phone        string
	SignupAt     string
	LastActiveAt string
	Status       int
}

// TableUserFile 用户文件表结构
type TableUserFile struct {
	UserName    string
	FileHash    string
	FileName    string
	FileSize    int64
	UploadAt    string
	LastUpdated string
}

// ExecResult sql函数执行的结果
type ExecResult struct {
	Suc  bool        `json:"suc"`
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}
