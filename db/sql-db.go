package db

import (
	"database/sql"
	"oap-reposts/env"
)

type SqlDb struct {
	*sql.DB
}

func NewSqlDb() (*SqlDb, error) {
	dsn := env.GetDbDsn()
	db, err := sql.Open("postgres", dsn)
	return &SqlDb{db}, err
}
