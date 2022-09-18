package main

import (
	"net"
	"os"

	"github.com/cloudflare/cfssl/log"
	"github.com/hashicorp/go-hclog"
	"github.com/satoshi-u/go-microservices/currency/pb"
	"github.com/satoshi-u/go-microservices/currency/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {

	// Register CurrencyServer with grpcServer & currencyServer instances
	gs := grpc.NewServer()
	cs := server.NewCurrency(hclog.Default())
	pb.RegisterCurrencyServer(gs, cs)

	// solution | Failed to list services: server does not support the reflection API
	reflection.Register(gs)

	// start a grpc server( grpcServer has a method Serve | similar to httpServer.ListenAndServe)
	listener, err := net.Listen("tcp", ":9092")
	if err != nil {
		log.Error("Unable to listen", "error", err)
		os.Exit(1)
	}
	gs.Serve(listener)

	/*
		grpcurl --plaintext localhost:9092 list
		grpcurl --plaintext localhost:9092 list pb.Currency
		grpcurl --plaintext localhost:9092 describe pb.Currency.GetRate
		grpcurl --plaintext localhost:9092 describe pb.RateRequest
		grpcurl --plaintext localhost:9092 describe pb.RateResponse

		-> when base & destination of type string
		grpcurl --plaintext -d '{"Base":"GBP", "Destination":"INR"}' localhost:9092 pb.Currency.GetRate
	*/
}
