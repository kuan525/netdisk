package orm

import (
	"database/sql"
	mydb "github.com/kuan525/netdisk/dbclient/conn"
	"log"
)

// OnFileUploadFinished 文件上传完成，保存meta
func OnFileUploadFinished(filehash string, filename string, filesize int64, fileaddr string) (res ExecResult) {
	stmt, err := mydb.DBConn().Prepare(
		"insert ignore into tbl_file " +
			"(`file_sha1`, `file_name`, `file_size`, `file_addr`, `status`) " +
			"values (?,?,?,?,1)")
	if err != nil {
		log.Println("failed to prepare statement, err:" + err.Error())
		res.Suc = false
		return
	}
	defer stmt.Close()

	// Exec : 执行sql
	ret, err := stmt.Exec(filehash, filename, filesize, fileaddr)
	if err != nil {
		log.Println(err.Error())
		res.Suc = false
		return
	}
	// RowsAffected 返回这次sql执行对数据库影响的行数
	if rf, err := ret.RowsAffected(); nil == err {
		if rf <= 0 {
			log.Printf("File with hash:%s has been uploaded before", filehash)
		}
		res.Suc = true
		return
	}
	res.Suc = false
	return
}

// GetFileMeta 从mysql获取文件元信息
func GetFileMeta(filehash string) (res ExecResult) {
	stmt, err := mydb.DBConn().Prepare(
		"select file_sha1, file_addr, file_name, file_size " +
			"from tbl_file " +
			"where file_sha1=? and status=1 limit 1")
	if err != nil {
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	// stmt.QueryRow(filehash) 执行了一个查询操作，并返回查询结果的单行数据
	//.Scan()将查询结果的字段值扫描到相应的变量中
	tfile := TableFile{}
	err = stmt.QueryRow(filehash).Scan(&tfile.FileHash, &tfile.FileAddr, &tfile.FileName, &tfile.FileSize)
	if err != nil {
		if err != sql.ErrNoRows {
			// 查不到对应记录，返回参数以及错误均为nil
			res.Suc = true
			res.Data = nil
			return
		} else {
			log.Println(err.Error())
			res.Suc = false
			res.Msg = err.Error()
			return
		}
	}
	res.Suc = true
	res.Data = tfile
	return
}

// GetFileMetaList 从mysql中批量获取文件元信息
func GetFileMetaList(limit int64) (res ExecResult) {
	stmt, err := mydb.DBConn().Prepare(
		"select file_sha1, file_addr, file_name, file_size " +
			"from tbl_file " +
			"where status=1 limit ?")
	if err != nil {
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(limit)
	if err != nil {
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}

	// 将rows的数据写入tfiles中
	// 获取列名，同时得到多少列
	cloumns, _ := rows.Columns()
	values := make([]sql.RawBytes, len(cloumns))
	var tfiles []TableFile
	// 用了rows.Next()一次之后到达第一行
	for i := 0; i < len(values) && rows.Next(); i++ {
		tfile := TableFile{}
		err = rows.Scan(&tfile.FileHash, &tfile.FileAddr, &tfile.FileName, &tfile.FileSize)
		if err != nil {
			log.Println(err.Error())
			break
		}
		tfiles = append(tfiles, tfile)
	}
	res.Suc = true
	res.Data = tfiles
	return
}

// UpdateFileLocation 更新文件的存储地址（比如文件被转移了）
func UpdateFileLocation(filehash string, fileaddr string) (res ExecResult) {
	stmt, err := mydb.DBConn().Prepare(
		"update tbl_file set `file_addr`=?  where `file_sha1`=? limit 1;")
	if err != nil {
		log.Println("预编译sql失败，err：" + err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	ret, err := stmt.Exec(fileaddr, filehash)
	if err != nil {
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}

	if rf, err := ret.RowsAffected(); nil == err {
		if rf <= 0 {
			log.Printf("更新文件location失败，filehash:%s", filehash)
			res.Suc = false
			res.Msg = "无记录更新"
			return
		}
		res.Suc = true
		return
	} else {
		res.Suc = false
		res.Msg = err.Error()
		return
	}
}
