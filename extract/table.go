package extract

import (
	"github.com/pingcap/parser/ast"
	_ "github.com/pingcap/tidb/types/parser_driver"
	"strings"
)

type tableVisitor struct {
	sb        strings.Builder
	tableList []string
}

func (v *tableVisitor) Enter(in ast.Node) (out ast.Node, skipChildren bool) {
	switch n := in.(type) {
	case *ast.TableName:
		v.tableList = append(v.tableList, n.Schema.O+"."+n.Name.O)
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
