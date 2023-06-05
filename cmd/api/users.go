package main

import (
	"errors"
	"net/http"
	"time"

	"greenlight.jaswanthp.com/internal/data"
	"greenlight.jaswanthp.com/internal/validator"
)

func (app *application) registerUser(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.decodeJson(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &data.User{
		Name:      input.Name,
		Email:     input.Email,
		Activated: false,
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()

	if data.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateMail):
			v.AddError("email", "a user with this mail already exists")
			app.failedValidationResponse(w, r, v.Errors)

		default:
			app.serverErrorResponse(w, r, err)

		}
		return
	}

	token, err := app.models.Tokens.New(int64(user.Id), 3*24*time.Hour, data.ScopeActivation)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.background(func() {
		// envolope of data for sending both the userID and password.
		data := map[string]interface{}{
			"activationToken": token.Plaintext,
			"userID":          user.Id,
		}

		err = app.mailer.Send(user.Email, "user_welcome.tmpl.html", data)
		if err != nil {
			// Importantly, if there is an error sending the email then we use the
			// app.logger.PrintError() helper to manage it, instead of the
			// app.serverErrorResponse() helper like before.
			app.logger.PrintError(err, nil)
		}
	})

	err = app.writeJson(w, http.StatusCreated, envolope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
