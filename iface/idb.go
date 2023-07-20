package iface

import "database/sql"

type Idb interface {
	Query(string, ...any) (*sql.Rows, error)
	Begin() (*sql.Tx, error)
}

// type Itx interface {
// 	Idb
// 	Rollback() error
// 	Commit() error
// }
