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

	//config := PostgresConnConfig{
	//	Host:     osutil.GetEnv(ENV_HOST, "localhost"),
	//	Port:     osutil.GetEnv(ENV_PORT, "5432"),
	//	User:     osutil.GetEnv(ENV_USER, "user"),
	//	Password: osutil.GetEnv(ENV_PASSWORD, "password"),
	//	DBName:   osutil.GetEnv(ENV_DBNAME, "seldon"),
	//	SSLMode:  osutil.GetEnv(ENV_SSLMODE, "require"),
	//	LogLevel: logLevel,
	//}

	db, err := database.PostgresGormDB(database.PostgresConnConfig{
		Address: database.Address{
			Host:   "",
			Port:   "",
			DBName: "",
		},
		Auth: database.Auth{
			User:     "",
			Password: "",
		},
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
