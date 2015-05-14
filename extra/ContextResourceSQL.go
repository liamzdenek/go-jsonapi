package jsonapie;

import "database/sql"

type ContextResourceSQL interface {
    GetSQLTransaction(db *sql.DB) (*sql.Tx, error)
}
