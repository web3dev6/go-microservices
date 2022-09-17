package server

import (
	"context"

	"github.com/hashicorp/go-hclog"
	"github.com/satoshi-u/go-microservices/currency/pb"
)

// implements CurrencyServer
type Currency struct {
	log hclog.Logger
	*pb.UnimplementedCurrencyServer
}

func NewCurrency(l hclog.Logger) *Currency {
	return &Currency{l, &pb.UnimplementedCurrencyServer{}}
}

func (c *Currency) GetRate(ctx context.Context, rr *pb.RateRequest) (*pb.RateResponse, error) {
	c.log.Info("Handle GetRate", "base", rr.GetBase(), "destination", rr.GetDestination())

	return &pb.RateResponse{Rate: 0.5}, nil
}
