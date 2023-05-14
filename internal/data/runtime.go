package data

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Runtime int32

var ErrInvalidRuntime = errors.New("invalud runtime format")

func (r Runtime) MarshalJSON() ([]byte, error) {
	jsonValue := fmt.Sprintf("%d mins", r)

	quotedJsonValue := strconv.Quote(jsonValue)

	return []byte(quotedJsonValue), nil
}

func (r *Runtime) UnmarshalJSON(jsonvalue []byte) error {
	unquotedjsonvalue, err := strconv.Unquote(string(jsonvalue))
	if err != nil {
		return ErrInvalidRuntime
	}

	parts := strings.Split(unquotedjsonvalue, " ")

	if len(parts) != 2 || parts[1] != "mins" {
		return ErrInvalidRuntime
	}
	i, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return ErrInvalidRuntime
	}

	*r = Runtime(i)

	return nil

}
