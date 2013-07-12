package models

import (
	"github.com/robfig/revel"
	"github.com/russross/blackfriday"
	"html/template"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"strings"
	"time"
)

const trimLength = 300

type Article struct {
	Id        bson.ObjectId `bson:"_id,omitempty"`
	Author_id bson.ObjectId `bson:"Author_id"`
	Published bool          `bson:"Published"`
	Posted    time.Time     `bson: "Posted"`
	Title     string        `bson:"Title"`
	Body      string        `bson:"Body"`
	Tags      string        `bson:"Tags"`
	Alias     string        `bson:"Alias"`
	Meta      map[string]interface{}
}

func (article *Article) All(s *mgo.Session) []*Article {
	articles := []*Article{}

	query := Collection(article, s).Find(nil).Sort("-Published").Limit(10)
	query.All(&articles)

	for _, a := range articles {
		a.AddMeta(s)
	}
	return articles
}

func (article *Article) AddMeta(s *mgo.Session) {
	if article.Meta == nil {
		article.Meta = make(map[string]interface{})
	}
	article.Meta["author"] = article.GetAuthor(s).String()
	article.Meta["markdown"] = template.HTML(string(blackfriday.MarkdownBasic([]byte(article.Body))))

	if len(article.Body) > trimLength {
		article.Meta["teaser"] = template.HTML(string(blackfriday.MarkdownBasic([]byte(article.Body[0:trimLength]))))
	} else {
		article.Meta["teaser"] = template.HTML(string(blackfriday.MarkdownBasic([]byte(article.Body[0:len(article.Body)]))))
	}
	article.Meta["id"] = article.Id.Hex()
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

	query := Collection(article, s).Find(bson.M{"Title": Title})
	query.One(articles)

	return articles
}

func (article *Article) GetById(s *mgo.Session, Id bson.ObjectId) *Article {
	a := new(Article)

	query := Collection(article, s).FindId(Id)
	query.One(a)
	a.AddMeta(s)
	return a
}

func (article *Article) GetByIdString(s *mgo.Session, Id string) *Article {
	ObjectId := bson.ObjectIdHex(Id)
	return article.GetById(s, ObjectId)
}

func (article *Article) preSave() {

}

func (article *Article) Save(s *mgo.Session) error {
	coll := Collection(article, s)
	_, err := coll.Upsert(bson.M{"_id": article.Id}, article)
	if err != nil {
		revel.WARN.Printf("Unable to save article: %v error %v", article, err)
	}
	return err
}

func (article *Article) SplitTags() []string {
	return strings.Split(article.Tags, ",")
}
