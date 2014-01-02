package models

import (
	"github.com/c9s/inflect"
	"github.com/jgraham909/bloggo/app"
	"github.com/robfig/revel"
	"labix.org/v2/mgo"
	"reflect"
	"strings"
)

var (
	Collections map[string]string
)

// Empty struct to embed in models that will provide application default funcs.
type Model struct{}

func Collection(m interface{}, s *mgo.Session) *mgo.Collection {
	typ := reflect.TypeOf(m)
	i := strings.LastIndex(typ.String(), ".") + 1
	n := typ.String()[i:len(typ.String())]
	c := inflect.Tableize(n)

	revel.TRACE.Printf("Using collection %s for %s", c, typ.String())

	return s.DB(app.DB).C(c)
}

// It is expected that each model will embed the type 'Model' and then extend
// or override the following functions to enforce corresponding business rules.
func (m *Model) CanBeCreatedBy(s *mgo.Session, u *User) bool {
	// Default nobody can create.
	return false
}

func (m *Model) CanBeReadBy(s *mgo.Session, u *User) bool {
	// Default everybody can read.
	return true
}

func (m *Model) CanBeDeletedBy(s *mgo.Session, u *User) bool {
	// Default nobody can delete.
	return false
}

func (m *Model) CanBeUpdatedBy(s *mgo.Session, u *User) bool {
	// Default nobody can update.
	return false
}
