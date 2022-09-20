package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-openapi/runtime/middleware"
	gorHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/nicholasjackson/env"
	"github.com/satoshi-u/go-microservices/currency/pb"
	"github.com/satoshi-u/go-microservices/product-api/data"
	"github.com/satoshi-u/go-microservices/product-api/handlers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// go mod init github.com/satoshi-u/go-microservices
// go run main.go
// hello   -> curl -v localhost:9090 -d sarthak
// GET     -> curl -v localhost:9090/products | jq
// GET     -> curl -v "localhost:9090/products?currency=INR" | jq
// GET     -> curl -v localhost:9090/products/2 | jq
// GET     -> curl -v "localhost:9090/products/2?currency=INR" | jq
// POST    -> curl -v localhost:9090/products -d '{"name": "Indian Tea", "description": "nice cup of tea", "price": 3.14, "sku": "prod-bev-003"}'| jq
// POST    -> curl -v localhost:9090/products -d '{"name": "coffee $1", "description": "cheap coffee", "price": 1.00, "sku": "prod-bev-004"}'| jq
// PUT   	 -> curl -v localhost:9090/products -XPUT -d '{"id": 1, "name": "Cappuccino", "description": "steamed milk foam", "price": 5.00, "sku": "prod-bev-001"}'| jq
// DELETE  -> curl -v localhost:9090/products/4 -XDELETE | jq

// create swagger.yaml       -> make swagger
// codegen from swagger.yaml -> swagger generate client -f ../swagger.yaml -A product-api

var bindAddress = env.String("BIND_ADDRESS", false, ":9090", "Bind address for the server")

func main() {

	env.Parse()

	// logger
	// l := log.New(os.Stdout, "product-api", log.LstdFlags)
	l := hclog.Default()
	// validation
	v := data.NewValidation()
	// client gRPC conn
	conn, err := grpc.Dial("localhost:9092", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	cc := pb.NewCurrencyClient(conn)
	// ProductsDB instance
	pdb := data.NewProductsDB(cc, l)
	// handler instantiate with constructor dependency injection : logger, validation, ProductsDB
	ph := handlers.NewProducts(l, v, pdb)
	// hh := handlers.NewHello(l)

	// new std lib mux : create mux and register handlers
	// sm := http.NewServeMux()
	// sm.Handle("/", hh)
	// sm.Handle("/", ph)

	// gorilla mux : create mux and register GET|POST|PUT|DELETE handlers
	sm := mux.NewRouter()
	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/products", ph.GetProducts)
	getRouter.HandleFunc("/products", ph.GetProducts).Queries("currency", "{[A-Z]{3}}")
	getRouter.HandleFunc("/products/{id:[0-9]+}", ph.GetProduct)
	getRouter.HandleFunc("/products/{id:[0-9]+}", ph.GetProduct).Queries("currency", "{[A-Z]{3}}")

	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/products", ph.UpdateProducts)
	putRouter.Use(ph.MiddlewareValidateProduct)

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/products", ph.AddProducts)
	postRouter.Use(ph.MiddlewareValidateProduct)

	deleteRouter := sm.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/products/{id:[0-9]+}", ph.DeleteProducts)

	// ReDocs- Swagger
	// make swagger
	// swagger generate client -f ../swagger.yaml -A product-api
	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(opts, nil)
	getRouter.HandleFunc("/docs", sh.ServeHTTP)
	getRouter.HandleFunc("/swagger.yaml", http.FileServer(http.Dir("./")).ServeHTTP)

	// CORS
	cors := gorHandlers.CORS(gorHandlers.AllowedOrigins([]string{"http://localhost:3000"})) // "http://localhost:3000"   *

	// new server- address, handler, tls, timeouts
	s := &http.Server{
		Addr:         *bindAddress,
		Handler:      cors(sm),
		ErrorLog:     l.StandardLogger(&hclog.StandardLoggerOptions{}),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// start server listen as a non-blocking separate go routine
	go func() {
		l.Info("Started server on port 9090...")
		err := s.ListenAndServe()
		if err != nil {
			l.Error("error starting server", "error", err)
			os.Exit(1)
		}
	}()

	// graceful shutdown with os signal -> set signal notification on our sig channel
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// graceful shutdown when recieve input in sigChan, blocking operation
	sig := <-sigChan
	l.Info("Recieved terminate in sigChan, initiating graceful shutdown... sig:", sig)
	tc, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	// Even though ctx will be expired, it is good practice to call its
	// cancellation function in any case. Failure to do so may keep the
	// context and its parent alive longer than necessary.
	defer cancel()
	s.Shutdown(tc)
}
