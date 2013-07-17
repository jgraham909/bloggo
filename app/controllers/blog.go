package controllers

import (
	"github.com/jgraham909/bloggo/app/models"
	"github.com/robfig/revel"
	"labix.org/v2/mgo/bson"
	"strings"
	"time"
)

type Blog struct {
	Application
}

func (c Blog) Index() revel.Result {
	articles := models.GetArticlesByDate(c.MongoSession, 10)
	return c.Render(articles)
}

func (c Blog) Tag(t string) revel.Result {
	articles := models.GetArticlesByTag(c.MongoSession, t)
	return c.Render(articles)
}

func (c Blog) EditLinks(id bson.ObjectId) revel.Result {
	canEdit := false
	article := &models.Article{}
	if c.User != nil {
		article = models.GetArticleByObjectId(c.MongoSession, id)
		if article.CanEdit(c.User) {
			canEdit = true
		}
	}
	return c.Render(canEdit, article)
}

func (c Blog) Add() revel.Result {
	if c.User != nil {
		article := models.Article{}
		action := "/Blog/Create"
		actionButton := "Create"
		return c.Render(action, article, actionButton)
	}
	return c.Forbidden("You must be logged in to create articles.")
}

func (c Blog) Create(article *models.Article) revel.Result {
	if c.User != nil {
		article.Tags = strings.Split(c.Params.Values["article.Tags"][0], ",")
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
		article.Save(c.MongoSession)
	}
	return c.Redirect(Application.Index)
}

func (c Blog) View(id bson.ObjectId) revel.Result {
	article := models.GetArticleByObjectId(c.MongoSession, id)
	return c.Render(article)
}

func (c Blog) Edit(id bson.ObjectId) revel.Result {
	if c.User != nil {
		article := models.GetArticleByObjectId(c.MongoSession, id)
		if article.CanEdit(c.User) {
			action := "/Blog/Update"
			actionButton := "Update"
			return c.Render(action, article, actionButton)
		}
		return c.Forbidden("You do not have permission to edit this resource.")
	}
	return c.Redirect(User.Login)
}

func (c Blog) Update(article *models.Article) revel.Result {
	if c.User != nil {
		check := models.GetArticleByObjectId(c.MongoSession, article.Id)
		if check.CanEdit(c.User) {
			article.Tags = strings.Split(c.Params.Values["article.Tags"][0], ",")
			article.Validate(c.Validation)
			if c.Validation.HasErrors() {
				c.Validation.Keep()
				c.FlashParams()
				c.Flash.Error("Please correct the errors below.")
				return c.Redirect("/blog/%s/edit", article.Id.Hex())
			}

			// @todo properly implement this
			article.Author_id = check.Author_id
			article.Published = check.Published
			article.Posted = check.Posted

			article.Save(c.MongoSession)
			return c.Redirect("/blog/%s", article.Id.Hex())
		}
		return c.Forbidden("You do not have permission to edit this resource.")
	}
	return c.Redirect(User.Login)
}
