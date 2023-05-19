package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"greenlight.jaswanthp.com/internal/validator"
)

type envolope map[string]interface{}

// func (app *application) readParamID(r *http.Request) (int64, error) {
// 	params := httprouter.ParamsFromContext(r.Context())

// 	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
// 	if err != nil || id < 1 {
// 		return 0, errors.New("invalid parameters")
// 	}
// 	return id, nil
// }

func (app *application) writeJson(w http.ResponseWriter, status int, data interface{}, headers http.Header) error {

	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil

}

func (app *application) decodeJson(w http.ResponseWriter, r *http.Request, dst interface{}) error {

	max_bytes := 1_048_576
	// limiting the body size to 1Mb
	r.Body = http.MaxBytesReader(w, r.Body, int64(max_bytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unMarshallError *json.UnmarshalTypeError
		var invalidUnMarshalError *json.InvalidUnmarshalError

		switch {
		// Use the errors.As() function to check whether the error has the type
		// *json.SyntaxError. If it does, then return a plain-english error message
		// which includes the location of the problem.
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

		// In some circumstances Decode() may also return an io.ErrUnexpectedEOF error
		// for syntax errors in the JSON. So we check for this using errors.Is() and
		// return a generic error message. There is an open issue regarding this at
		// https://github.com/golang/go/issues/25956.
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly formatted json")

		// Likewise, catch any *json.UnmarshalTypeError errors. These occur when the
		// JSON value is the wrong type for the target destination. If the error relates
		// to a specific field, then we include that in our error message to make it
		// easier for the client to debug.

		case errors.As(err, &unMarshallError):
			if unMarshallError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unMarshallError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unMarshallError.Offset)

		// An io.EOF error will be returned by Decode() if the request body is empty. We
		// check for this with errors.Is() and return a plain-english error message
		// instead.
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		// unknown feild error handling
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown field %s", fieldName)

		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %d bytes", max_bytes)
		// A json.InvalidUnmarshalError error will be returned if we pass a non-nil
		// pointer to Decode(). We catch this and panic, rather than returning an error
		// to our handler. At the end of this chapter we'll talk about panicking
		// versus returning errors, and discuss why it's an appropriate thing to do in
		// this specific situation.
		case errors.As(err, &invalidUnMarshalError):
			panic(err)
		// For anything else, return the error message as-is.
		default:
			return err
		}

	}

	// checking so that we have only one json value in the request
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body should contain only one json value")
	}

	return nil

}
func (app *application) readString(qs url.Values, key string, defaultvalue string) string {

	s := qs.Get(key)

	if s == "" {
		return defaultvalue
	}
	return s
}

func (app *application) readCSV(qs url.Values, key string, defaultvalue []string) []string {
	s := qs.Get(key)

	if s == "" {
		return defaultvalue
	}

	return strings.Split(s, ",")
}

func (app *application) readInt(qs url.Values, key string, defaultvalue int, v *validator.Validator) int {
	i := qs.Get(key)

	if i == "" {
		return defaultvalue
	}

	n, err := strconv.Atoi(i)
	if err != nil {
		v.AddError(key, "must be an integer")
		return defaultvalue
	}
	return n
}
