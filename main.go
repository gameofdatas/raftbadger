package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"raftbadger/pkg/config"
	"raftbadger/pkg/factory"
	"raftbadger/pkg/proto"
	"raftbadger/server"

	"google.golang.org/grpc/reflection"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	svr, err := server.NewServer(cfg, cfg.BootstrapNode)
	if err != nil {
		log.Fatalf("error initializing server %v", err)
	}
	grpcServer := factory.InitializeGRPCServer()

	svr.Tm.Register(grpcServer)

	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", cfg.Port, err)
	}

	proto.RegisterBadgerServiceServer(grpcServer, svr)

	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop

		log.Println("Shutting down gRPC server...")
		grpcServer.GracefulStop()
		log.Println("gRPC server stopped.")
	}()

	log.Printf("Starting gRPC server on port %s", cfg.Port)
	if err = grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
