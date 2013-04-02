package controllers

import (
	"code.google.com/p/go.crypto/bcrypt"
	"fmt"
	"github.com/jgraham909/bloggo/app/models"
	"github.com/robfig/revel"
	"labix.org/v2/mgo/bson"
)

type User struct {
	Application
}

func (c User) Index() revel.Result {
	if c.User != nil {
		action := "/User/SaveExistingUser"
		ObjectId := bson.ObjectId.Hex(c.User.Id)
		return c.Render(action, ObjectId)
	}
	return c.Redirect(User.Login)
}

func (c User) SaveExistingUser(user *models.User, password models.Password, ObjectId string) revel.Result {
	// Weak access control (only let users change their own account)
	if c.User.Id == bson.ObjectIdHex(ObjectId) {
		// Don't trust user submitted id... load from session.
		user.Id = c.User.Id
		user.Validate(c.Validation)

		// Only validate the password if either is non-empty
		if password.Pass != "" || password.PassConfirm != "" {
			user.ValidatePassword(c.Validation, password)
		}

		if c.Validation.HasErrors() {
			c.Validation.Keep()
			c.FlashParams()
			c.Flash.Error("Please correct the errors below.")
			return c.Redirect(User.Index)
		}

		user.Save(c.MongoSession, password)

		// Refresh the session in case the email address was changed.
		c.Session["user"] = user.Email

		c.Flash.Success("Successfully updated account")
		return c.Redirect(Application.Index)
	}
	return c.Forbidden("You can only edit your own account. ")
}

func (c User) SaveNewUser(user *models.User, password models.Password) revel.Result {
	if exists := user.GetByEmail(c.MongoSession, user.Email); exists.Email == user.Email {
		msg := fmt.Sprint("Account with ", user.Email, " already exists.")
		c.Validation.Required(user.Email != exists.Email).
			Message(msg)
	} else {
		user.Id = bson.NewObjectId()
	}

	user.Validate(c.Validation)
	user.ValidatePassword(c.Validation, password)

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		c.Flash.Error("Please correct the errors below.")
		return c.Redirect(User.RegisterForm)
	}

	user.Save(c.MongoSession, password)

	c.Session["user"] = user.Email
	c.Flash.Success("Welcome, " + user.String())
	return c.Redirect(Application.Index)
}

func (c User) Login(Email, Password string) revel.Result {
	user := new(models.User)
	user = user.GetByEmail(c.MongoSession, Email)

	if user.Email != "" {
		err := bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(Password))
		if err == nil {
			c.Session["user"] = Email
			c.Flash.Success("Welcome, " + user.String())
			return c.Redirect(Application.Index)
		}
	}

	c.Flash.Out["mail"] = Email
	c.Flash.Error("Incorrect email address or password.")
	return c.Redirect(User.LoginForm)
}

func (c User) LoginForm() revel.Result {
	if c.UserAuthenticated() == false {
		return c.Render()
	}

	// User already logged in bounce to main page
	return c.Redirect(Application.Index)
}

func (c User) RegisterForm() revel.Result {
	action := "/User/SaveNewUser"
	return c.Render(action)
}

func (c User) Logout() revel.Result {
	for k := range c.Session {
		delete(c.Session, k)
	}
	return c.Redirect(Application.Index)
}
