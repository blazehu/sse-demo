package main

import (
	"context"
	"github.com/blazehu/sse-demo/gen/proto"
	"github.com/blazehu/sse-demo/util"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	googleproto "google.golang.org/protobuf/proto"
	"log"
	"net/http"
)

const (
	grpcEndpoint = "localhost:50051"
	httpPort     = ":8080"
)

func myFilter(ctx context.Context, w http.ResponseWriter, resp googleproto.Message) error {
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Cache-Control", "no-cache")
	return nil
}

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
	gwmux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, util.NewCustomTranscoder(&runtime.JSONPb{})),
		runtime.WithForwardResponseOption(myFilter),
	)
	err = chat.RegisterChatServiceHandler(ctx, gwmux, grpcConn)
	if err != nil {
		log.Fatalf("Failed to register gateway: %v", err)
	}

	// 设置 CORS 策略
	mux := http.NewServeMux()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		gwmux.ServeHTTP(w, r)
	})
	mux.Handle("/", handler)

	// 启动 gRPC-Gateway 服务器
	log.Printf("Starting gRPC-Gateway on %s", httpPort)
	if err := http.ListenAndServe(httpPort, mux); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
