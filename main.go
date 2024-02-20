package main

import (
	"context"
	"github.com/blazehu/sse-demo/gen/proto"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

const (
	grpcEndpoint = "localhost:50051"
	httpPort     = ":8080"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// 创建 gRPC 连接
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	grpcConn, err := grpc.DialContext(ctx, grpcEndpoint, opts...)
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}
	defer grpcConn.Close()

	// 创建 gRPC-Gateway 服务器
	mux := runtime.NewServeMux()
	err = chat.RegisterChatServiceHandler(ctx, mux, grpcConn)
	if err != nil {
		log.Fatalf("Failed to register gateway: %v", err)
	}

	// 启动 gRPC-Gateway 服务器
	log.Printf("Starting gRPC-Gateway on %s", httpPort)
	if err := http.ListenAndServe(httpPort, mux); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
