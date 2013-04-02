package controllers

import (
	"github.com/jgraham909/bloggo/app/models"
	m "github.com/jgraham909/revmgo/app/controllers"
	"github.com/robfig/revel"
)

type Application struct {
	*revel.Controller
	m.MongoController
	User *models.User
}

func init() {
	revel.InterceptMethod((*Application).Setup, revel.BEFORE)
	revel.TemplateFuncs["nil"] = func(a interface{}) bool {
		return a == nil
	}
}

// Responsible for doing any necessary setup for each web request.
func (c *Application) Setup() revel.Result {
	// If there is an active user session load the User data for this user.
	if email, ok := c.Session["user"]; ok {
		c.User = c.User.GetByEmail(c.MongoSession, email)
		c.RenderArgs["user"] = c.User
	}
	return nil
}

func (c Application) Index() revel.Result {
	return c.Redirect(Blog.Index)
}

func (c Application) UserAuthenticated() bool {
	if _, ok := c.Session["user"]; ok {
		return true
	}
	return false
}
