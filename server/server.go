package main

import (
	"github.com/blazehu/sse-demo/gen/proto"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net"
	"time"
)

// Server provides chat service
type Server struct {
	chat.UnimplementedChatServiceServer
}

// Chat returns chat content
func (s *Server) Chat(_ *emptypb.Empty, stream chat.ChatService_ChatServer) error {
	for {
		msg := chat.Message{
			User:    "blazehu",
			Content: time.Now().Format(time.RFC3339),
		}
		if err := stream.Send(&msg); err != nil {
			return err
		}
		time.Sleep(time.Second * 1)
	}
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	chat.RegisterChatServiceServer(grpcServer, &Server{})
	grpcServer.Serve(lis)
}
