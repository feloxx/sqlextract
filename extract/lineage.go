package extract

import (
	"github.com/pingcap/parser/ast"
	. "github.com/pingcap/parser/format"
	"strings"
)

type LineageVisitor struct {
	sb              strings.Builder
	InputTableList  []string
	OutputTableList []string
}

func (v *LineageVisitor) Enter(in ast.Node) (out ast.Node, skipChildren bool) {
	switch n := in.(type) {
	case *ast.CreateTableStmt:
		v.OutputTableList = append(v.OutputTableList, n.Table.Text())
	case *ast.SelectStmt:
		if n.From != nil {
			v.traverseJoin(n.From.TableRefs)
		}
	case *ast.UnionStmt:
		for _, uni := range n.SelectList.Selects {
			if uni.From != nil {
				v.traverseJoin(uni.From.TableRefs)
			}
		}
	case *ast.InsertStmt:
		_ = n.Table.TableRefs.Restore(NewRestoreCtx(DefaultRestoreFlags, &v.sb))
		v.OutputTableList = append(v.OutputTableList, v.sb.String())
		v.sb.Reset()
		if n.Select != nil {
			switch sel := n.Select.(type) {
			case *ast.SelectStmt:
				if sel.From != nil {
					v.traverseJoin(sel.From.TableRefs)
				}
			case *ast.UnionStmt:
				for _, uni := range sel.SelectList.Selects {
					if uni.From != nil {
						v.traverseJoin(uni.From.TableRefs)
					}
				}
			}
		}
	}
	return in, false
}

func (v *LineageVisitor) Leave(in ast.Node) (out ast.Node, ok bool) {
	return in, true
}

func (v *LineageVisitor) traverseJoin(n *ast.Join) {
	switch n.Right.(type) {
	case *ast.TableSource:
		r := n.Right.(*ast.TableSource)
		switch sel := r.Source.(type) {
		case *ast.SelectStmt: //右节点是个子查询
			if sel.From != nil {
				v.traverseJoin(sel.From.TableRefs)
			}
		default: //右节点是个表
			_ = r.Source.Restore(NewRestoreCtx(DefaultRestoreFlags, &v.sb))
			dbTable := strings.Split(v.sb.String(), ".")
			v.sb.Reset()

			if len(dbTable) == 2 && n.On != nil {
				switch n.On.Expr.(type) {
				case *ast.BinaryOperationExpr:
					db := dbTable[0]
					table := dbTable[1]
					v.InputTableList = append(v.InputTableList, db+"."+table)
				}
			}
		}
	}
	switch n.Left.(type) {
	case *ast.TableSource:
		l := n.Left.(*ast.TableSource)
		switch source := l.Source.(type) {
		case *ast.SelectStmt:
			if source.From != nil {
				v.traverseJoin(source.From.TableRefs)
			}
		default:
			_ = l.Source.Restore(NewRestoreCtx(DefaultRestoreFlags, &v.sb))
			dbTable := strings.Split(v.sb.String(), ".")
			v.sb.Reset()

			if len(dbTable) == 2 {
				v.InputTableList = append(v.InputTableList, dbTable[0]+"."+dbTable[1])
			}
		}
	case *ast.Join:
		v.traverseJoin(n.Left.(*ast.Join))
	}
}

func (v *LineageVisitor) Parser(sql string) ([]string, []string) {
	stmts, _, _ := p.Parse(sql, "", "")
	for _, stmt := range stmts {
		stmt.Accept(v)
	}
	return v.InputTableList, v.OutputTableList
}

func NewLineage() *LineageVisitor {
	return &LineageVisitor{}
}
