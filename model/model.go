package model

type Input struct {
	Multiples []struct {
		StrX string `validate:"required"`
		IntX int    `validate:"gt=0"`
	} `validate:"gte=0,lte=2,dive"`
	Limit int `validate:"gt=0,lte=1000000"`
}
