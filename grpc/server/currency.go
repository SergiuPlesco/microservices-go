package server

import (
	"context"

	"github.com/SergiuPlesco/microservices-go/grpc/data"
	"github.com/SergiuPlesco/microservices-go/grpc/pb"
	"github.com/hashicorp/go-hclog"
)

type Currency struct {
	pb.UnimplementedCurrencyServer
	r *data.ExchangeRates
	l hclog.Logger
}

func NewCurrency(r *data.ExchangeRates, l hclog.Logger) *Currency {
	return &Currency{l: l, r: r}
}

func (c *Currency) GetRate(ctx context.Context, req *pb.RateRequest) (*pb.RateResponse, error) {
	c.l.Debug("Handle GetRate", "base", req.GetBase(), "destination", req.GetDestination())

	rate, err := c.r.GetRate(req.GetBase().String(), req.GetDestination().String())
	if err != nil {
		return nil, err
	}

	return &pb.RateResponse{Rate: rate}, nil
}
