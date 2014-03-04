package models

import (
	"github.com/revel/revel"
	"github.com/russross/blackfriday"
	"html/template"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"time"
)

const trimLength = 300

type Article struct {
	Model     `bson:",inline"`
	Author_id bson.ObjectId          `bson:"Author_id"`
	Published bool                   `bson:"Published"`
	Posted    time.Time              `bson:"Posted"`
	Title     string                 `bson:"Title"`
	Body      string                 `bson:"Body"`
	Tags      []string               `bson:"Tags"`
	Alias     string                 `bson:"Alias"`
	Meta      map[string]interface{} `bson:",omitempty"`
}

func (article *Article) AddMeta(s *mgo.Session) {
	if article.Meta == nil {
		article.Meta = make(map[string]interface{})
	}
	auth := article.GetAuthor(s)
	article.Meta["author"] = auth
	article.Meta["author_id"] = auth.Id.Hex()
	article.Meta["markdown"] = template.HTML(string(blackfriday.MarkdownBasic([]byte(article.Body))))

	if len(article.Body) > trimLength {
		article.Meta["teaser"] = template.HTML(string(blackfriday.MarkdownBasic([]byte(article.Body[0:trimLength]))))
	} else {
		article.Meta["teaser"] = template.HTML(string(blackfriday.MarkdownBasic([]byte(article.Body[0:len(article.Body)]))))
	}
}

func (article *Article) GetAuthor(s *mgo.Session) *User {
	auth := GetUserByObjectId(s, article.Author_id)
	return auth
}

func (article *Article) Validate(v *revel.Validation) {
	v.Check(article.Title,
		revel.Required{},
		revel.MinSize{1},
		revel.MaxSize{256},
	)
}

func GetArticlesByTitle(s *mgo.Session, Title string) []Article {
	articles := []Article{}
	article := new(Article)
	query := Collection(article, s).Find(bson.M{"Title": Title})
	query.All(articles)

	return articles
}

func GetArticleByObjectId(s *mgo.Session, Id bson.ObjectId) *Article {
	a := new(Article)

	query := Collection(a, s).FindId(Id)
	query.One(a)
	a.AddMeta(s)
	return a
}

func GetArticleById(s *mgo.Session, Id string) *Article {
	ObjectId := bson.ObjectIdHex(Id)
	return GetArticleByObjectId(s, ObjectId)
}

func GetArticlesByDate(s *mgo.Session, limit int) []*Article {
	articles := []*Article{}
	a := new(Article)
	query := Collection(a, s).Find(bson.M{"Published": true}).Sort("-Posted").Limit(limit)
	query.All(&articles)

	for _, a := range articles {
		a.AddMeta(s)
	}
	return articles
}

func GetArticlesByTag(s *mgo.Session, t string) []*Article {
	articles := []*Article{}
	a := new(Article)
	query := Collection(a, s).Find(bson.M{"Published": true, "Tags": t}).Sort("-Posted")
	query.All(&articles)

	for _, a := range articles {
		a.AddMeta(s)
	}
	return articles
}

func (article *Article) Save(s *mgo.Session) error {
	coll := Collection(article, s)
	_, err := coll.Upsert(bson.M{"_id": article.Id}, article)
	if err != nil {
		revel.WARN.Printf("Unable to save article: %v error %v", article, err)
	}
	return err
}

func (article *Article) Delete(s *mgo.Session) error {
	coll := Collection(article, s)
	err := coll.RemoveId(article.Id)
	if err != nil {
		revel.WARN.Printf("Undable to delete article: %v error %v", article, err)
	}
	return err
}

func (article *Article) CanBeUpdatedBy(s *mgo.Session, u *User) bool {
	if u.Id == article.Author_id {
		return true
	}
	return false
}

func (article *Article) CanBeDeletedBy(s *mgo.Session, u *User) bool {
	return article.CanBeUpdatedBy(s, u)
}

func (article *Article) CanBeCreatedBy(s *mgo.Session, u *User) bool {
	// All valid users can create articles
	if u != nil {
		return u.Id.Valid()
	}
	return false
}
