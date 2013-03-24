package models

import (
	"fmt"
	"github.com/robfig/revel"
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
	return fmt.Sprintf("User(%s %s)", u.Firstname, u.Lastname)
}

var userRegex = regexp.MustCompile("^\\w*$")

func (user *User) Validate(v *revel.Validation) {
	v.Check(user.Firstname,
		revel.Required{},
		revel.MinSize{1},
		revel.MaxSize{64},
		revel.Match{userRegex},
	)

	v.Check(user.Lastname,
		revel.Required{},
		revel.MinSize{1},
		revel.MaxSize{64},
		revel.Match{userRegex},
	)

	v.Check(user.Email,
		revel.Required{},
		revel.Email{},
	)

	ValidatePassword(v, user.Password).
		Key("user.Password")

}

func ValidatePassword(v *revel.Validation, password string) *revel.ValidationResult {
	return v.Check(password,
		revel.Required{},
		revel.MinSize{8},
	)
}
