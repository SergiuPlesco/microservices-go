package server

import (
	"context"
	"log"

	"github.com/SergiuPlesco/microservices-go/grpc/pb"
)

type Currency struct {
	pb.UnimplementedCurrencyServer
	l *log.Logger
}

func NewCurrency(l *log.Logger) *Currency {
	return &Currency{l: l}
}

func (c *Currency) GetRate(ctx context.Context, req *pb.RateRequest) (*pb.RateResponse, error) {
	c.l.Println("Handle GetRate", "base", req.GetBase(), "destination", req.GetDestination())

	return &pb.RateResponse{Rate: 0.5}, nil
}
