package main

import "C"
import (
	"github.com/pingcap/parser"
)

var p *parser.Parser

func init() {
	p = parser.New()
}

//export GetTableList
func GetTableList(sql *C.char) *C.char {
	return C.CString("")
}

func main() {
	
}
