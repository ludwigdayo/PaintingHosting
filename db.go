package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var db = new(sql.DB)

// 数据库路径
func openDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	return db, err
}
