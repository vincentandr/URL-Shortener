package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
			res := GetMapKeys(tc.items)
			assert.ElementsMatch(t, res, tc.expected)
		})
	}
}

// Append quantity purchased to the product details and return it as response
func TestAppendItemToResponse(t *testing.T) {
	productSlices := []catalogpb.GetProductsByIdsResponse{
		{
			Products: []*catalogpb.GetProductByIdsResponse{
				{ProductId: "itemid1", Name: "itemname1", Price: 800},
				{ProductId: "itemid2", Name: "itemname2", Price: 45},
				{ProductId: "itemid3", Name: "itemname3", Price: 120},
			},
		},
		{
			Products: []*catalogpb.GetProductByIdsResponse{
				{ProductId: "itemid1", Name: "itemname1", Price: 800},
			},
		},
		{},
	}

	maps := []map[string]string{
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
					ProductId: productSlices[0].Products[0].ProductId,
					Name: productSlices[0].Products[0].Name,
					Price: productSlices[0].Products[0].Price,
					Qty: 36,
				},
				{
					ProductId: productSlices[0].Products[1].ProductId,
					Name: productSlices[0].Products[1].Name,
					Price: productSlices[0].Products[1].Price,
					Qty: 8,
				},
				{
					ProductId: productSlices[0].Products[2].ProductId,
					Name: productSlices[0].Products[2].Name,
					Price: productSlices[0].Products[2].Price,
					Qty: 12,
				},
			},
		},
		{
			Products: []*pb.ItemResponse{
				{
					ProductId: productSlices[1].Products[0].ProductId,
					Name: productSlices[1].Products[0].Name,
					Price: productSlices[1].Products[0].Price,
					Qty: 36,
				},
			},
		},
		{},
	}

	tcs := []struct{
		name string
		products *catalogpb.GetProductsByIdsResponse
		cartItems map[string]string
		expected *pb.ItemsResponse
	}{
		{
			"append more than 1 items",
			&productSlices[0],
			maps[0],
			&expectedSlices[0],
		},
		{
			"append just 1 item",
			&productSlices[1],
			maps[1],
			&expectedSlices[1],
		},
		{
			"append no item",
			&productSlices[2],
			maps[2],
			&expectedSlices[2],
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			res, err := AppendItemToResponse(tc.products, tc.cartItems)
			if assert.NoError(t, err) {
				assert.ElementsMatch(t, res.Products, tc.expected.Products)
			}
		})
	}
}