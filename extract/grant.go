package extract

import (
	"github.com/pingcap/parser/ast"
)

type GrantVisitor struct {
	db    string
	table string
}

func (v *GrantVisitor) Enter(in ast.Node) (out ast.Node, skipChildren bool) {
	switch n := in.(type) {
	case *ast.GrantStmt:
		v.db = n.Level.DBName
		if n.Level.TableName == "" {
			v.table = "*"
		} else {
			v.table = n.Level.TableName
		}
	}
	return in, false
}

func (v *GrantVisitor) Leave(in ast.Node) (out ast.Node, ok bool) {
	return in, true
}

func (v *GrantVisitor) Parser(sql string) (string, string) {
	stmts, _, _ := p.Parse(sql, "", "")
	for _, stmt := range stmts {
		stmt.Accept(v)
	}
	return v.db, v.table
}

func NewGrant() *GrantVisitor {
	return &GrantVisitor{}
}
