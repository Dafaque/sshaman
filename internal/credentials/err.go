package credentials

import "fmt"

const (
	codeCredentialsExist = 1
)

type Error struct {
	msg     string
	code    int
	details map[string]interface{}
}

func (e Error) Error() string {
	msg := e.msg
	for k, v := range e.details {
		msg += fmt.Sprintf(" %s=%v", k, v)
	}
	return msg
}
func (e Error) Is(target error) bool {
	if target == nil {
		return false
	}
	te, ok := target.(Error)
	if !ok {
		return false
	}
	return e.code == te.code
}

var ErrCredentialsExist = Error{
	msg:  "credentials already exists",
	code: codeCredentialsExist,
}

func NewCredentialExistsError(alias string) Error {
	err := ErrCredentialsExist
	err.details = map[string]interface{}{"alias": alias}
	return err
}
