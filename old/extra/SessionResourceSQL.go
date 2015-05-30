package jsonapie;

import "database/sql"

type SessionResourceSQL interface {
    GetSQLTransaction(db *sql.DB) (*sql.Tx, error)
}
