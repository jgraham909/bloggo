package models

import (
	"code.google.com/p/go.crypto/bcrypt"
	"fmt"
	"github.com/robfig/revel"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type User struct {
	Model          `bson:",inline"`
	Firstname      string              `bson:"Firstname"`
	Lastname       string              `bson:"Lastname"`
	Email          string              `bson:"Email"`
	HashedPassword []byte              `bson:"HashedPassword"`
	Meta           map[string][]string `bson:",omitempty"`
}

type Password struct {
	Pass        string
	PassConfirm string
}

func (user *User) String() string {
	return fmt.Sprintf("%s %s", user.Firstname, user.Lastname)
}

func (user *User) Validate(v *revel.Validation) {
	v.Check(user.Firstname,
		revel.Required{},
		revel.MinSize{1},
		revel.MaxSize{64},
	)

	v.Check(user.Lastname,
		revel.Required{},
		revel.MinSize{1},
		revel.MaxSize{64},
	)

	v.Check(user.Email,
		revel.Required{},
	)

	v.Email(user.Email)
}

func (user *User) ValidatePassword(v *revel.Validation, password Password) {
	v.Check(password.Pass,
		revel.MinSize{8},
	)
	v.Check(password.PassConfirm,
		revel.MinSize{8},
	)
	v.Required(password.Pass == password.PassConfirm).Message("The passwords do not match.")
}

// Save a user to the database. If a struct with p.Pass != nil is passed this
// will update the user's password as well.
// This returns the error value from mgo.Upsert()
func (user *User) Save(s *mgo.Session, p Password) error {
	// Calculate the new password hash or load the existing one so we don't clobber it on save.
	if p.Pass != "" {
		user.HashedPassword, _ = bcrypt.GenerateFromPassword([]byte(p.Pass), bcrypt.DefaultCost)
	} else {
		existing := GetUserByObjectId(s, user.Id)
		if existing.HashedPassword != nil {
			user.HashedPassword = existing.HashedPassword
		}
	}

	_, err := Collection(user, s).Upsert(bson.M{"_id": user.Id}, user)
	if err != nil {
		go revel.WARN.Printf("Unable to save user account: %v error %v", user, err)
	}
	return err
}

func (user *User) Delete(s *mgo.Session) error {
	var err error
	return err
}

func GetUserByEmail(s *mgo.Session, Email string) *User {
	acct := new(User)
	query := Collection(acct, s).Find(bson.M{"Email": Email})
	query.One(acct)

	return acct
}

func GetUserByObjectId(s *mgo.Session, Id bson.ObjectId) *User {
	u := new(User)
	query := Collection(u, s).FindId(Id)
	query.One(u)
	return u
}

func (user *User) CanBeCreatedBy(s *mgo.Session, u *User) bool {
	if u == nil {
		return false
	}

	// Only admin can create, set via app.conf:bloggo.admin hex value
	if a, found := revel.Config.String("bloggo.admin"); found && u != nil {
		if bson.IsObjectIdHex(a) && a == u.Id.Hex() {
			return true
		}
	}
	return false
}

func (user *User) CanBeReadBy(s *mgo.Session, u *User) bool {
	// Default everybody can read.
	return true
}

func (user *User) CanBeDeletedBy(s *mgo.Session, u *User) bool {
	// Only admin can create, set via app.conf:bloggo.admin hex value
	if a, found := revel.Config.String("bloggo.admin"); found {
		if bson.IsObjectIdHex(a) && a == u.Id.Hex() {
			return true
		}
	}
	return false
}

func (user *User) CanBeUpdatedBy(s *mgo.Session, u *User) bool {
	if u == nil {
		return false
	}
	// Only admin can create, set via app.conf:bloggo.admin hex value
	if a, found := revel.Config.String("bloggo.admin"); found {
		if bson.IsObjectIdHex(a) && a == u.Id.Hex() {
			return true
		}
	}

	// User's can edit their own account
	if user.Id == u.Id {
		return true
	}
	return false
}
