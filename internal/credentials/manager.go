package credentials

type Manager interface {
	Get(alias string) (*Credentials, error)
	Set(cred *Credentials, force bool) error
	List() ([]*Credentials, error)
	Del(alias string) error
	Drop() error
	Done() error
}
