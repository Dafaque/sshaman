package auth

import "context"

type permissions struct {
	read      []string
	write     []string
	delete    []string
	overwrite []string
	su        bool
	uid       int64
}

func (perms *permissions) UID() int64 {
	return perms.uid
}

func (perms *permissions) Read(space string) bool {
	if perms.isWildcard(space) {
		return true
	}
	for _, s := range perms.read {
		if s == space {
			return true
		}
	}
	return false
}

func (perms *permissions) Write(space string) bool {
	if perms.isWildcard(space) {
		return true
	}
	for _, s := range perms.write {
		if s == space {
			return true
		}
	}
	return false
}

func (perms *permissions) Delete(space string) bool {
	if perms.isWildcard(space) {
		return true
	}
	for _, s := range perms.delete {
		if s == space {
			return true
		}
	}
	return false
}

func (perms *permissions) Overwrite(space string) bool {
	if perms.isWildcard(space) {
		return true
	}
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

func (perms *permissions) isWildcard(space string) bool {
	return space == "*"
}

type RPCCredentials struct {
	token                    string
	requireTransportSecurity bool
}

func NewRPCCredentials(token string, requireTransportSecurity bool) *RPCCredentials {
	return &RPCCredentials{token: token, requireTransportSecurity: requireTransportSecurity}
}

func (c *RPCCredentials) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{
		"Authorization": "Bearer " + c.token,
	}, nil
}

func (c *RPCCredentials) RequireTransportSecurity() bool {
	return c.requireTransportSecurity
}
