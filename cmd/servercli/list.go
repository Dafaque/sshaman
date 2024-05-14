package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	server "github.com/Dafaque/sshaman/pkg/server/api"
)

func listRoles(ctx context.Context, client server.RemoteCredentialsManagerClient) error {
	r, err := client.ListRoles(ctx, &server.ListRolesRequest{})
	if err != nil {
		return err
	}
	tw := tabwriter.NewWriter(os.Stdout, 1, 1, 2, ' ', 0)
	fmt.Fprintf(tw, "ID\tName\tPermissions\tSpaces\n")
	for _, role := range r.GetRoles() {
		perms := bytes.NewBuffer(nil)
		if role.GetRead() {
			perms.WriteString("r")
		}
		if role.GetWrite() {
			perms.WriteString("w")
		}
		if role.GetDelete() {
			perms.WriteString("d")
		}
		if role.GetOverwrite() {
			perms.WriteString("o")
		}
		if role.GetSu() {
			perms.Reset()
			perms.WriteString("su")
		}
		fmt.Fprintf(
			tw,
			"%d\t%s\t%s\t%s\n",
			role.GetId(),
			role.GetName(),
			perms.String(),
			strings.Join(role.GetSpaces(), ", "),
		)
	}
	tw.Flush()
	return nil
}

func listUsers(ctx context.Context, client server.RemoteCredentialsManagerClient) error {
	r, err := client.ListUsers(ctx, &server.ListUsersRequest{})
	if err != nil {
		return err
	}
	users := r.GetUsers()

	roles, err := client.ListRoles(ctx, &server.ListRolesRequest{})
	if err != nil {
		return err
	}
	roleMap := make(map[int64]string)
	for _, role := range roles.GetRoles() {
		roleMap[role.GetId()] = role.GetName()
	}

	tw := tabwriter.NewWriter(os.Stdout, 1, 1, 2, ' ', 0)
	fmt.Fprintf(tw, "ID\tName\tRoles\n")
	for _, user := range users {
		var roles []string
		for _, role := range user.GetRoles() {
			roleName, exists := roleMap[role]
			if !exists {
				roleName = fmt.Sprintf("unknown(id=%d)", role)
			}
			roles = append(roles, roleName)
		}
		fmt.Fprintf(tw, "%d\t%s\t%s\n", user.GetId(), user.GetName(), strings.Join(roles, ", "))
	}
	tw.Flush()
	return nil
}
