package extract

import (
	"github.com/pingcap/parser/ast"
	. "github.com/pingcap/parser/format"
	"strings"
)

type DdlVisitor struct {
	sb         strings.Builder
	ddlColumns []string
}

func (v *DdlVisitor) Enter(in ast.Node) (out ast.Node, skipChildren bool) {
	switch n := in.(type) {
	case *ast.CreateTableStmt:
		for _, col := range n.Cols {
			err := col.Restore(NewRestoreCtx(DefaultRestoreFlags, &v.sb))
			if err != nil {
				panic(err)
			}
			v.ddlColumns = append(v.ddlColumns, v.sb.String())
			v.sb.Reset()
		}
	}
	return in, false
}

func (v *DdlVisitor) Leave(in ast.Node) (out ast.Node, ok bool) {
	return in, true
}

func (v *DdlVisitor) Parser(sql string) []string {
	stmts, _, _ := p.Parse(sql, "", "")
	for _, stmt := range stmts {
		stmt.Accept(v)
	}
	return v.ddlColumns
}

func (v *DdlVisitor) Clear() {
	v.ddlColumns = []string{}
}

func NewDDL() *DdlVisitor {
	return &DdlVisitor{}
}
