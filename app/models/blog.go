package models

import (
	"github.com/robfig/revel"
	"labix.org/v2/mgo/bson"
	"time"
)

type Blog struct {
	Id        bson.ObjectId `bson:"_id,omitempty"`
	Author_id bson.ObjectId `Author_id`
	Published bool          `Published`
	Posted    time.Time     `Posted`
	Title     string        `Title`
	Body      string        `Body`
	Tags      []string      `Tags`
	Alias     string        `Alias`
}

func (blog *Blog) Validate(v *revel.Validation) {
	v.Check(blog.Title,
		revel.Required{},
		revel.MinSize{1},
		revel.MaxSize{256},
	)
}
