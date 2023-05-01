// This package process data from action class handlers
package manager

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/arckadious/fizzbuzz/model"
	"github.com/arckadious/fizzbuzz/repository"
	"github.com/arckadious/fizzbuzz/response"
)

type Fizz struct {
	*Manager //Fizz class has attributes and methods from manager parent class
	repoFizz *repository.RepositoryFizz
}

// NewFizz constructor Fizz
func NewFizz(m *Manager, repo *repository.RepositoryFizz) *Fizz {
	return &Fizz{
		Manager:  m,
		repoFizz: repo,
	}
}

func (m *Fizz) HandleFizz(w http.ResponseWriter, input model.Input) {

	res := m.GetApiResponse()

	var tab []string
	for i := 1; i <= input.Limit; i++ {
		var elem string
		for _, multiple := range input.Multiples {
			if i%multiple.IntX == 0 {
				elem += multiple.StrX
			}
		}

		//no multiple -> add number
		if elem == "" {
			tab = append(tab, strconv.Itoa(i))
		} else {
			tab = append(tab, elem)
		}

	}

	output := strings.Join(tab, ",")

	res.SetData(output).WriteJSONResponse(w)
}

func (m *Fizz) HandleStatistics(w http.ResponseWriter) {

	res := m.GetApiResponse()

	msg, hits, noRows, err := m.repoFizz.GetMostRequestUsed()
	if err != nil {
		if noRows {
			res.SetStatusCode(http.StatusPartialContent).WriteJSONResponse(w)
			return
		}
		res.SetInternalServerErrorResponse([]response.ApiError{{Code: response.ErrorInternalServerError, Message: err.Error()}}).WriteJSONResponse(w)
		return
	}

	var msgStruct model.Input
	err = json.Unmarshal([]byte(msg), &msgStruct)
	if err != nil {
		res.SetInternalServerErrorResponse([]response.ApiError{{Code: response.ErrorInternalServerError, Message: err.Error()}}).WriteJSONResponse(w)
		return
	}

	res.SetData(model.Output{Hits: hits, Request: msgStruct}).WriteJSONResponse(w)
}
