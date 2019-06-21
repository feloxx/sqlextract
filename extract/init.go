package extract

import "github.com/pingcap/parser"

var p *parser.Parser

func init() {
	p = parser.New()
}
