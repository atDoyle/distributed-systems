package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/atDoyle/distributed-systems/key-value-store"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	// conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewKeyValueStoreClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	key := "A"
	value := "Example value."

	// Set the key-value pair
	_, err = c.Set(ctx, &pb.SetRequest{Key: key, Value: []byte(value)})
	if err != nil {
		log.Fatalf("could not set: %v", err)
	}
	log.Printf("Set: key=%s, value=%v", key, value)

	// Get the value for the key
	r, err := c.Get(ctx, &pb.GetRequest{Key: key})
	if err != nil {
		log.Fatalf("could not get: %v", err)
	}
	retrievedValue := string(r.GetValue())
	log.Printf("Get: key=%s, value=%v", key, retrievedValue)
}
