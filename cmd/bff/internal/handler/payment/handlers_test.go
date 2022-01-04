package paymentHandlers_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	cartHandlers "github.com/vincentandr/shopping-microservice/cmd/bff/internal/handler/cart"
	catalogHandlers "github.com/vincentandr/shopping-microservice/cmd/bff/internal/handler/catalog"
	paymentHandlers "github.com/vincentandr/shopping-microservice/cmd/bff/internal/handler/payment"
	cartMock "github.com/vincentandr/shopping-microservice/cmd/bff/internal/mock/cart"
	catalogMock "github.com/vincentandr/shopping-microservice/cmd/bff/internal/mock/catalog"
	paymentMock "github.com/vincentandr/shopping-microservice/cmd/bff/internal/mock/payment"
	routes "github.com/vincentandr/shopping-microservice/cmd/bff/internal/route"
	paymentPb "github.com/vincentandr/shopping-microservice/internal/proto/payment"
)

// Test if the routes and status are correct

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