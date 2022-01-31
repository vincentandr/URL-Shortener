package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vincentandr/shopping-microservice/cmd/cart/internal/server"
	pb "github.com/vincentandr/shopping-microservice/internal/proto/cart"
	catalogpb "github.com/vincentandr/shopping-microservice/internal/proto/catalog"
)

// Test to get item ids from cart items map
func TestGetMapKeys(t *testing.T) {
	tcs := []struct{
		name string
		items map[string]string
		expected []string
	}{
		{
			"more than 1 items",
			map[string]string{
				"itemId1":"4",
				"itemId2":"2",
				"itemId3":"9",
			},
			[]string{"itemId1", "itemId2", "itemId3"},
		},
		{
			"just 1 item",
			map[string]string{
				"itemId1":"4",
			},
			[]string{"itemId1"},
		},
		{
			"no item",
			map[string]string{},
			[]string{},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			res := server.GetMapKeys(tc.items)
			assert.ElementsMatch(t, res, tc.expected)
		})
	}
}

// Append quantity purchased to the product details and return it as response
func TestAppendItemToResponse(t *testing.T) {
	productSlices := []catalogpb.GetProductsResponse{
		{
			Products: []*catalogpb.GetProductResponse{
				{ProductId: "itemid1", Name: "itemname1", Price: 800, Qty: 50, Desc: "desc1", Image: "image1"},
				{ProductId: "itemid2", Name: "itemname2", Price: 45, Qty: 75, Desc: "desc1", Image: "image1"},
				{ProductId: "itemid3", Name: "itemname3", Price: 120, Qty: 100, Desc: "desc1", Image: "image1"},
			},
		},
		{
			Products: []*catalogpb.GetProductResponse{
				{ProductId: "itemid1", Name: "itemname1", Price: 800, Qty: 50, Desc: "desc1", Image: "image1"},
			},
		},
		{},
	}

	qtyMaps := []map[string]string{
		{
			"itemid1": "36",
			"itemid2": "8",
			"itemid3": "12",
		},
		{
			"itemid1": "36",
		},
		{},
	}

	expectedSlices := []pb.ItemsResponse{
		{
			Products: []*pb.ItemResponse{
				{
					ProductId: "itemid1",
					Name: productSlices[0].Products[0].Name,
					Price: 800,
					Qty: 36,
					Stock: 50,
					Desc: productSlices[0].Products[0].Desc,
					Image: productSlices[0].Products[0].Image,
				},
				{
					ProductId: "itemid2",
					Name: productSlices[0].Products[1].Name,
					Price: 45,
					Qty: 8,
					Stock: 75,
					Desc: productSlices[0].Products[1].Desc,
					Image: productSlices[0].Products[1].Image,
				},
				{
					ProductId: "itemid3",
					Name: productSlices[0].Products[2].Name,
					Price: 120,
					Qty: 12,
					Stock: 100,
					Desc: productSlices[0].Products[2].Desc,
					Image: productSlices[0].Products[2].Image,
				},
			},
		},
		{
			Products: []*pb.ItemResponse{
				{
					ProductId: "itemid1",
					Name: productSlices[1].Products[0].Name,
					Price: 800,
					Qty: 36,
					Stock: 50,
					Desc: productSlices[1].Products[0].Desc,
					Image: productSlices[1].Products[0].Image,
				},
			},
		},
		{},
	}

	tcs := []struct{
		name string
		products *catalogpb.GetProductsResponse
		cartItems map[string]string
		expected *pb.ItemsResponse
	}{
		{
			"append more than 1 items",
			&productSlices[0],
			qtyMaps[0],
			&expectedSlices[0],
		},
		{
			"append just 1 item",
			&productSlices[1],
			qtyMaps[1],
			&expectedSlices[1],
		},
		{
			"append no item",
			&productSlices[2],
			qtyMaps[2],
			&expectedSlices[2],
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			res, err := server.AppendItemToResponse(tc.products, tc.cartItems)
			if assert.NoError(t, err) {
				assert.ElementsMatch(t, res.Products, tc.expected.Products)
			}
		})
	}
}