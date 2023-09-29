package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/majolo/web-app-starter/gateway"
	"github.com/majolo/web-app-starter/gen/diary/v1"
	"github.com/majolo/web-app-starter/services"
	"google.golang.org/grpc"
	"log"
)

var (
	httpPort = flag.Int("port", 8080, "The server http port")
)

func main() {
	flag.Parse()
	httpServer, err := gateway.NewServer(gateway.Args{
		Port: uint(*httpPort),
	})
	if err != nil {
		log.Fatalf("failed to create http server: %v", err)
	}
	grpcServer := grpc.NewServer()

	// Register services
	diaryService := services.NewDiaryService()
	diary.RegisterDiaryServiceServer(grpcServer, diaryService)
	httpServer.Register(context.Background(), diaryService)

	// Start servers (just http for now)
	err = httpServer.Start()
	if err != nil {
		fmt.Println("did not start http server cleanly")
	}
}
