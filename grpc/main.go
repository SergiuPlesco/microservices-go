package main

import (
	"fmt"
	"net"
	"os"

	"github.com/SergiuPlesco/microservices-go/grpc/data"
	"github.com/SergiuPlesco/microservices-go/grpc/pb"
	"github.com/SergiuPlesco/microservices-go/grpc/server"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	log := hclog.Default()

	rates, err := data.NewRates(log)
	if err != nil {
		log.Error("unable to generate rates", "error", err)
		os.Exit(1)
	}
	gs := grpc.NewServer()
	cs := server.NewCurrency(rates, log)

	pb.RegisterCurrencyServer(gs, cs)

	// only on dev mode, shuold be disabled in production
	reflection.Register(gs)

	l, err := net.Listen("tcp", ":9101")
	fmt.Println("GRPC Server started at 9101")
	if err != nil {
		log.Error("Unable to listen, error: ", err)
		os.Exit(1)
	}

	gs.Serve(l)
}
