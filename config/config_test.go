package config

import (
	"reflect"
	"testing"

	"github.com/go-playground/validator/v10"
)

func TestNew(t *testing.T) {
	type args struct {
		fileName string
	}
	tests := []struct {
		name string
		args args
		want *Config
	}{
		{
			"Load config",
			args{
				"../parameters/parameters.json",
			},
			New("../parameters/parameters.json", *validator.New()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.fileName, *validator.New()); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
