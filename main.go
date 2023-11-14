package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "back-grpc/proto"
)

type server struct {
	pb.GreeterServer
}

func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("failed to get metadata")
	}

	// Mengambil nilai token dari metadata
	token := ""
	field := []string{}

	if len(md["authorization"]) > 0 {
		token = md["authorization"][0]
	}

	if len(md["field"]) > 0 {
		field = append(field, md["field"]...)
	}

	fmt.Println("token", token)
	fmt.Println("[]field", field)
	return &pb.HelloResponse{Greeting: "Hello, " + req.GetName()}, nil
}

func unaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	log.Println("--> unary interceptor: ", info.FullMethod)
	return handler(ctx, req)
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(unaryInterceptor),
	)
	pb.RegisterGreeterServer(s, &server{})

	fmt.Println("gRPC server is running on port 50051...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
