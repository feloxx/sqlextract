package extract

import (
	"github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"
	. "github.com/pingcap/parser/format"
	"strings"
)

var p *parser.Parser

func init() {
	p = parser.New()
}

type tableVisitor struct {
	sb        strings.Builder
	tableList []string
}

func (v *tableVisitor) Enter(in ast.Node) (out ast.Node, skipChildren bool) {
	switch n := in.(type) {
	case *ast.TableName:
		_ = n.Restore(NewRestoreCtx(DefaultRestoreFlags, &v.sb))
		v.tableList = append(v.tableList, v.sb.String())
	}
	return in, false
}

func (v *tableVisitor) Leave(in ast.Node) (out ast.Node, ok bool) {
	return in, true
}

func (v *tableVisitor) Parser(sql string) []string {
	stmts, _, _ := p.Parse(sql, "", "")
	for _, stmt := range stmts {
		stmt.Accept(v)
	}
	return v.tableList
}

func NewTable() *tableVisitor {
	return &tableVisitor{}
}
