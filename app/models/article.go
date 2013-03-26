package models

import (
	"github.com/robfig/revel"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"time"
)

type Article struct {
	Id        bson.ObjectId `bson:"_id,omitempty"`
	Author_id bson.ObjectId `Author_id`
	Published bool          `Published`
	Posted    time.Time     `Posted`
	Title     string        `Title`
	Body      string        `Body`
	Tags      []string      `Tags`
	Alias     string        `Alias`
}

func (article *Article) Validate(v *revel.Validation) {
	v.Check(article.Title,
		revel.Required{},
		revel.MinSize{1},
		revel.MaxSize{256},
	)
}

func (a *Article) GetByTitle(s *mgo.Session, Title string) []Article {
	articles := []Article{}

	coll := s.DB("bloggo").C("blogs")
	query := coll.Find(bson.M{"Title": Title})
	query.One(articles)

	return articles
}

func (a *Article) GetById(s *mgo.Session, Id bson.ObjectId) *Article {
	article := new(Article)
	coll := s.DB("bloggo").C("articles")
	query := coll.FindId(Id)
	query.One(article)

	return article
}

func (a *Article) Save(s *mgo.Session) error {
	coll := s.DB("bloggo").C("articles")
	_, err := coll.Upsert(bson.M{"_id": a.Id}, a)
	if err != nil {
		revel.WARN.Printf("Unable to save user account: %v error %v", a, err)
	}
	return err
}
