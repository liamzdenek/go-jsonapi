package jsonapie;

import (. "..";"database/sql";"fmt");

func init() {
    // sanity check to ensure this satisfies the interface at compile time
    var c Session = &SessionSimple{};
    _ = c; // compiler stfu about unuse
    
    var crs SessionResourceSQL = &SessionSimple{};
    _ = crs; // compiler stfu abot unuse
}

type SessionSimple struct {
    Transactions map[*sql.DB]*sql.Tx;
}

func NewSessionSimple() *SessionSimple {
    return &SessionSimple{
        Transactions: make(map[*sql.DB]*sql.Tx),
    }
}

func (ctx *SessionSimple) Begin() error {
    fmt.Printf("BEGIN\n");
    return nil;
}

func (ctx *SessionSimple) Success() error {
    fmt.Printf("SUCCESS\n");
    for _,tx := range ctx.Transactions {
        err := tx.Commit();
        if err != nil { // TODO: is this best? should we attempt to roll them all back and send all the errors at once at the end?
            return err;
        }
    }
    return nil;
}

func (ctx *SessionSimple) Failure() error {
    fmt.Printf("FAILURE\n");
    for _,tx := range ctx.Transactions {
        err := tx.Rollback();
        if err != nil { // TODO: is this best? should we attempt to roll them all back and send all the errors at once at the end?
            return err;
        }
    }
    return nil;
}

func (ctx *SessionSimple) GetSQLTransaction(db *sql.DB) (*sql.Tx, error) {
    fmt.Printf("GETSQLTX\n");
    if tx, ok := ctx.Transactions[db]; ok && tx != nil {
        return tx, nil;
    }
    res,err := db.Begin();
    if err != nil {
        return nil, err;
    }
    ctx.Transactions[db] = res;
    return res, nil;
}
