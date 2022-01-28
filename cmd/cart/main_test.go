package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vincentandr/shopping-microservice/cmd/cart/internal/server"
	pb "github.com/vincentandr/shopping-microservice/internal/proto/cart"
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
	productMaps := []map[string]map[string]string{
		{
		"itemid1": {"name": "itemname1", "price": "800", "desc":"desc1", "image":"image1"},
		"itemid2": {"name": "itemname2", "price": "45", "desc":"desc2", "image":"image2"},
		"itemid3": {"name": "itemname3", "price": "120", "desc":"desc3", "image":"image3"},
		},
		{
			"itemid1": {"name": "itemname1", "price": "800", "desc":"desc1", "image":"image1"},
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
					Name: productMaps[0]["itemid1"]["name"],
					Price: 800,
					Qty: 36,
					Desc: productMaps[0]["itemid1"]["desc"],
					Image: productMaps[0]["itemid1"]["image"],
				},
				{
					ProductId: "itemid2",
					Name: productMaps[0]["itemid2"]["name"],
					Price: 45,
					Qty: 8,
					Desc: productMaps[0]["itemid2"]["desc"],
					Image: productMaps[0]["itemid2"]["image"],
				},
				{
					ProductId: "itemid3",
					Name: productMaps[0]["itemid3"]["name"],
					Price: 120,
					Qty: 12,
					Desc: productMaps[0]["itemid3"]["desc"],
					Image: productMaps[0]["itemid3"]["image"],
				},
			},
		},
		{
			Products: []*pb.ItemResponse{
				{
					ProductId: "itemid1",
					Name: productMaps[1]["itemid1"]["name"],
					Price: 800,
					Qty: 36,
					Desc: productMaps[1]["itemid1"]["desc"],
					Image: productMaps[1]["itemid1"]["image"],
				},
			},
		},
		{},
	}

	tcs := []struct{
		name string
		products map[string]map[string]string
		cartItems map[string]string
		expected *pb.ItemsResponse
	}{
		{
			"append more than 1 items",
			productMaps[0],
			qtyMaps[0],
			&expectedSlices[0],
		},
		{
			"append just 1 item",
			productMaps[1],
			qtyMaps[1],
			&expectedSlices[1],
		},
		{
			"append no item",
			productMaps[2],
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