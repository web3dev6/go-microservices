package server

import (
	"context"

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

func NewCurrency(l hclog.Logger, er *data.ExchangeRates) *Currency {
	return &Currency{l, er, &pb.UnimplementedCurrencyServer{}}
}

func (c *Currency) GetRate(ctx context.Context, rr *pb.RateRequest) (*pb.RateResponse, error) {
	c.log.Info("Handle GetRate", "base", rr.GetBase(), "destination", rr.GetDestination())
	rate, err := c.rates.GetRate(rr.GetBase().String(), rr.GetDestination().String())
	if err != nil {
		return nil, err
	}

	return &pb.RateResponse{Rate: rate}, nil
}
