package controllers

import (
	"code.google.com/p/go.crypto/bcrypt"
	"fmt"
	"github.com/jgraham909/bloggo/app/models"
	"github.com/revel/revel"
	"labix.org/v2/mgo/bson"
)

type User struct {
	Application
}

func (c User) GetUpdate(id bson.ObjectId) revel.Result {
	if c.ActiveUser != nil {
		action := "/User/" + id.Hex()
		user := c.ActiveUser
		if user.CanBeUpdatedBy(c.MongoSession, c.ActiveUser) {
			return c.Render(action, user)
		}
		return c.Forbidden("You are not allowed to edit this resource.")
	}
	return c.Redirect(User.GetLogin)
}

func (c User) PostUpdate(id bson.ObjectId, user *models.User, password models.Password) revel.Result {
	if user.CanBeUpdatedBy(c.MongoSession, c.ActiveUser) {
		// Don't trust user submitted id... load from session.
		user.Id = c.ActiveUser.Id
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

func (c User) PostCreate(user *models.User, password models.Password) revel.Result {
	if user.CanBeCreatedBy(c.MongoSession, c.ActiveUser) {
		if exists := models.GetUserByEmail(c.MongoSession, user.Email); exists.Email == user.Email {
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
			return c.Redirect(User.GetCreate)
		}

		user.Save(c.MongoSession, password)

		c.Session["user"] = user.Email
		c.Flash.Success("Welcome, " + user.String())
		return c.Redirect(Application.Index)
	} else {
		return c.Forbidden("You are not allowed to create user accounts.")
	}
}

func (c User) PostLogin(Email, Password string) revel.Result {
	user := models.GetUserByEmail(c.MongoSession, Email)

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
	return c.Redirect(User.GetLogin)
}

func (c User) GetLogin() revel.Result {
	if c.UserAuthenticated() == false {
		return c.Render()
	}

	// User already logged in bounce to main page
	return c.Redirect(Application.Index)
}

func (c User) GetCreate() revel.Result {
	action := "/User/PostCreate"
	user := models.User{}
	if user.CanBeCreatedBy(c.MongoSession, c.ActiveUser) {
		return c.Render(action, user)
	}
	return c.Forbidden("You are not allowed to create user accounts.")
}

func (c User) GetLogout() revel.Result {
	for k := range c.Session {
		delete(c.Session, k)
	}
	return c.Redirect(Application.Index)
}

func (c User) GetDelete(id bson.ObjectId) revel.Result {
	user := models.GetUserByObjectId(c.MongoSession, id)
	if user != nil && user.CanBeDeletedBy(c.MongoSession, c.ActiveUser) {
		return c.Render()
	}
	return c.Forbidden("No user matching this id")
}

func (c User) GetRead(id bson.ObjectId) revel.Result {
	user := models.GetUserByObjectId(c.MongoSession, id)
	if user != nil && user.CanBeReadBy(c.MongoSession, c.ActiveUser) {
		return c.Render(user)
	}

	return c.NotFound("No user matching this id")
}

func (c User) Delete(id bson.ObjectId) revel.Result {
	user := models.GetUserByObjectId(c.MongoSession, id)
	if user != nil && user.CanBeDeletedBy(c.MongoSession, c.ActiveUser) {
		user.Delete(c.MongoSession)
		c.Flash.Success("Deleted user ", user.String())
	}
	return c.Redirect(Application.Index)
}

func (c User) EditLinks(id bson.ObjectId) revel.Result {
	canEdit := false
	u := &models.User{}
	if c.ActiveUser != nil {
		u = models.GetUserByObjectId(c.MongoSession, id)
		if u.CanBeUpdatedBy(c.MongoSession, c.ActiveUser) {
			canEdit = true
		}
	}
	return c.Render(canEdit, u)
}
