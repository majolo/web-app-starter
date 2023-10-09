package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/majolo/web-app-starter/gateway"
	"github.com/majolo/web-app-starter/gen/diary/v1"
	"github.com/majolo/web-app-starter/services"
	"github.com/nedpals/supabase-go"
	"google.golang.org/grpc"
	"log"
	"os"
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

	supabaseUrl := os.Getenv("NEXT_PUBLIC_SUPABASE_URL")
	supabaseAnonKey := os.Getenv("NEXT_PUBLIC_SUPABASE_ANON_KEY")
	supabaseClient := supabase.CreateClient(supabaseUrl, supabaseAnonKey)

	// Register services
	diaryService := services.NewDiaryService(supabaseClient)
	diary.RegisterDiaryServiceServer(grpcServer, diaryService)
	httpServer.Register(context.Background(), diaryService)

	// Start servers (just http for now)
	fmt.Println("starting http server")
	err = httpServer.Start()
	if err != nil {
		fmt.Println("did not start http server cleanly")
	}
}
