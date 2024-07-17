package main

import (
	"authentication_service/config"
	"authentication_service/pkg"
	pb "authentication_service/genproto/authentication_service"
	"authentication_service/service"
	repositories"authentication_service/storege/postgres"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {

	config := config.Load()

	db, err := pkg.ConnectDB()
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer db.Close()
	listener, err := net.Listen("tcp", ":"+config.URL_PORT)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	defer listener.Close()

	log.Printf("Server started on port " + config.URL_PORT)

	authStorage := repositories.NewUserRepo(db)

	as := service.NewAuthService(authStorage)

	s := grpc.NewServer()
	pb.RegisterEcoServiceServer(s, as)

	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
