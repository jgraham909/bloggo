package app

import "github.com/robfig/revel"
import "fmt"

var (
	Db string
)

func AppInit() {
	var found bool
	if Db, found = revel.Config.String("bloggo.db"); !found {
		Db = "bloggo"
	}
}
