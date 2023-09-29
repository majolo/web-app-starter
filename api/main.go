package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"
)

var (
	grpcPort = flag.Int("port", 50051, "The server grpc port")
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("server listening at %v", lis.Addr())

	chartServer := implementation.NewChartService(fsClient, websiteDomain)
	grpcServer := grpc.NewServer()
	charts.RegisterChartsServer(grpcServer, chartServer)

	httpServer.Register(context.Background(), chartServer)

	// Start services
	serviceWaiter := sync.WaitGroup{}
	serviceWaiter.Add(1)
	go func() {
		err := httpServer.Start()
		if err != nil {
			fmt.Println("unableToStartMsg")
		}
		serviceWaiter.Done()
	}()

	// Graceful shutdown
	serviceWaiter.Wait()
}
