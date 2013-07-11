package models

import (
	"github.com/jgraham909/bloggo/app"
	"github.com/robfig/revel"
	"labix.org/v2/mgo"
	"reflect"
	"strings"
)

var (
	Collections map[string]string
)

func Collection(m interface{}, s *mgo.Session) *mgo.Collection {
	typ := reflect.TypeOf(m)
	i := strings.LastIndex(typ.String(), ".") + 1
	n := typ.String()[i:len(typ.String())]

	var found bool
	var c string
	if c, found = revel.Config.String("bloggo.db.collection." + n); !found {
		c = n
	}
	return s.DB(app.DB).C(c)
}
