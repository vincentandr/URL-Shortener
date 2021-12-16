package bff

import (
	"testing"
)

func TestHome(t *testing.T){
	testcases := []struct{
		name string
		input int
		expected int
	}{
		{"positive", 1, 3},
		{"zero", 0, 2},
		{"negative", -1, 1},
		{"long", 9999, 10001},
	}

	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			if res, _ := Home(test.input); test.expected != res{
				t.Error("Calculation is not correct!")
			}
		})
	}
}