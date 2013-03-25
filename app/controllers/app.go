package controllers

import (
	"github.com/jgraham909/bloggo/app/models"
	m "github.com/jgraham909/revmgo/app/controllers"
	"github.com/robfig/revel"
	"labix.org/v2/mgo/bson"
)

type Application struct {
	m.MgoController
	User *models.User
}

func init() {
	revel.InterceptMethod((*Application).SetUser, revel.BEFORE)
	revel.TemplateFuncs["nil"] = func(a interface{}) bool {
		return a == nil
	}
}

func (c Application) Index() revel.Result {
	if c.User != nil {
		user := c.User
		return c.Render(user)
	}
	return c.Render()
}

func (c *Application) SetUser() revel.Result {
	c.User = c.GetCurrentUser()
	return nil
}

func (c Application) UserAuthenticated() bool {
	if _, ok := c.Session["user"]; ok {
		return true
	}
	return false
}

func (c Application) GetCurrentUser() *models.User {
	var u *models.User

	if email, ok := c.Session["user"]; ok {
		u = c.GetUser(email)
	}
	return u
}

func (c Application) GetUser(Email string) *models.User {
	u := models.User{}

	coll := c.MSession.DB("bloggo").C("users")
	query := coll.Find(bson.M{"Email": Email})
	query.One(&u)

	return &u
}
