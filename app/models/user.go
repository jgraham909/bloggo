package models

import (
	"code.google.com/p/go.crypto/bcrypt"
	"fmt"
	"github.com/robfig/revel"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"regexp"
)

type User struct {
	Id             bson.ObjectId `bson:"_id,omitempty"`
	Firstname      string        `Firstname`
	Lastname       string        `Lastname`
	Email          string        `Email`
	HashedPassword []byte        `HashedPassword`
	Password       string
}

func (u *User) String() string {
	return fmt.Sprintf("%s %s", u.Firstname, u.Lastname)
}

var nameRegex = regexp.MustCompile("^\\w*$")

func (user *User) Validate(v *revel.Validation) {
	v.Check(user.Firstname,
		revel.Required{},
		revel.MinSize{1},
		revel.MaxSize{64},
		revel.Match{nameRegex},
	)

	v.Check(user.Lastname,
		revel.Required{},
		revel.MinSize{1},
		revel.MaxSize{64},
		revel.Match{nameRegex},
	)

	v.Check(user.Email,
		revel.Required{},
	)

	v.Email(user.Email)
}

func ValidatePassword(v *revel.Validation, password string) *revel.ValidationResult {
	return v.Check(password,
		revel.Required{},
		revel.MinSize{8},
	)
}

func (u *User) GetByEmail(s *mgo.Session, Email string) *User {
	acct := new(User)

	coll := s.DB("bloggo").C("users")
	query := coll.Find(bson.M{"Email": Email})
	query.One(acct)

	return acct
}

func (u *User) GetById(s *mgo.Session, Id bson.ObjectId) *User {
	acct := new(User)
	coll := s.DB("bloggo").C("users")
	query := coll.FindId(Id)
	query.One(acct)

	return acct
}

func (u *User) Save(s *mgo.Session) error {
	// Calculate the new password hash or load the existing one so we don't clobber it on save.
	if u.Password != "" {
		u.HashedPassword, _ = bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	} else {
		existing := u.GetById(s, u.Id)
		if existing.HashedPassword != nil {
			u.HashedPassword = existing.HashedPassword
		}
	}

	// Empty out the unhashed password to ensure it is not stored in plaintext
	u.Password = ""

	coll := s.DB("bloggo").C("users")
	_, err := coll.Upsert(bson.M{"_id": u.Id}, u)
	if err != nil {
		revel.WARN.Printf("Unable to save user account: %v error %v", u, err)
	}
	return err
}
