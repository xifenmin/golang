package fsyn

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
)

type Mysql struct {
	dbobj  *sql.DB
	dbrows *sql.Rows
}

func NewMysql(ip string, port int, username string, passwd string, dbname string) (mysqlobj *Mysql) {
	var mysql = username + ":" + passwd + "@tcp(" + ip + ":" + strconv.Itoa(port) + ")/" + dbname + "?charset=utf8"
	mysqlobj = new(Mysql)
	mysqlobj.dbobj, _ = sql.Open("mysql", mysql)
	return
}

func (mysqlobj *Mysql) Query(context string) map[string]string {
	mysqlobj.dbrows, _ = mysqlobj.dbobj.Query(context)
	record_map := mysqlobj.Rows()
	return record_map
}

func (mysqlobj *Mysql) Rows() map[string]string {

	var ac int
	var name string
	var passwd string
	var ip string
	var m1, a1, c1 string

	record := make(map[string]string)

	defer mysqlobj.dbrows.Close()

	if mysqlobj.dbrows == nil {
		return record
	}

	for mysqlobj.dbrows.Next() {
		err := mysqlobj.dbrows.Scan(&ac, &name, &passwd, &ip, &m1, &a1, &c1)

		if err != nil {
			fmt.Println(err)
		}

		info := fmt.Sprintf("ac:%d,username:%s,passwd:%s,ip:%s,m1:%s,a1:%s,c1:%s", ac, name, passwd, ip, m1, a1, c1)
		record[strconv.Itoa(ac)] = info
	}

	return record
}
