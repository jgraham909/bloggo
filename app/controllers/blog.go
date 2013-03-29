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
	article := new(models.Article)
	articles := article.All(c.MSession)

	return c.Render(articles)
}

func (c Blog) Add() revel.Result {
	if c.User != nil {
		article := models.Article{}
		ObjectId := bson.ObjectId.Hex(article.Id)
		action := "/Blog/Create"
		return c.Render(action, ObjectId, article)
	}
	return c.Forbidden("You must be logged in to create articles.")
}

func (c Blog) Create(article *models.Article) revel.Result {
	if c.User != nil {
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
	}
	return c.Redirect(Application.Index)
}
