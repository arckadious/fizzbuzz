package model

import (
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

type Input struct {
	Limit     int        `validate:"gt=0,lte=1000000"`
	Multiples []Multiple `validate:"gte=0,lte=2,dive"`
}

type Multiple struct {
	IntX int    `validate:"gt=0"`
	StrX string `validate:"required"`
}
type Output struct {
	Request Input `json:"request"`
	Hits    int   `json:"hits"`
}

func (input *Input) String() (val string) {
	if input == nil {
		logrus.Error("input.String() : struct nil")
		return
	}
	val += "{["
	for _, elem := range input.Multiples {
		val += "{" + elem.StrX + " " + strconv.Itoa(elem.IntX) + "} "
	}

	return strings.TrimRight(val, " ") + "] " + strconv.Itoa(input.Limit) + "}"
}
