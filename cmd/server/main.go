package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/DarRo9/pvz_service/config"
	"github.com/DarRo9/pvz_service/internal/db"
	internal_grpc "github.com/DarRo9/pvz_service/internal/grpc"
	"github.com/DarRo9/pvz_service/internal/grpc/pvz/pvz_v1"
	handler "github.com/DarRo9/pvz_service/internal/handler"
	internal_middleware "github.com/DarRo9/pvz_service/internal/middleware"
	"github.com/DarRo9/pvz_service/internal/repository"
	"github.com/DarRo9/pvz_service/internal/service"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

func main() {

	dbCfg := db.DatabaseConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
	}

	db, err := db.NewDatabase(&dbCfg)
	if err != nil {
		log.Fatalf("Error connecting to the database:", err)
		return
	}
	defer db.Close()

	config, err := config.LoadConfig(
		"config/config.yaml",
	)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	repo := repository.NewPostgresRepository(db)
	service := service.NewService(repo, config)
	httpHandler := handler.NewHTTPHandler(service)
	grpcHandler := internal_grpc.NewGRPCHandler(service)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	// Запускаем HTTP сервер
	wg.Add(1)
	go func() {
		defer wg.Done()
		startHTTPServer(ctx, httpHandler)
	}()

	// Запускаем gRPC сервер
	wg.Add(1)
	go func() {
		defer wg.Done()
		startGRPCServer(ctx, grpcHandler)
	}()

	// Запускаем Metrics сервер
	wg.Add(1)
	go func() {
		defer wg.Done()
		startMetricsServer(ctx)
	}()

	log.Println("Servers started")

	<-done
	log.Println("Servers stopping...")

	cancel()

	wg.Wait()
	log.Println("Servers stopped")
}

func startMetricsServer(ctx context.Context) {
	srv := &http.Server{
		Addr:    ":9000",
		Handler: promhttp.Handler(),
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Metrics server error: %s\n", err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Metrics server shutdown error: %+v\n", err)
	} else {
		log.Println("Metrics server stopped gracefully")
	}
}

func startHTTPServer(ctx context.Context, h *handler.HTTPHandler) {
	r := chi.NewRouter()

	wrapper := handler.ServerInterfaceWrapper{
		Handler:            h,
		HandlerMiddlewares: nil,
		ErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		},
	}

	r.Use(middleware.Logger)
	r.Use(internal_middleware.PrometheusMiddleware)
	r.Post("/dummyLogin", wrapper.PostDummyLogin)
	r.Post("/login", wrapper.PostLogin)
	r.Post("/register", wrapper.PostRegister)

	r.Route("/", func(r chi.Router) {
		r.Use(internal_middleware.AuthMiddleware)
		r.Post("/products", wrapper.PostProducts)
		r.Get("/pvz", wrapper.GetPvz)
		r.Post("/pvz", wrapper.PostPvz)
		r.Post("/pvz/{pvzId}/close_last_reception", wrapper.PostPvzPvzIdCloseLastReception)
		r.Post("/pvz/{pvzId}/delete_last_product", wrapper.PostPvzPvzIdDeleteLastProduct)
		r.Post("/receptions", wrapper.PostReceptions)
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %s\n", err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server shutdown error: %+v\n", err)
	} else {
		log.Println("HTTP server stopped gracefully")
	}
}

func startGRPCServer(ctx context.Context, userHandler *internal_grpc.GRPCHandler) {
	grpcServer := grpc.NewServer()
	pvz_v1.RegisterPVZServiceServer(grpcServer, userHandler)

	lis, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("gRPC server error: %s\n", err)
		}
	}()

	<-ctx.Done()

	stopped := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		close(stopped)
	}()

	select {
	case <-stopped:
		log.Println("gRPC server stopped gracefully")
	case <-time.After(5 * time.Second):
		grpcServer.Stop()
		log.Println("gRPC server stopped (force)")
	}
}
