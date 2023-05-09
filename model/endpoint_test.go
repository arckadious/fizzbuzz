// This package contains golang main models for API endpoints (input/output)
package model

import (
	"testing"
)

func TestEndpoint(t *testing.T) {

	////////////////////
	// Input.String() //
	////////////////////

	var tests = []struct {
		name string
		i    *Input
		want string
	}{
		{"Struct nil", nil, ""},
		{"Struct empty", &Input{}, "{[] 0}"},
		{"Struct tab size 0", &Input{Limit: 2}, "{[] 2}"},
		{"Struct tab size 1", &Input{
			Limit: 2,
			Multiples: []Multiple{
				{
					IntX: 3,
					StrX: "fizz",
				},
			},
		}, "{[{fizz 3}] 2}"},
		{"Struct tab size 2", &Input{
			Multiples: []Multiple{
				{
					StrX: "fizz",
					IntX: 3,
				},
				{
					IntX: 5,
					StrX: "buzz",
				},
			},
			Limit: 2,
		}, "{[{fizz 3} {buzz 5}] 2}"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val := tt.i.String()
			if val != tt.want {
				t.Errorf("got %s, want %s. test name : %s", val, tt.want, tt.name)
			}
		})
	}
}
