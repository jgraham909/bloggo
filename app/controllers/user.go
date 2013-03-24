package controllers

import (
	"code.google.com/p/go.crypto/bcrypt"
	"fmt"
	"github.com/jgraham909/bloggo/app/models"
	m "github.com/jgraham909/revmgo/app/controllers"
	"github.com/robfig/revel"
	"labix.org/v2/mgo/bson"
)

type User struct {
	m.MgoController
	models.User
}

func (c User) connected() bool {
	if _, ok := c.Session["user"]; ok {
		return true
	}
	return false
}

func (c User) getUser(Email string) (models.User, error) {
	u := models.User{}

	var err error
	coll := c.MSession.DB("bloggo").C("users")
	query := coll.Find(bson.M{"Email": Email})
	query.One(&u)

	return u, err
}

func (c User) Index() revel.Result {
	if c.connected() != false {
		return c.Redirect(User.Login)
	}
	c.Flash.Error("Please log in first")
	return c.Render()
}

func (c User) SaveUser(user models.User, verifyPassword string) revel.Result {
	c.Validation.Required(verifyPassword)
	c.Validation.Required(verifyPassword == user.Password).
		Message("Password does not match")
	user.Validate(c.Validation)

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(Application.Index)
	}

	user.HashedPassword, _ = bcrypt.GenerateFromPassword(
		[]byte(user.Password), bcrypt.DefaultCost)

	c.Session["user"] = user.Email
	c.Flash.Success("Welcome, " + user.String())
	return c.Redirect(Application.Index)
}

func (c User) Login(Email, Password string) revel.Result {
	user, _ := c.getUser(Email)

	if user.Email != "" {
		err := bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(Password))
		if err == nil {
			c.Session["user"] = Email
			c.Flash.Success("Welcome, " + Email)
			return c.Redirect(Application.Index)
		}
	}

	c.Flash.Out["email"] = Email
	c.Flash.Error("Incorrect email address or password.")
	return c.Redirect(User.LoginForm)
}

func (c User) LoginForm() revel.Result {
	derp, _ := bcrypt.GenerateFromPassword([]byte("1234"), bcrypt.DefaultCost)
	fmt.Printf("hashed: %v\n", derp)
	if c.connected() == false {
		return c.Render(User.LoginForm)
	}
	return c.Redirect(Application.Index)
}

func (c User) RegisterForm() revel.Result {
	return c.Render(User.RegisterForm)
}

func (c User) Register() revel.Result {
	return c.Todo()
}

func (c User) Logout() revel.Result {
	for k := range c.Session {
		delete(c.Session, k)
	}
	return c.Redirect(Application.Index)
}
