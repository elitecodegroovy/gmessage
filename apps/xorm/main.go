package main

import (
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	"os"
)

var engine *xorm.Engine

func NewEngine() {

	var err error
	engine, err = xorm.NewEngine("mysql", "swisse:swisse@10.50.115.114:16052/swisse?charset=utf8")
	if err != nil {
		panic("can't create the mysql engine")
	}
	engine.ShowSQL(true)
	engine.Logger().SetLevel(core.LOG_DEBUG)

	f, err := os.Create("sql.log")
	if err != nil {
		println(err.Error())
		return
	}
	engine.SetLogger(xorm.NewSimpleLogger(f))
}
func main() {

}
