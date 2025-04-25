package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/atDoyle/distributed-systems/key-value-store"
)

// server is our implementation of the KeyValueStoreServer interface
type server struct {
	pb.UnimplementedKeyValueStoreServer
	store map[string][]byte
}

func newServer() *server {
	return &server{
		store: make(map[string][]byte),
	}
}

// Set implements the Set RPC method
func (s *server) Set(ctx context.Context, req *pb.SetRequest) (*emptypb.Empty, error) {
	s.store[req.Key] = req.Value
	fmt.Printf("Set: key=%s, value=%v\n", req.Key, req.Value)
	return &emptypb.Empty{}, nil
}

// Get implements the Get RPC method
func (s *server) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	value, ok := s.store[req.Key]
	if !ok {
		return nil, fmt.Errorf("key not found: %s", req.Key)
	}
	fmt.Printf("Get: key=%s, value=%v\n", req.Key, value)
	return &pb.GetResponse{Value: value}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterKeyValueStoreServer(s, newServer())
	reflection.Register(s) // Enable reflection for gRPC clients like `grpcurl`
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
