package main

import (
	"context"
	"errors"
	"log"
	"strconv"
	"strings"

	server "github.com/Dafaque/sshaman/pkg/server/api"
)

func createRole(ctx context.Context, client server.RemoteCredentialsManagerClient) error {
	if *flagName == "" {
		return errors.New("name is required")
	}
	if *flagSpaces == "" {
		return errors.New("spaces are required")
	}
	if !*flagRead && !*flagWrite && !*flagDel && !*flagSuper {
		return errors.New("user must have at least one permission")
	}
	r, err := client.CreateRole(ctx, &server.CreateRoleRequest{
		Role: &server.Role{
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
		log.Printf("Role %s created", *flagName)
	}
	return nil
}

func createUser(ctx context.Context, client server.RemoteCredentialsManagerClient) error {
	if *flagName == "" {
		return errors.New("name is required")
	}
	if *flagRoles == "" {
		return errors.New("roles are required")
	}
	roles, err := makeRolesArray(*flagRoles)
	if err != nil {
		return err
	}
	r, err := client.CreateUser(ctx, &server.CreateUserRequest{
		User: &server.User{
			Name:  *flagName,
			Roles: roles,
		},
	})
	if err != nil {
		return err
	}
	if r.Success {
		log.Printf("User %s created", *flagName)
	}
	return nil
}

func makeRolesArray(roles string) ([]int64, error) {
	roleIDs := strings.Split(roles, ",")
	var iroles []int64
	for _, id := range roleIDs {
		roleID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return nil, err
		}
		iroles = append(iroles, roleID)
	}
	return iroles, nil
}
