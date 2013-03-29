package models

import (
	"github.com/robfig/revel"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"time"
)

type Article struct {
	Id        bson.ObjectId `bson:"_id,omitempty"`
	Author_id bson.ObjectId `bson:"Author_id"`
	Published bool          `bson:"Published"`
	Posted    time.Time     `bson: "Posted"`
	Title     string        `bson:"Title"`
	Body      string        `bson:"Body"`
	Tags      []string      `bson:"Tags"`
	Alias     string        `bson:"Alias"`
	Meta      map[string]string
}

// Return the appropriate collection instance for this user.
func (article *Article) Collection(s *mgo.Session) *mgo.Collection {
	return s.DB("bloggo").C("articles")
}

func (article *Article) All(s *mgo.Session) []*Article {
	articles := []*Article{}

	coll := article.Collection(s)
	query := coll.Find(nil).Sort("-Published").Limit(10)
	query.All(&articles)

	for i, a := range articles {
		if a.Meta == nil {
			a.Meta = make(map[string]string)
		}
		articles[i].Meta["author"] = a.GetAuthor(s).String()
	}
	return articles
}

func (article *Article) GetAuthor(s *mgo.Session) *User {
	author := new(User)
	auth := author.GetById(s, article.Author_id)
	return auth
}

func (article *Article) Validate(v *revel.Validation) {
	v.Check(article.Title,
		revel.Required{},
		revel.MinSize{1},
		revel.MaxSize{256},
	)
}

func (article *Article) GetByTitle(s *mgo.Session, Title string) []Article {
	articles := []Article{}

	coll := article.Collection(s)
	query := coll.Find(bson.M{"Title": Title})
	query.One(articles)

	return articles
}

func (article *Article) GetById(s *mgo.Session, Id bson.ObjectId) *Article {
	a := new(Article)
	coll := a.Collection(s)
	query := coll.FindId(Id)
	query.One(a)

	return a
}

func (article *Article) GetByIdString(s *mgo.Session, Id string) *Article {
	ObjectId := bson.ObjectIdHex(Id)
	return article.GetById(s, ObjectId)
}

func (article *Article) Save(s *mgo.Session) error {
	coll := article.Collection(s)
	_, err := coll.Upsert(bson.M{"_id": article.Id}, article)
	if err != nil {
		revel.WARN.Printf("Unable to save user account: %v error %v", article, err)
	}
	return err
}
