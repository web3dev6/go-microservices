package server

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/satoshi-u/go-microservices/currency/data"
	"github.com/satoshi-u/go-microservices/currency/pb"
)

// implements CurrencyServer
type Currency struct {
	log   hclog.Logger
	rates *data.ExchangeRates
	*pb.UnimplementedCurrencyServer
}

// NewCurrency - gives back a currency server
func NewCurrency(l hclog.Logger, er *data.ExchangeRates) *Currency {
	return &Currency{l, er, &pb.UnimplementedCurrencyServer{}}
}

// GetRate - calls the underlying data.ExchangeRates with RateRequest values to get & return a valid RateResponse
func (c *Currency) GetRate(ctx context.Context, rr *pb.RateRequest) (*pb.RateResponse, error) {
	c.log.Info("Handle GetRate", "base", rr.GetBase(), "destination", rr.GetDestination())
	rate, err := c.rates.GetRate(rr.GetBase().String(), rr.GetDestination().String())
	if err != nil {
		return nil, err
	}
	log.Println("base rate: ", rr.GetBase().String())
	log.Println("destination rate: ", rr.GetDestination().String())
	log.Println("rate: ", rate)
	return &pb.RateResponse{Rate: rate}, nil
}

// SubscribeRates - starts sending const RateResponse in never ending loop to a client who calls -> GRPC pb.Currency.SubscribeRates
//								- starts receiving RateRequest in never ending loop to a client when client writes in stdin of called GRPC pb.Currency.SubscribeRates
func (c *Currency) SubscribeRates(src pb.Currency_SubscribeRatesServer) error {
	go func() {
		// inbound from client
		for {
			rr, err := src.Recv()
			if err == io.EOF {
				c.log.Info("Client has closed connection")
				break
			}
			if err != nil {
				c.log.Error("Unable to read from client", "error", err)
				break
			}
			time.Sleep(5 * time.Second)
			c.log.Info("Handle client request", "rate-request", rr)
		}
	}()
	// outbound to client
	for {
		err := src.Send(&pb.RateResponse{Rate: 1.0})
		if err != nil {
			return err
		}
		time.Sleep(5 * time.Second)
	}
}
