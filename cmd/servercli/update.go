package main

import (
	"context"
	"errors"
	"log"
	"strings"

	server "github.com/Dafaque/sshaman/pkg/server/api"
)

func updateRole(ctx context.Context, client server.RemoteCredentialsManagerClient) error {
	if *flagID == -1 {
		return errors.New("id is required")
	}
	r, err := client.UpdateRole(ctx, &server.UpdateRoleRequest{
		Role: &server.Role{
			Id:          *flagID,
			Name:        *flagName,
			Description: *flagDescription,
			Read:        *flagRead,
			Write:       *flagWrite,
			Delete:      *flagDel,
			Overwrite:   *flagOverwrite,
			Spaces:      strings.Split(*flagSpaces, ","),
		},
	})
	if err != nil {
		return err
	}
	if r.Success {
		log.Printf("Role %s updated", *flagName)
	}
	return nil
}

func updateUser(ctx context.Context, client server.RemoteCredentialsManagerClient) error {
	if *flagID == -1 {
		return errors.New("id is required")
	}
	roles, err := makeRolesArray(*flagRoles)
	if err != nil {
		return err
	}
	u, err := client.UpdateUser(ctx, &server.UpdateUserRequest{
		User: &server.User{
			Id:    *flagID,
			Name:  *flagName,
			Roles: roles,
		},
	})
	if err != nil {
		return err
	}
	if u.Success {
		log.Printf("User %s updated", *flagName)
	}
	return nil
}
