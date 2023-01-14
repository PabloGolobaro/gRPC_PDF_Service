package main

import (
	"fmt"
	pb "github.com/pablogolobaro/pdfcomposer"
	"google.golang.org/grpc"
	"net"
	"service-pdf-compose-grpc/internal/handlers"
)

func main() {
	fmt.Println("Starting server ...")
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer()
	server := &handlers.Server{}
	pb.RegisterPdfComposeServer(s, server)
	if err := s.Serve(lis); err != nil {
		panic(err)
	}

}
