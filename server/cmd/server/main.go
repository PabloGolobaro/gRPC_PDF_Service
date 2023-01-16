package main

import (
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	pb "github.com/pablogolobaro/pdfcomposer"
	"google.golang.org/grpc"
	"log"
	"net"
	"service-pdf-compose-grpc/internal/handlers"
)

func main() {
	fmt.Println("Starting server ...")
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer(grpc_middleware.WithStreamServerChain(metricInterceptor, logInterceptor))
	server := &handlers.Server{}
	pb.RegisterPdfComposeServer(s, server)
	if err := s.Serve(lis); err != nil {
		panic(err)
	}
}
func metricInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	fmt.Println("Metric Interceptor")
	log.Println(info.FullMethod)
	err := handler(srv, ss)
	if err != nil {
		return err
	}
	return nil
}
func logInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	fmt.Println("Log Interceptor")

	err := handler(srv, ss)
	if err != nil {
		return err
	}
	return nil
}
