package roles

type Role struct {
	ID          int64
	Name        string
	Description string
	Read        bool     // read access
	Write       bool     // write access
	Delete      bool     // delete access
	Overwrite   bool     // overwrite access
	SU          bool     // super user access
	Spaces      []string // appliable spaces
}

type internalOperationsPermissions struct{}

func (p internalOperationsPermissions) SU() bool {
	return true
}
