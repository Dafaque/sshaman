package credentials

type (
	permissions interface {
		UID() int64
		Read(space string) bool
		Write(space string) bool
		Delete(space string) bool
		Overwrite(space string) bool
		SU() bool
	}

	repository interface {
		Create(cred *Credentials) error
		Get(cred *Credentials) error
		Update(cred *Credentials) error
		Delete(cred *Credentials) error
	}
)
