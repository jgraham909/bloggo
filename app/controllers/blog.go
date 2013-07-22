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

func (c Blog) Tag(tag string) revel.Result {
	articles := models.GetArticlesByTag(c.MongoSession, tag)
	return c.Render(articles, tag)
}

func (c Blog) GetDelete(id bson.ObjectId) revel.Result {
	if c.ActiveUser != nil {
		article := models.GetArticleByObjectId(c.MongoSession, id)
		if article.CanBeDeletedBy(c.MongoSession, c.ActiveUser) {
			return c.Render(article)
		}
		return c.Forbidden("You do not have permission to delete this resource.")
	}
	return c.Forbidden("Only site users may delete content.")
}

func (c Blog) Delete(id bson.ObjectId) revel.Result {
	if c.ActiveUser != nil {
		article := models.GetArticleByObjectId(c.MongoSession, id)
		if article.CanBeDeletedBy(c.MongoSession, c.ActiveUser) {
			article.Delete(c.MongoSession)
		}
	}
	return c.RenderText("")
}

func (c Blog) Links(id bson.ObjectId) revel.Result {
	links := false
	canUpdate := false
	canDelete := false
	article := &models.Article{}
	if c.ActiveUser != nil {
		article = models.GetArticleByObjectId(c.MongoSession, id)
		if article.CanBeUpdatedBy(c.MongoSession, c.ActiveUser) {
			canUpdate = true
			links = true
		}
		if article.CanBeDeletedBy(c.MongoSession, c.ActiveUser) {
			canDelete = true
			links = true
		}
	}
	return c.Render(article, links, canUpdate, canDelete)
}

func (c Blog) GetCreate() revel.Result {
	article := models.Article{}
	if article.CanBeCreatedBy(c.MongoSession, c.ActiveUser) {
		action := "/blog/create"
		actionButton := "Create"
		return c.Render(action, article, actionButton)
	}
	return c.Forbidden("You are not allowed to create articles.")
}

func (c Blog) PostCreate(article *models.Article) revel.Result {
	if c.ActiveUser != nil {
		article.Tags = strings.Split(c.Params.Values["article.Tags"][0], ",")
		article.Validate(c.Validation)
		if c.Validation.HasErrors() {
			c.Validation.Keep()
			c.FlashParams()
			c.Flash.Error("Please correct the errors below.")
			return c.Redirect(Blog.GetCreate)
		}

		// Set calculated fields
		article.Author_id = c.ActiveUser.Id
		article.Published = true
		article.Posted = time.Now()
		article.Id = bson.NewObjectId()
		article.Save(c.MongoSession)
	}
	return c.Redirect(Application.Index)
}

func (c Blog) GetRead(id bson.ObjectId) revel.Result {
	if id.Hex() != "" {
		article := models.GetArticleByObjectId(c.MongoSession, id)
		return c.Render(article)
	}
	return c.NotFound("Invalid article Id.")
}

func (c Blog) GetUpdate(id bson.ObjectId) revel.Result {
	if c.ActiveUser != nil {
		article := models.GetArticleByObjectId(c.MongoSession, id)
		if article.CanBeUpdatedBy(c.MongoSession, c.ActiveUser) {
			action := "/Blog/Update"
			actionButton := "Update"
			return c.Render(action, article, actionButton)
		}
		return c.Forbidden("You do not have permission to edit this resource.")
	}
	return c.Redirect(User.GetLogin)
}

func (c Blog) Update(article *models.Article) revel.Result {
	if c.ActiveUser != nil {
		check := models.GetArticleByObjectId(c.MongoSession, article.Id)
		if check.CanBeUpdatedBy(c.MongoSession, c.ActiveUser) {
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
	return c.Redirect(User.GetLogin)
}
