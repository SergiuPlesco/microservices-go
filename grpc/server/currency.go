package server

import (
	"context"
	"io"
	"time"

	"github.com/SergiuPlesco/microservices-go/grpc/data"
	"github.com/SergiuPlesco/microservices-go/grpc/pb"
	"github.com/hashicorp/go-hclog"
)

type Currency struct {
	pb.UnimplementedCurrencyServer
	rates         *data.ExchangeRates
	log           hclog.Logger
	subscriptions map[pb.Currency_SubscribeRatesServer][]*pb.RateRequest
}

func NewCurrency(r *data.ExchangeRates, l hclog.Logger) *Currency {
	c := &Currency{
		rates:         r,
		log:           l,
		subscriptions: make(map[pb.Currency_SubscribeRatesServer][]*pb.RateRequest),
	}
	go c.handleUpdates()
	return c

}

func (c *Currency) handleUpdates() {
	ru := c.rates.MonitorRates(5 * time.Second)
	for range ru {
		c.log.Info("Got Updated rates")

		// loop over subscribed clients
		for k, v := range c.subscriptions {

			// loop over subscribed rates
			for _, rr := range v {
				r, err := c.rates.GetRate(rr.GetBase().String(), rr.GetDestination().String())
				if err != nil {
					c.log.Error("Unable to get update rate", "base", rr.GetBase().String(), "destination", rr.GetDestination().String())
				}

				err = k.Send(&pb.RateResponse{Base: rr.Base, Destination: rr.Destination, Rate: r})
				if err != nil {
					c.log.Error("Unable to send updated rate", "base", rr.GetBase().String(), "destination", rr.GetDestination().String())
				}
			}
		}

	}
}

func (c *Currency) GetRate(ctx context.Context, req *pb.RateRequest) (*pb.RateResponse, error) {
	c.log.Debug("Handle GetRate", "base", req.GetBase(), "destination", req.GetDestination())

	rate, err := c.rates.GetRate(req.GetBase().String(), req.GetDestination().String())
	if err != nil {
		return nil, err
	}

	return &pb.RateResponse{Rate: rate}, nil
}

func (c *Currency) SubscribeRates(src pb.Currency_SubscribeRatesServer) error {

	for {
		// Recv is a blocking method which returns on client data
		rr, err := src.Recv()
		// io.EOF signals that the client has closed the connection
		if err == io.EOF {
			c.log.Info("client has closed connection")
			break
		}
		// any other error means the transport between the server and client is unavailable
		if err != nil {
			c.log.Error("unable to read from client", "error", err)
			break
		}
		c.log.Info("handle client request", "request_base", rr.GetBase(), "request_dest", rr.GetDestination())

		rrs, ok := c.subscriptions[src]
		if !ok {
			rrs = []*pb.RateRequest{}
		}

		rrs = append(rrs, rr)
		c.subscriptions[src] = rrs
	}

	return nil
}
