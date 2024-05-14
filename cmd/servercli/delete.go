package main

import (
	"context"
	"errors"
	"log"

	server "github.com/Dafaque/sshaman/pkg/server/api"
)

func deleteRole(ctx context.Context, client server.RemoteCredentialsManagerClient) error {
	if *flagID == -1 {
		return errors.New("id is required")
	}
	r, err := client.DeleteRole(ctx, &server.DeleteRoleRequest{
		Id: *flagID,
	})
	if err != nil {
		return err
	}
	if r.Success {
		log.Printf("Role %d deleted", *flagID)
	}
	return nil
}

func deleteUser(ctx context.Context, client server.RemoteCredentialsManagerClient) error {
	if *flagID == -1 {
		return errors.New("id is required")
	}
	r, err := client.DeleteUser(ctx, &server.DeleteUserRequest{
		Id: *flagID,
	})
	if err != nil {
		return err
	}
	if r.Success {
		log.Printf("User %d deleted", *flagID)
	}
	return nil
}
