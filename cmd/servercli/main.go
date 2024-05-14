package main

import (
	"context"
	"flag"
	"log"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/Dafaque/sshaman/internal/server/auth"
	server "github.com/Dafaque/sshaman/pkg/server/api"
)

const (
	address = "localhost:50051"
	token   = "eyJhbGciOiJFUzM4NCIsInR5cCI6IkpXVCJ9.eyJhdWQiOiI0OTk2MGRlNTg4MGU4YzY4NzQzNDE3MGY2NDc2NjA1YjhmZTRhZWI5YTI4NjMyYzc5OTVjZjNiYTgzMWQ5NzYzIiwiZXhwIjoxNzE1Nzg3NDczLCJpYXQiOjE3MTU3MDEwNzMsImlzcyI6IjQ5OTYwZGU1ODgwZThjNjg3NDM0MTcwZjY0NzY2MDViOGZlNGFlYjlhMjg2MzJjNzk5NWNmM2JhODMxZDk3NjMiLCJuYmYiOjE3MTU3MDEwNzMsInN1YiI6IjEifQ.xq6qX0FMS2VoNDmGyFxpxd1ZkqyoYJNY0nbwMW-03rOv0DIpa7nz3onHb0XcMlRTgTP4pcMyi1xm73nFdZCJeiKnf4tCt_0XMpHGKORINJakJTMmsL448WrmQ2L8NOcC"
)

const (
	entityRoles = "roles"
	entityUsers = "users"
)

var entities = []string{entityUsers, entityRoles}

var (
	flagList   string
	flagCreate string
	flagUpdate string
	flagDelete string

	flagID   = flag.Int64("id", -1, "ID of the entity to get, set or delete")
	flagName = flag.String("name", "", "name of the entity to create")

	// MARK: - flags for Roles
	flagDescription = flag.String("description", "", "description of the role to create")
	flagRead        = flag.Bool("r", false, "read permission for role")
	flagWrite       = flag.Bool("w", false, "write permission for role")
	flagDel         = flag.Bool("d", false, "delete permission for role")
	flagOverwrite   = flag.Bool("o", false, "overwrite permission for role")
	flagSuper       = flag.Bool("super", false, "super permission for role")
	flagSpaces      = flag.String("s", "", "comma separated list of spaces applied for role")

	// MARK: - flags for Users
	flagRoles = flag.String("roles", "", "comma separated list of role id's applied for user")
)

type operation func(ctx context.Context, client server.RemoteCredentialsManagerClient) error

func main() {
	flag.StringVar(&flagList, "list", "", "target entity to list: "+strings.Join(entities, ", "))
	flag.StringVar(&flagCreate, "create", "", "target entity to create: "+strings.Join(entities, ", "))
	flag.StringVar(&flagUpdate, "update", "", "target entity to update: "+strings.Join(entities, ", "))
	flag.StringVar(&flagDelete, "delete", "", "target entity to delete: "+strings.Join(entities, ", "))
	flagTimeout := flag.Duration("timeout", time.Second, "timeout for the operation")

	flag.Parse()

	var op operation

	switch {
	case flagList != "":
		switch flagList {
		case entityRoles:
			op = listRoles
		case entityUsers:
			op = listUsers
		default:
			log.Fatalf("invalid target entity: %s", flagList)
		}
	case flagCreate != "":
		switch flagCreate {
		case entityRoles:
			op = createRole
		case entityUsers:
			op = createUser
		default:
			log.Fatalf("invalid target entity: %s", flagCreate)
		}
	case flagUpdate != "":
		switch flagUpdate {
		case entityRoles:
			op = updateRole
		case entityUsers:
			op = updateUser
		default:
			log.Fatalf("invalid target entity: %s", flagUpdate)
		}
	case flagDelete != "":
		switch flagDelete {
		case entityRoles:
			op = deleteRole
		case entityUsers:
			op = deleteUser
		default:
			log.Fatalf("invalid target entity: %s", flagDelete)
		}
	default:
		log.Println("no operation selected")
		flag.Usage()
		os.Exit(1)
	}

	conn, err := grpc.Dial(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithPerRPCCredentials(auth.NewRPCCredentials(token, false)),
	)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	client := server.NewRemoteCredentialsManagerClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), *flagTimeout)
	defer cancel()

	if err := op(ctx, client); err != nil {
		log.Fatalf("failed to execute operation: %v", err)
	}
}
