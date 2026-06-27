// gRPC Client connects to the UserService and makes RPC calls.
//
// LEARNING NOTES:
// - gRPC clients get a type-safe stub generated from the .proto file
// - The stub looks like a regular function call but goes over the network
// - Context carries deadlines/timeouts across the wire (distributed timeout propagation!)
// - For C++ devs: the generated client stub is similar to COM/CORBA proxies
//
// Run: go run cmd/client/main.go
package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	pb "github.com/GauravAgarwalGarg/grin-with-golang/src/10_distributed/grpc_service/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Connect to gRPC server (insecure for local dev; use TLS in production)
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

	// --- Unary RPC: Create a user ---
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	created, err := client.CreateUser(ctx, &pb.CreateUserRequest{
		Name:     "Gaurav Agarwal",
		Email:    "gaurav@example.com",
		Password: "secure123",
	})
	if err != nil {
		log.Fatalf("CreateUser failed: %v", err)
	}
	fmt.Printf("Created: %+v\n", created)

	// --- Unary RPC: Get the user back ---
	fetched, err := client.GetUser(ctx, &pb.GetUserRequest{Id: created.Id})
	if err != nil {
		log.Fatalf("GetUser failed: %v", err)
	}
	fmt.Printf("Fetched: %+v\n", fetched)

	// --- Server-streaming RPC: List all users ---
	stream, err := client.ListUsers(ctx, &pb.ListUsersRequest{PageSize: 100})
	if err != nil {
		log.Fatalf("ListUsers failed: %v", err)
	}

	fmt.Println("\nAll users:")
	for {
		user, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("stream error: %v", err)
		}
		fmt.Printf("  - %s (%s)\n", user.Name, user.Email)
	}
}
