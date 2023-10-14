package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/majolo/web-app-starter/database"
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
	dbHost   = flag.String("db-host", "localhost", "The database host")
	dbPort   = flag.String("db-port", "5432", "The database port")
	dbUser   = flag.String("db-user", "postgres", "The database user")
	dbPass   = flag.String("db-pass", "password", "The database password")
	dbName   = flag.String("db-name", "diarydb", "The database name")
	//dbSSL    = flag.String("db-ssl", "require", "The database ssl mode")
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

	db, err := database.PostgresGormDB(database.PostgresConnConfig{
		Address: database.Address{
			Host:   *dbHost,
			Port:   *dbPort,
			DBName: *dbName,
		},
		Auth: database.Auth{
			User:     *dbUser,
			Password: *dbPass,
		},
		// TODO: these may be needed
		TLSConfig: database.TLSConfig{
			TLSMode:    "",
			CACertPath: "",
			KeyPath:    "",
			CertPath:   "",
		},
		PoolConfig: database.PoolConfig{
			MaxIdleConns:    0,
			MaxOpenConns:    0,
			ConnMaxLifetime: 0,
		},
	})
	if err != nil {
		log.Fatalf("failed to create db from config: %v", err)
	}

	// Register services
	diaryService, err := services.NewDiaryService(supabaseClient, db)
	if err != nil {
		log.Fatalf("failed to create diary service: %v", err)
	}
	diary.RegisterDiaryServiceServer(grpcServer, diaryService)
	httpServer.Register(context.Background(), diaryService)

	// Start servers (just http for now)
	fmt.Println("starting http server")
	err = httpServer.Start()
	if err != nil {
		fmt.Println("did not start http server cleanly")
	}
}
