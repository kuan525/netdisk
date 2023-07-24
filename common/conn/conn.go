package conn

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	cfg "github.com/kuan525/netdisk/config"
	"log"
	"os"
)

var db *sql.DB

func InitDBConn() {
	db, _ = sql.Open("mysql", cfg.MySQLSource)
	// "database/sql" 自带连接池，这里设置连接最大数量
	db.SetMaxOpenConns(1000)
	err := db.Ping()
	if err != nil {
		fmt.Println("Failed to connect to mysql, err:" + err.Error())
		os.Exit(1)
	}
}

// DBConn 返回数据库连接对象
func DBConn() *sql.DB {
	return db
}

// ParseRows 解析*sql.Rows的函数，转换成[]map[string]interface{}
// 每个map对应一行数据，键是列名，值是对应的值
func ParseRows(rows *sql.Rows) []map[string]interface{} {
	// 获取查询结果的列名
	columns, _ := rows.Columns()

	// 创建用于扫描和存储数据的相关变量
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))

	// 将values中的每个元的指针赋值给scanArgs，用于接受每个列的值
	for j := range values {
		scanArgs[j] = &values[j]
	}

	// 创建一个map用于保存一行数据
	record := make(map[string]interface{})
	//创建一个切片，用于保存所有行的数据
	records := make([]map[string]interface{}, 0)
	for rows.Next() {
		// scanArgs存储的是values的指针，这样将数据写入到values中去，按顺序将数据写入
		err := rows.Scan(scanArgs...)
		checkErr(err)

		for i, col := range values {
			if col != nil {
				record[columns[i]] = col
			}
		}
		// 值拷贝
		records = append(records, record)
	}
	return records
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}
