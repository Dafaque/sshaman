package users

type User struct {
	ID    int64
	Name  string
	Roles []int64
}

type internalOperationsPermissions struct{}

func (p internalOperationsPermissions) SU() bool {
	return true
}
