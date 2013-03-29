package models

import (
	"code.google.com/p/go.crypto/bcrypt"
	"fmt"
	"github.com/robfig/revel"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type User struct {
	Id             bson.ObjectId `bson:"_id,omitempty"`
	Firstname      string        `bson:"Firstname"`
	Lastname       string        `bson:"Lastname"`
	Email          string        `bson:"Email"`
	HashedPassword []byte        `bson:"HashedPassword"`
	Meta           map[string][]string
}

type Password struct {
	Pass        string
	PassConfirm string
}

// Return the appropriate collection instance for this user.
func (user *User) Collection(s *mgo.Session) *mgo.Collection {
	return s.DB("bloggo").C("users")
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
		existing := user.GetById(s, user.Id)
		if existing.HashedPassword != nil {
			user.HashedPassword = existing.HashedPassword
		}
	}

	coll := user.Collection(s)
	_, err := coll.Upsert(bson.M{"_id": user.Id}, user)
	if err != nil {
		go revel.WARN.Printf("Unable to save user account: %v error %v", user, err)
	}
	return err
}

func (user *User) Delete(s *mgo.Session) error {
	var err error
	return err
}

func (user *User) GetByEmail(s *mgo.Session, Email string) *User {
	acct := new(User)

	coll := user.Collection(s)
	query := coll.Find(bson.M{"Email": Email})
	query.One(acct)

	return acct
}

func (user *User) GetById(s *mgo.Session, Id bson.ObjectId) *User {
	acct := new(User)
	coll := user.Collection(s)
	query := coll.Find(bson.M{"_id": Id})
	err := query.One(acct)
	if err != nil {
		revel.WARN.Printf("Unable to load user by Id: %v error %v", Id, err)
	}
	return acct
}
