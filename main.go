package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/miya-masa/go-micro/apigateway"
	"github.com/miya-masa/go-micro/products"
	"github.com/miya-masa/go-micro/users"
)

func main() {

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stdout)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}
	r := mux.NewRouter()

	userproxy := apigateway.NewUserService([]string{
		"localhost:8080", "localhost:8081", "localhost:8082",
	})
	userproxy = users.Logging(logger, "localhost:8000")(userproxy)
	productproxy := apigateway.NewProductService([]string{
		"localhost:8083", "localhost:8084", "localhost:8085",
	})
	productproxy = products.Logging(logger, "localhost:8000")(productproxy)

	r.Handle("/users/{name}", httptransport.NewServer(
		users.NewEndpoints(userproxy).UserByName,
		decodeUserByNameRequest,
		encodeResponse))
	r.Handle("/products/{name}", httptransport.NewServer(
		products.NewEndpoints(productproxy).ProductByName,
		decodeProductByNameRequest,
		encodeResponse))

	serverApiGateway := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	usrv := users.NewService()
	userinstances := []string{"localhost:8080", "localhost:8081", "localhost:8082"}
	for _, v := range userinstances {
		r := mux.NewRouter()
		srv := users.Logging(logger, v)(usrv)
		r.Handle("/{name}", httptransport.NewServer(
			users.NewEndpoints(srv).UserByName,
			decodeUserByNameRequest,
			encodeResponse))
		s := &http.Server{
			Handler: r,
			Addr:    v,
			// Good practice: enforce timeouts for servers you create!
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		}
		go func(ins string) {
			logger.Log("start", ins)
			logger.Log(s.ListenAndServe())
		}(v)
	}

	psrv := products.NewService()
	prodinstances := []string{"localhost:8083", "localhost:8084", "localhost:8085"}

	for _, v := range prodinstances {
		r := mux.NewRouter()
		srv := products.Logging(logger, v)(psrv)
		r.Handle("/{name}", httptransport.NewServer(
			products.NewEndpoints(srv).ProductByName,
			decodeProductByNameRequest,
			encodeResponse))
		s := &http.Server{
			Handler: r,
			Addr:    v,
			// Good practice: enforce timeouts for servers you create!
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		}
		go func(ins string) {
			logger.Log("start", ins)
			logger.Log(s.ListenAndServe())
		}(v)
	}
	logger.Log(serverApiGateway.ListenAndServe())
}

func decodeProductByNameRequest(_ context.Context, r *http.Request) (interface{}, error) {
	if name, ok := mux.Vars(r)["name"]; ok {
		return name, nil
	}
	return nil, errors.New("name is required")
}

func decodeUserByNameRequest(_ context.Context, r *http.Request) (interface{}, error) {
	if name, ok := mux.Vars(r)["name"]; ok {
		return name, nil
	}
	return nil, errors.New("name is required")
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
