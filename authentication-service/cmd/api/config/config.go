package config

import "database/sql"

var Db *sql.DB

func SetDb(_db *sql.DB) {
	Db = _db
}
func GetDb() *sql.DB {
	return Db
}
