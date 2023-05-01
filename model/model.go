// This package contains golang main models for API endpoints (input/output)
package model

import (
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

// Input is '/v1/fizzbuzz' model endpoint
type Input struct {
	Limit     int        `validate:"gt=0,lte=1000000"`
	Multiples []Multiple `validate:"gte=0,lte=2,dive"`
}

type Multiple struct {
	IntX int    `validate:"gt=0"`
	StrX string `validate:"required"`
}

// Output is '/v1/statistics' model endpoint
type Output struct {
	Request Input `json:"request"`
	Hits    int   `json:"hits"`
}

// String returns data in String format
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
