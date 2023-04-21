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

// Manager fizz
type Fizz struct {
	*Manager //Fizz class has attributes and methods from manager parent class
	repoFizz *repository.RepositoryFizz
}

// NewFizz Constructor
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

	res.SetData(output).WriteJsonResponse(w)
}

func (m *Fizz) HandleStatistics(w http.ResponseWriter) {

	res := m.GetApiResponse()

	msg, hits, noRows, err := m.repoFizz.GetMostRequestUsed()
	if err != nil {
		if noRows {
			res.SetStatusCode(http.StatusPartialContent).WriteJsonResponse(w)
			return
		}
		res.SetInternalServerErrorResponse([]response.ApiError{{Code: response.ErrorInternalServerError, Message: err.Error()}}).WriteJsonResponse(w)
		return
	}

	var msgStruct model.Input
	err = json.Unmarshal([]byte(msg), &msgStruct)
	if err != nil {
		res.SetInternalServerErrorResponse([]response.ApiError{{Code: response.ErrorInternalServerError, Message: err.Error()}}).WriteJsonResponse(w)
		return
	}

	//TO DO
	res.SetData(model.Output{Hits: hits, Request: msgStruct}).WriteJsonResponse(w)
}
