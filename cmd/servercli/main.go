package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/Dafaque/sshaman/internal/remote/auth"
	remote "github.com/Dafaque/sshaman/pkg/remote/api"
)

const (
	address = "localhost:50051"
	token   = "eyJhbGciOiJFUzM4NCIsInR5cCI6IkpXVCJ9.eyJhdWQiOiI0OTk2MGRlNTg4MGU4YzY4NzQzNDE3MGY2NDc2NjA1YjhmZTRhZWI5YTI4NjMyYzc5OTVjZjNiYTgzMWQ5NzYzIiwiZXhwIjoxNzE1MDA3MTY3LCJpYXQiOjE3MTQ5MjA3NjcsImlzcyI6IjQ5OTYwZGU1ODgwZThjNjg3NDM0MTcwZjY0NzY2MDViOGZlNGFlYjlhMjg2MzJjNzk5NWNmM2JhODMxZDk3NjMiLCJuYmYiOjE3MTQ5MjA3NjcsInN1YiI6IjEifQ.ps7eJ7VRdJW-MbaPTbltxMCJkdt7gSnrpm7dJXgEU6qYv4HQMdAK0eIS-rJVzMMBu6kZBUujXZ9eICyNaAgPA9SqtpMeTYDWyAi2np59IPKOe6Eor00F05MmKgdpQtPT"
)

func main() {
	conn, err := grpc.Dial(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithPerRPCCredentials(auth.NewRPCCredentials(token, false)),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := remote.NewRemoteCredentialsManagerClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Example call to ListRoles
	r, err := client.ListRoles(ctx, &remote.ListRolesRequest{})
	if err != nil {
		log.Fatalf("could not list roles: %v", err)
	}
	log.Printf("Roles: %v", r.GetRoles())
}
