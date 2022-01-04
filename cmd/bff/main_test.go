package main_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	cartMock "github.com/vincentandr/shopping-microservice/cmd/bff/internal/mock/cart"
	catalogMock "github.com/vincentandr/shopping-microservice/cmd/bff/internal/mock/catalog"
	paymentMock "github.com/vincentandr/shopping-microservice/cmd/bff/internal/mock/payment"
	cartHandlers "github.com/vincentandr/shopping-microservice/cmd/bff/internal/web/handler/cart"
	catalogHandlers "github.com/vincentandr/shopping-microservice/cmd/bff/internal/web/handler/catalog"
	paymentHandlers "github.com/vincentandr/shopping-microservice/cmd/bff/internal/web/handler/payment"
	routes "github.com/vincentandr/shopping-microservice/cmd/bff/internal/web/route"
	cartPb "github.com/vincentandr/shopping-microservice/internal/proto/cart"
	catalogPb "github.com/vincentandr/shopping-microservice/internal/proto/catalog"
	paymentPb "github.com/vincentandr/shopping-microservice/internal/proto/payment"
)

// Test if the routes and status are correct

// Cart Handlers Tests

func TestGetCartItems(t *testing.T) {
	input := &cartPb.GetCartItemsRequest{UserId: "user1"}
	expected := &cartPb.ItemsResponse{
		Products: []*cartPb.ItemResponse{
				{
					ProductId: "productid1",
					Name: "item1",
					Price: 350,
					Qty: 12,
				},
				{
					ProductId: "productid2",
					Name: "item2",
					Price: 120,
					Qty: 7,
				},
			},
	}
	httpMethod := http.MethodGet
	mockMethod := "Grpc_GetCartItems"
	route := "/cart/user1"

	var res *cartPb.ItemsResponse

	// Config mock
	m := new(cartMock.GrpcMock)

	m.On(mockMethod, mock.Anything, input, mock.Anything).Return(expected, nil)

	// Handlers for routers
	handlerMap := map[string]interface{}{
		"cart": cartHandlers.NewGrpcClient(m),
		"catalog": catalogHandlers.NewGrpcClient(new(catalogMock.GrpcMock)),
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

func TestAddOrUpdateCartQty(t *testing.T) {
	input := &cartPb.AddOrUpdateCartRequest{
		UserId: "user1", 
		ProductId: "productid1", 
		NewQty: 12,
	}
	expected := &cartPb.ItemsResponse{
		Products: []*cartPb.ItemResponse{
				{
					ProductId: "productid1",
					Name: "item1",
					Price: 350,
					Qty: 18,
				},
			},
	}
	httpMethod := http.MethodPut
	mockMethod := "Grpc_AddOrUpdateCart"
	route := "/cart/user1/productid1?qty=12"

	var res *cartPb.ItemsResponse

	// Config mock
	m := new(cartMock.GrpcMock)

	m.On(mockMethod, mock.Anything, input, mock.Anything).Return(expected, nil)

	// Handlers for routers
	handlerMap := map[string]interface{}{
		"cart": cartHandlers.NewGrpcClient(m),
		"catalog": catalogHandlers.NewGrpcClient(new(catalogMock.GrpcMock)),
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

func TestRemoveCartItem(t *testing.T) {
	input := &cartPb.RemoveItemFromCartRequest{
		UserId: "user1", 
		ProductId: "productid1",
	}
	expected := &cartPb.ItemsResponse{
		Products: []*cartPb.ItemResponse{
				{
					ProductId: "productid2",
					Name: "item2",
					Price: 120,
					Qty: 7,
				},
			},
	}
	httpMethod := http.MethodDelete
	mockMethod := "Grpc_RemoveItemFromCart"
	route := "/cart/user1/productid1"

	var res *cartPb.ItemsResponse

	// Config mock
	m := new(cartMock.GrpcMock)

	m.On(mockMethod, mock.Anything, input, mock.Anything).Return(expected, nil)

	// Handlers for routers
	handlerMap := map[string]interface{}{
		"cart": cartHandlers.NewGrpcClient(m),
		"catalog": catalogHandlers.NewGrpcClient(new(catalogMock.GrpcMock)),
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

func TestRemoveAllCartItems(t *testing.T) {
	input := &cartPb.RemoveAllCartItemsRequest{UserId: "user1",}
	expected := &cartPb.ItemsResponse{}
	httpMethod := http.MethodDelete
	mockMethod := "Grpc_RemoveAllCartItems"
	route := "/cart/user1"

	var res *cartPb.ItemsResponse

	// Config mock
	m := new(cartMock.GrpcMock)

	m.On(mockMethod, mock.Anything, input, mock.Anything).Return(expected, nil)

	// Handlers for routers
	handlerMap := map[string]interface{}{
		"cart": cartHandlers.NewGrpcClient(m),
		"catalog": catalogHandlers.NewGrpcClient(new(catalogMock.GrpcMock)),
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
		assert.ElementsMatch(t, res, expected)
	}
}

func TestCheckout(t *testing.T) {
	input := &cartPb.CheckoutRequest{UserId: "user1",}
	expected := &cartPb.CheckoutResponse{OrderId: "orderid1"}
	httpMethod := http.MethodGet
	mockMethod := "Grpc_Checkout"
	route := "/cart/checkout/user1"

	var res *cartPb.CheckoutResponse

	// Config mock
	m := new(cartMock.GrpcMock)

	m.On(mockMethod, mock.Anything, input, mock.Anything).Return(expected, nil)

	// Handlers for routers
	handlerMap := map[string]interface{}{
		"cart": cartHandlers.NewGrpcClient(m),
		"catalog": catalogHandlers.NewGrpcClient(new(catalogMock.GrpcMock)),
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
		assert.Equal(t, res, expected)
	}
}

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

// Payment Handlers Tests

func TestGetOrders(t *testing.T) {
	input := &paymentPb.GetOrdersRequest{UserId: ""}
	expected := &paymentPb.GetOrdersResponse{
		Orders: []*paymentPb.GetOrderResponse{
				{
					OrderId: "orderid1",
					UserId: "user1",
					Items: []*paymentPb.ItemResponse{
						{
							ProductId: "productid1",
							Name: "item1",
							Price: 350,
							Qty: 7,
						},
						{
							ProductId: "productid2",
							Name: "item2",
							Price: 150,
							Qty: 10,
						},
					},
					Status: "paid",
				},
			},
	}
	httpMethod := http.MethodGet
	mockMethod := "Grpc_GetOrders"
	route := "/payment"

	var res *paymentPb.GetOrdersResponse

	// Config mock
	m := new(paymentMock.GrpcMock)

	m.On(mockMethod, mock.Anything, input, mock.Anything).Return(expected, nil)

	// Handlers for routers
	handlerMap := map[string]interface{}{
		"cart": cartHandlers.NewGrpcClient(new(cartMock.GrpcMock)),
		"catalog": catalogHandlers.NewGrpcClient(new(catalogMock.GrpcMock)),
		"payment": paymentHandlers.NewGrpcClient(m),
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
		assert.ElementsMatch(t, res.Orders, expected.Orders)
	}
}

func TestGetOrdersWithUserId(t *testing.T) {
	input := &paymentPb.GetOrdersRequest{UserId: "user1"}
	expected := &paymentPb.GetOrdersResponse{
		Orders: []*paymentPb.GetOrderResponse{
				{
					OrderId: "orderid1",
					UserId: "user1",
					Items: []*paymentPb.ItemResponse{
						{
							ProductId: "productid1",
							Name: "item1",
							Price: 350,
							Qty: 7,
						},
						{
							ProductId: "productid2",
							Name: "item2",
							Price: 150,
							Qty: 10,
						},
					},
					Status: "paid",
				},
			},
	}
	httpMethod := http.MethodGet
	mockMethod := "Grpc_GetOrders"
	route := "/payment/user1"

	var res *paymentPb.GetOrdersResponse

	// Config mock
	m := new(paymentMock.GrpcMock)

	m.On(mockMethod, mock.Anything, input, mock.Anything).Return(expected, nil)

	// Handlers for routers
	handlerMap := map[string]interface{}{
		"cart": cartHandlers.NewGrpcClient(new(cartMock.GrpcMock)),
		"catalog": catalogHandlers.NewGrpcClient(new(catalogMock.GrpcMock)),
		"payment": paymentHandlers.NewGrpcClient(m),
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
		assert.ElementsMatch(t, res.Orders, expected.Orders)
	}
}

func TestMakePayment(t *testing.T) {
	input := &paymentPb.PaymentRequest{OrderId: "orderid1"}
	expected := &paymentPb.PaymentResponse{}
	httpMethod := http.MethodPut
	mockMethod := "Grpc_MakePayment"
	route := "/payment/orderid1"

	var res *paymentPb.PaymentResponse

	// Config mock
	m := new(paymentMock.GrpcMock)

	m.On(mockMethod, mock.Anything, input, mock.Anything).Return(expected, nil)

	// Handlers for routers
	handlerMap := map[string]interface{}{
		"cart": cartHandlers.NewGrpcClient(new(cartMock.GrpcMock)),
		"catalog": catalogHandlers.NewGrpcClient(new(catalogMock.GrpcMock)),
		"payment": paymentHandlers.NewGrpcClient(m),
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
		assert.ElementsMatch(t, res, expected)
	}
}