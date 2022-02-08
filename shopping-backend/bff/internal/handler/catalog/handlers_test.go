package catalogHandlers_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	cartHandlers "github.com/vincentandr/shopping-microservice/bff/internal/handler/cart"
	catalogHandlers "github.com/vincentandr/shopping-microservice/bff/internal/handler/catalog"
	paymentHandlers "github.com/vincentandr/shopping-microservice/bff/internal/handler/payment"
	cartMock "github.com/vincentandr/shopping-microservice/bff/internal/mock/cart"
	catalogMock "github.com/vincentandr/shopping-microservice/bff/internal/mock/catalog"
	paymentMock "github.com/vincentandr/shopping-microservice/bff/internal/mock/payment"
	routes "github.com/vincentandr/shopping-microservice/bff/internal/route"
	catalogPb "github.com/vincentandr/shopping-microservice/internal/proto/catalog"
)

// Test if the routes and status are correct

// Catalog Handlers Tests

func TestGetProducts(t *testing.T) {
	input := &catalogPb.EmptyRequest{}
	expected := &catalogPb.GetProductsResponse{
		Products: []*catalogPb.GetProductResponse{
				{
					ProductId: "productid1",
					Name: "item1",
					Price: 120,
					Qty: 7,
				},
				{
					ProductId: "productid2",
					Name: "item2",
					Price: 350,
					Qty: 10,
				},
			},
	}
	httpMethod := http.MethodGet
	mockMethod := "Grpc_GetProducts"
	route := "/products"

	var res *catalogPb.GetProductsResponse

	// Config mock
	m := new(catalogMock.GrpcMock)

	m.On(mockMethod, mock.Anything, input, mock.Anything).Return(expected, nil)

	// Handlers for routers
	handlerMap := map[string]interface{}{
		"cart": cartHandlers.NewGrpcClient(new(cartMock.GrpcMock)),
		"catalog": catalogHandlers.NewGrpcClient(m),
		"payment": paymentHandlers.NewGrpcClient(new(paymentMock.GrpcMock)),
	}

	r := routes.NewRouter(handlerMap)

	server := httptest.NewServer(r)
	defer server.Close()

	// Config request
	url := server.URL + route

	req, err := http.NewRequestWithContext(context.Background(), httpMethod, url, nil)
	if err != nil {
		t.Fatalf("could not create %s request: %v", httpMethod, err)
	}
	req.Close = true

	// Send request
	client := &http.Client{}
	
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("could not send %s request: %v", httpMethod, err)
	}
	defer resp.Body.Close()

	// Check for http status match
	if assert.Equal(t, http.StatusOK, resp.StatusCode) {
		// Decode response
		err = json.NewDecoder(resp.Body).Decode(&res)
		if err != nil {
			t.Fatalf("failed to decode %s response: %v", httpMethod, err)
		}
		
		// Check for result struct match
		assert.ElementsMatch(t, res.Products, expected.Products)
	}
}

func TestGetProductsWithError(t *testing.T) {
	input := &catalogPb.EmptyRequest{}
	expected := http.StatusInternalServerError
	httpMethod := http.MethodGet
	mockMethod := "Grpc_GetProducts"
	route := "/products"

	// Config mock
	m := new(catalogMock.GrpcMock)

	m.On(mockMethod, mock.Anything, input, mock.Anything).Return(nil, errors.New("error"))

	// Handlers for routers
	handlerMap := map[string]interface{}{
		"cart": cartHandlers.NewGrpcClient(new(cartMock.GrpcMock)),
		"catalog": catalogHandlers.NewGrpcClient(m),
		"payment": paymentHandlers.NewGrpcClient(new(paymentMock.GrpcMock)),
	}

	r := routes.NewRouter(handlerMap)

	server := httptest.NewServer(r)
	defer server.Close()

	// Config request
	url := server.URL + route

	req, err := http.NewRequestWithContext(context.Background(), httpMethod, url, nil)
	if err != nil {
		t.Fatalf("could not create %s request: %v", httpMethod, err)
	}
	req.Close = true

	// Send request
	client := &http.Client{}
	
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("could not send %s request: %v", httpMethod, err)
	}
	defer resp.Body.Close()

	// Check for http status match
	assert.Equal(t, expected, resp.StatusCode)
}

func TestGetProductsByName(t *testing.T) {
	input := &catalogPb.GetProductsByNameRequest{Name: "item2"}
	expected := &catalogPb.GetProductsResponse{
		Products: []*catalogPb.GetProductResponse{
				{
					ProductId: "productid2",
					Name: "item2",
					Price: 350,
					Qty: 10,
				},
			},
	}
	httpMethod := http.MethodGet
	mockMethod := "Grpc_GetProductsByName"
	route := "/products/search?name=item2"

	var res *catalogPb.GetProductsResponse

	// Config mock
	m := new(catalogMock.GrpcMock)

	m.On(mockMethod, mock.Anything, input, mock.Anything).Return(expected, nil)

	// Handlers for routers
	handlerMap := map[string]interface{}{
		"cart": cartHandlers.NewGrpcClient(new(cartMock.GrpcMock)),
		"catalog": catalogHandlers.NewGrpcClient(m),
		"payment": paymentHandlers.NewGrpcClient(new(paymentMock.GrpcMock)),
	}

	r := routes.NewRouter(handlerMap)

	server := httptest.NewServer(r)
	defer server.Close()

	// Config request
	url := server.URL + route

	req, err := http.NewRequestWithContext(context.Background(), httpMethod, url, nil)
	if err != nil {
		t.Fatalf("could not create %s request: %v", httpMethod, err)
	}
	req.Close = true

	// Send request
	client := &http.Client{}
	
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("could not send %s request: %v", httpMethod, err)
	}
	defer resp.Body.Close()

	// Check for http status match
	if assert.Equal(t, http.StatusOK, resp.StatusCode) {
		// Decode response
		err = json.NewDecoder(resp.Body).Decode(&res)
		if err != nil {
			t.Fatalf("failed to decode %s response: %v", httpMethod, err)
		}
		
		// Check for result struct match
		assert.ElementsMatch(t, res.Products, expected.Products)
	}
}

func TestGetProductsByNameWithError(t *testing.T) {
	input := &catalogPb.GetProductsByNameRequest{Name: "item2"}
	expected := http.StatusInternalServerError
	httpMethod := http.MethodGet
	mockMethod := "Grpc_GetProductsByName"
	route := "/products/search?name=item2"

	// Config mock
	m := new(catalogMock.GrpcMock)

	m.On(mockMethod, mock.Anything, input, mock.Anything).Return(nil, errors.New("error"))

	// Handlers for routers
	handlerMap := map[string]interface{}{
		"cart": cartHandlers.NewGrpcClient(new(cartMock.GrpcMock)),
		"catalog": catalogHandlers.NewGrpcClient(m),
		"payment": paymentHandlers.NewGrpcClient(new(paymentMock.GrpcMock)),
	}

	r := routes.NewRouter(handlerMap)

	server := httptest.NewServer(r)
	defer server.Close()

	// Config request
	url := server.URL + route

	req, err := http.NewRequestWithContext(context.Background(), httpMethod, url, nil)
	if err != nil {
		t.Fatalf("could not create %s request: %v", httpMethod, err)
	}
	req.Close = true

	// Send request
	client := &http.Client{}
	
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("could not send %s request: %v", httpMethod, err)
	}
	defer resp.Body.Close()

	// Check for http status match
	assert.Equal(t, expected, resp.StatusCode)
}