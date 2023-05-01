package model

import (
	"fmt"
	"testing"
)

func TestString(t *testing.T) {
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
			Limit: 2,
			Multiples: []Multiple{
				{
					IntX: 3,
					StrX: "fizz",
				},
				{
					IntX: 5,
					StrX: "buzz",
				},
			},
		}, "{[{fizz 3} {buzz 5}] 2}"},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%s", tt.name)
		t.Run(testname, func(t *testing.T) {
			val := tt.i.String()
			if val != tt.want {
				t.Errorf("got %s, want %s", val, tt.want)
			}
		})
	}
}
