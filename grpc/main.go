package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/SergiuPlesco/microservices-go/grpc/pb"
	"github.com/SergiuPlesco/microservices-go/grpc/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	log := log.New(os.Stdout, "grpc server", log.LstdFlags)
	gs := grpc.NewServer()
	cs := server.NewCurrency(log)

	pb.RegisterCurrencyServer(gs, cs)

	// only on dev mode, shuold be disabled in production
	reflection.Register(gs)

	l, err := net.Listen("tcp", ":9100")
	fmt.Println("Server started at 9100")
	if err != nil {
		log.Println("Unable to listen, error: ", err)
		os.Exit(1)
	}

	gs.Serve(l)
}
