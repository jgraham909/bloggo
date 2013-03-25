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
	// TODO Remove this when debugging isn't needed
	// keep go from complaining about unused fmt
	fmt.Printf("test")

	if _, ok := c.Session["user"]; ok {
		return true
	}
	return false
}

func (c User) getCurrentUser() models.User {
	u := models.User{}

	if email, ok := c.Session["user"]; ok {
		u = c.getUser(email)
	}
	return u
}

func (c User) getUser(Email string) models.User {
	u := models.User{}

	coll := c.MSession.DB("bloggo").C("users")
	query := coll.Find(bson.M{"Email": Email})
	query.One(&u)

	return u
}

func (c User) Index() revel.Result {
	if c.connected() != false {
		return c.Redirect(User.Login)
	}

	// Todo change to redirect to view/edit user account
	c.Flash.Error("Please log in first")
	return c.Render()
}

func (c User) SaveNewUser(user *models.User, verifyPassword string) revel.Result {
	fmt.Printf("SNU %v \n", user.Email)

	if exists := c.getUser(user.Email); exists.Email == user.Email {
		msg := fmt.Sprint("Account with ", user.Email, " already exists.")
		c.Validation.Required(user.Email != exists.Email).
			Message(msg)
	} else {
		user.Id = bson.NewObjectId()
	}
	return c.SaveUser(user, verifyPassword)
}

func (c User) SaveUser(user *models.User, verifyPassword string) revel.Result {
	fmt.Printf("SaveUser(user): %v\n", user)
	c.Validation.Required(verifyPassword)
	c.Validation.Required(verifyPassword == user.Password).
		Message("Password does not match")

	user.Validate(c.Validation)

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		c.Flash.Error("Please correct the errors below.")
		return c.Redirect(User.RegisterForm)
	}

	user.HashedPassword, _ = bcrypt.GenerateFromPassword(
		[]byte(user.Password), bcrypt.DefaultCost)

	// Empty out the unhashed password to ensure it is not stored
	user.Password = ""

	coll := c.MSession.DB("bloggo").C("users")
	_, err := coll.Upsert(bson.M{"_id": user.Id}, user)
	if err != nil {
		revel.WARN.Printf("Unable to save user account: %v error", user, err)
		c.Flash.Error("Unable to save user account, please contact the site administrator")
	}

	c.Session["user"] = user.Email
	c.Flash.Success("Welcome, " + user.String())
	return c.Redirect(Application.Index)
}

func (c User) Login(Email, Password string) revel.Result {
	user := c.getUser(Email)

	if user.Email != "" {
		err := bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(Password))
		if err == nil {
			c.Session["user"] = Email
			c.Flash.Success("Welcome, " + Email)
			return c.Redirect(Application.Index)
		}
	}

	c.Flash.Out["mail"] = Email
	c.Flash.Error("Incorrect email address or password.")
	return c.Redirect(User.LoginForm)
}

func (c User) LoginForm() revel.Result {
	if c.connected() == false {
		return c.Render()
	}

	// User already logged in bounce to main page
	return c.Redirect(Application.Index)
}

func (c User) RegisterForm() revel.Result {
	return c.Render()
}

func (c User) Logout() revel.Result {
	for k := range c.Session {
		delete(c.Session, k)
	}
	return c.Redirect(Application.Index)
}
