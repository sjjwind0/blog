package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

const (
	kDatabaseName  = "blog"
	kConnectString = "root:123456@tcp(localhost:3306)/blog?charset=utf8"
)

type Database struct {
	DB  *sql.DB
	ref int
}

var database *Database = nil

func DatabaseInstance() *Database {
	if database == nil {
		database = &Database{}
		database.Open()
	}
	return database
}

func (this *Database) Open() error {
	var err error = nil
	if this.ref == 0 {
		this.DB, err = sql.Open("mysql", kConnectString)
		if err != nil {
			fmt.Printf("connect err", err)
			return err
		}
	}
	fmt.Println("open success")
	this.ref++
	return nil
}

func (this *Database) Close() {
	this.ref--
	if this.ref == 0 {
		this.DB.Close()
	}
}

func (this *Database) DoesTableExist(tableName string) bool {
	rows, err := this.DB.Query("select * from `INFORMATION_SCHEMA`.`TABLES` where table_name =? and TABLE_SCHEMA=?", tableName, kDatabaseName)
	fmt.Println("err = ", err)
	defer rows.Close()
	if err == nil {
		if rows.Next() {
			return true
		}
	}
	return false
}
