package rest

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/handlers"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	v1 "github.com/mark/todo/services/pkg/api/v1"
	"google.golang.org/grpc"
)

// RunServer Run HTTP/REST API Gateway
func RunServer(ctx context.Context, grpcPort, httpPort string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	if err := v1.RegisterToDoServiceHandlerFromEndpoint(ctx, mux, "localhost:"+grpcPort, opts); err != nil {
		log.Fatalf("Failed to start HTTP gateway: %v", err)
	}
	headersOk := handlers.AllowedHeaders([]string{"Content-Type"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	srv := &http.Server{
		Addr:    ":" + httpPort,
		Handler: handlers.CORS(originsOk, headersOk, methodsOk)(mux),
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			log.Printf("Gracefully stopping HTTP/REST gateway...")
			srv.SetKeepAlivesEnabled(false)
			if err := srv.Shutdown(ctx); err != nil {
				log.Fatalf("Could not gracefully shutdown the server: %v\n", err)
			}
			<-ctx.Done()
		}
		_, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}()
	log.Println("Starting HTTP/REST gateway...")
	return srv.ListenAndServe()
}
