package gateway

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"net/http"
)

type registrant interface {
	RegisterGRPCGateway(ctx context.Context, mux *runtime.ServeMux)
}

type GenericHTTPServer struct {
	mux  *runtime.ServeMux
	port uint
}

type Args struct {
	Port uint `json:"port"`
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, r *http.Request) {
		// Note that in reality you would likely not want to allow all.
		response.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		response.Header().Set("Access-Control-Allow-Credentials", "true")
		response.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			response.Header().Set("Access-Control-Allow-Methods", "*")
			response.Write([]byte(""))
		} else {
			next.ServeHTTP(response, r)
		}
	})
}

func NewServer(args Args) (*GenericHTTPServer, error) {
	c := &GenericHTTPServer{
		mux:  runtime.NewServeMux(),
		port: args.Port,
	}
	return c, nil
}

func (s *GenericHTTPServer) Register(ctx context.Context, r registrant) {
	r.RegisterGRPCGateway(ctx, s.mux)
}

func (s *GenericHTTPServer) Start() error {
	address := fmt.Sprintf(":%d", s.port)
	err := http.ListenAndServe(address, corsMiddleware(s.mux))
	if err != http.ErrServerClosed {
		fmt.Println("unable to start or stop cleanly")
	}
	return err
}
