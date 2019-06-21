package main

import (
	"C"
	"github.com/feloxx/sqlextract/extract"
	"strings"
)

//export GetTableList
func GetTableList(sql *C.char) *C.char {
	sqlStr := C.GoString(sql)
	tableList := extract.NewTable().Parser(strings.TrimSpace(sqlStr))
	result := strings.Join(tableList, "|")
	return C.CString(result)
}

//export GetSqlLineageInput
func GetSqlLineageInput(sql *C.char) *C.char {
	sqlStr := C.GoString(sql)
	inputTableList, _ := extract.NewLineage().Parser(strings.TrimSpace(sqlStr))
	result := strings.Join(inputTableList, "|")
	return C.CString(result)
}

//export GetSqlLineageOutput
func GetSqlLineageOutput(sql *C.char) *C.char {
	sqlStr := C.GoString(sql)
	_, outputTableList := extract.NewLineage().Parser(strings.TrimSpace(sqlStr))
	result := strings.Join(outputTableList, "|")
	return C.CString(result)
}

func main() {

}
