package resource;

import "fmt"

type SQLQuery struct {
    Query string;
    FmtArguments []interface{}
    SqlArguments []interface{}
}

func(q *SQLQuery) PrepareQuery() string {
    return fmt.Sprintf(q.Query, q.FmtArguments...);
}

type SQLExpression interface {
    Express(*SQLQuery)
}

type SQLLogic struct {
    Expressions []SQLExpression
    Keyword string
}

func (logic *SQLLogic) Express(q *SQLQuery) {
    if len(logic.Expressions) > 1 {
        q.Query += "(";
    }
    for i,e := range logic.Expressions {
        if i != 0 {
            q.Query += logic.Keyword+" ";
        }
        e.Express(q);
    }
    if len(logic.Expressions) > 1 {
        q.Query += ")";
    }
}

func NewSQLLogic(keyword string, exprs []SQLExpression) *SQLLogic {
    if len(exprs) == 0 {
        return nil;
    }
    return &SQLLogic{
        Expressions: exprs,
        Keyword: keyword,
    }
}

func NewSQLAnd(exprs ...SQLExpression) *SQLLogic {
    return NewSQLLogic("AND", exprs)
}

func NewSQLOr(exprs ...SQLExpression) *SQLLogic {
    return NewSQLLogic("OR", exprs)
}

type SQLWhere struct {
    Expression SQLExpression
}

func NewSQLWhere(expr SQLExpression) *SQLWhere {
    if expr == nil {
        return nil;
    }
    return &SQLWhere{
        Expression: expr,
    }
}

func (where *SQLWhere) Express(q *SQLQuery) {
    if where == nil || where.Expression == nil {
        return;
    }
    q.Query += "WHERE ";
    where.Expression.Express(q);
}

type SQLLiteral struct {
    Expression string
}

func(literal *SQLLiteral) Express(q *SQLQuery) {
    q.Query += literal.Expression+" ";
}

type SQLEquals struct {
    Field string
    Value interface{}
}

func(equals *SQLEquals) Express(q *SQLQuery) {
    q.Query += equals.Field+"=? ";
    q.SqlArguments = append(q.SqlArguments, equals.Value);
}
