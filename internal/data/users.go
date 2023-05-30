package data

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"greenlight.jaswanthp.com/internal/validator"
)

type User struct {
	Id        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  password  `json:"-"`
	Activated bool      `json:"activated"`
	Version   int       `json:"-"`
}

type password struct {
	plaintext *string
	hash      []byte
}

func (p *password) Set(plainword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainword), 12)
	if err != nil {
		return err
	}
	p.plaintext = &plainword
	p.hash = hash
	return nil
}

func (p *password) Matches(plainword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plainword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "email must be provided")
	v.Check(validator.Matches(email, validator.EmailExp), "email", "must be a valid email")
}

func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "passwrord", "password cannot be empty")
	v.Check(len(password) >= 8 && len(password) <= 72, "password", "password has to be in between 8,72 chars")
}

func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.Name != "", "name", "must be provided")
	v.Check(len(user.Name) <= 500, "name", "must not be more than 500 bytes long")
	// Call the standalone ValidateEmail() helper.
	ValidateEmail(v, user.Email)
	// If the plaintext password is not nil, call the standalone
	// ValidatePasswordPlaintext() helper.
	if user.Password.plaintext != nil {
		ValidatePasswordPlaintext(v, *user.Password.plaintext)
	}
	// If the password hash is ever nil, this will be due to a logic error in our
	// codebase (probably because we forgot to set a password for the user). It's a
	// useful sanity check to include here, but it's not a problem with the data
	// provided by the client. So rather than adding an error to the validation map we
	// raise a panic instead.
	if user.Password.hash == nil {
		panic("missing password hash for user")
	}
}
