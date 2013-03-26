package controllers

import (
	"github.com/jgraham909/bloggo/app/models"
	"github.com/robfig/revel"
	"labix.org/v2/mgo/bson"
	"time"
)

type Blog struct {
	Application
}

func (c Blog) Index() revel.Result {
	articles := []models.Article{}

	coll := c.MSession.DB("bloggo").C("articles")
	query := coll.Find(nil).Sort("-Published").Limit(10)
	query.All(&articles)
	return c.Render(articles)
}

func (c Blog) Add() revel.Result {
	if c.User != nil {
		user := c.User
		article := models.Article{}
		ObjectId := bson.ObjectId.Hex(article.Id)
		action := "/Blog/Create"
		return c.Render(user, action, ObjectId, article)
	}
	return c.Forbidden("You must be logged in to create articles.")
}

func (c Blog) Create(article *models.Article) revel.Result {
	article.Validate(c.Validation)
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		c.Flash.Error("Please correct the errors below.")
		return c.Redirect(Blog.Add)
	}

	// Set calculated fields
	article.Author_id = c.User.Id
	article.Published = true
	article.Posted = time.Now()
	article.Id = bson.NewObjectId()
	article.Save(c.MSession)
	return c.Redirect(Application.Index)
}
