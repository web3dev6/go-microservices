package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-openapi/runtime/middleware"
	gorHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
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
// GET/{id}-> curl -v localhost:9090/products/2 | jq
// POST    -> curl -v localhost:9090/products -d '{"name": "Indian Tea", "description": "nice cup of tea", "price": 3.14, "sku": "prod-bev-003"}'| jq
// POST    -> curl -v localhost:9090/products -d '{"name": "coffee $1", "description": "cheap coffee", "price": 1.00, "sku": "prod-bev-004"}'| jq
// PUT   	 -> curl -v localhost:9090/products -XPUT -d '{"id": 1, "name": "Cappuccino", "description": "steamed milk foam", "price": 5.00, "sku": "prod-bev-001"}'| jq
// DELETE  -> curl -v localhost:9090/products/4 -XDELETE | jq

// create swagger.yaml       -> make swagger
// codegen from swagger.yaml -> swagger generate client -f ../swagger.yaml -A product-api
func main() {
	// logger
	l := log.New(os.Stdout, "product-api", log.LstdFlags)
	// validation
	v := data.NewValidation()
	// client gRPC conn
	conn, err := grpc.Dial("localhost:9092", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	cc := pb.NewCurrencyClient(conn)
	// handler instantiate with constructor dependency injection : logger, validation, conn
	ph := handlers.NewProducts(l, v, cc)
	// hh := handlers.NewHello(l)

	// new std lib mux : create mux and register handlers
	// sm := http.NewServeMux()
	// sm.Handle("/", hh)
	// sm.Handle("/", ph)

	// gorilla mux : create mux and register GET|POST|PUT|DELETE handlers
	sm := mux.NewRouter()
	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/products", ph.GetProducts)
	getRouter.HandleFunc("/products/{id}", ph.GetProduct)

	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/products", ph.UpdateProducts)
	putRouter.Use(ph.MiddlewareValidateProduct)

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/products", ph.AddProducts)
	postRouter.Use(ph.MiddlewareValidateProduct)

	deleteRouter := sm.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/products/{id}", ph.DeleteProducts)

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
		Addr:         ":9090",
		Handler:      cors(sm),
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	// start server listen as a non-blocking separate go routine
	go func() {
		log.Printf("Started http server at 9090...")
		err := s.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// graceful shutdown with os signal -> set signal notification on our sig channel
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// graceful shutdown when recieve input in sigChan, blocking operation
	sig := <-sigChan
	log.Println("Recieved terminate in sigChan, initiating graceful shutdown... sig:", sig)
	tc, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	// Even though ctx will be expired, it is good practice to call its
	// cancellation function in any case. Failure to do so may keep the
	// context and its parent alive longer than necessary.
	defer cancel()
	s.Shutdown(tc)
}
