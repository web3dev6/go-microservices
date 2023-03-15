package server

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/satoshi-u/go-microservices/currency/data"
	"github.com/satoshi-u/go-microservices/currency/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// implements CurrencyServer
type Currency struct {
	log           hclog.Logger
	rates         *data.ExchangeRates
	subscriptions map[pb.Currency_SubscribeRatesServer][]*pb.RateRequest
	*pb.UnimplementedCurrencyServer
}

// NewCurrency - gives back a currency server
func NewCurrency(l hclog.Logger, er *data.ExchangeRates) *Currency {
	c := &Currency{l, er, make(map[pb.Currency_SubscribeRatesServer][]*pb.RateRequest), &pb.UnimplementedCurrencyServer{}}
	go c.handleUpdates()
	return c
}

// handleUpdates
func (c *Currency) handleUpdates() {
	ru := c.rates.MonitorRates(5 * time.Second)
	for range ru {
		c.log.Info("Got updated rates")
		// loop over subscribed clients
		for client, subscription := range c.subscriptions {
			// loop over rates for a specific client
			for _, rr := range subscription {
				r, err := c.rates.GetRate(rr.GetBase().String(), rr.GetDestination().String())
				if err != nil {
					c.log.Error("Unable to get updated rates", "base", rr.GetBase(), "destination", rr.GetDestination())
				}
				err = client.Send(&pb.RateResponse{Base: rr.GetBase(), Destination: rr.GetDestination(), Rate: r}) // @gRPC stream{server -> client}
				if err != nil {
					c.log.Error("Unable to send updated rates", "base", rr.GetBase(), "destination", rr.GetDestination())
				}
			}
		}
	}
}

// GetRate - calls the underlying data.ExchangeRates with RateRequest values to get & return a valid RateResponse
func (c *Currency) GetRate(ctx context.Context, rr *pb.RateRequest) (*pb.RateResponse, error) {
	c.log.Info("Handle GetRate", "base", rr.GetBase(), "destination", rr.GetDestination())

	// validation to learn - gRPC Error messages in Unary RPCs - at server side
	if rr.Base == rr.Destination {
		// return nil, fmt.Errorf("error : Base can not be the same as destination")
		// return nil, status.Errorf(
		// 	codes.InvalidArgument,
		// 	"Base currency %s can not be the same as destination currency %s",
		// 	rr.Base.String(),
		// 	rr.Destination.String(),
		// )
		err := status.Newf(
			codes.InvalidArgument,
			"Base currency %s can not be the same as destination currency %s",
			rr.Base.String(),
			rr.Destination.String(),
		)
		err, wdErr := err.WithDetails(rr)
		if wdErr != nil {
			return nil, wdErr
		}
		return nil, err.Err()
	}

	rate, err := c.rates.GetRate(rr.GetBase().String(), rr.GetDestination().String())
	if err != nil {
		return nil, err
	}
	log.Println("base rate: ", rr.GetBase().String())
	log.Println("destination rate: ", rr.GetDestination().String())
	log.Println("rate: ", rate)
	return &pb.RateResponse{Base: rr.Base, Destination: rr.Destination, Rate: rate}, nil
}

// SubscribeRates - starts sending const RateResponse in never ending loop to a client who calls -> GRPC pb.Currency.SubscribeRates
//								- starts receiving RateRequest in never ending loop to a client when client writes in stdin of called GRPC pb.Currency.SubscribeRates
func (c *Currency) SubscribeRates(src pb.Currency_SubscribeRatesServer) error {
	// inbound from client - handle client messages
	for {
		// Recv is a blocking method which returns on client data
		rr, err := src.Recv() // @gRPC stream{server <- client}
		if err == io.EOF {
			c.log.Info("Client has closed connection")
			break
		}
		// any other err - transport between the server and client is unavailable
		if err != nil {
			c.log.Error("Unable to read from client", "error", err)
			return err
		}

		// time.Sleep(5 * time.Second)
		c.log.Info("Handle client subscribe request", "rate-request", rr, "request_base", rr.GetBase(), "request_dest", rr.GetDestination())

		rrs, ok := c.subscriptions[src]
		if !ok {
			rrs = []*pb.RateRequest{}
		}
		rrs = append(rrs, rr)
		c.subscriptions[src] = rrs
	}
	// outbound to client - handle server responses
	// we block here to keep the connection open
	// for {
	// 	err := src.Send(&pb.RateResponse{Rate: 1.0})
	// 	if err != nil {
	// 		return err
	// 	}
	// 	time.Sleep(5 * time.Second)
	// }
	return nil
}
