package auth

type permissions struct {
	read      []string
	write     []string
	delete    []string
	overwrite []string
	su        bool
}

func (perms *permissions) Read(space string) bool {
	for _, s := range perms.read {
		if s == space {
			return true
		}
	}
	return false
}

func (perms *permissions) Write(space string) bool {
	for _, s := range perms.write {
		if s == space {
			return true
		}
	}
	return false
}

func (perms *permissions) Delete(space string) bool {
	for _, s := range perms.delete {
		if s == space {
			return true
		}
	}
	return false
}

func (perms *permissions) Overwrite(space string) bool {
	for _, s := range perms.overwrite {
		if s == space {
			return true
		}
	}
	return false
}

func (perms *permissions) SU() bool {
	return perms.su
}
