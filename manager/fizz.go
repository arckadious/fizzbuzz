// This package process data from action class handlers
package manager

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	cst "github.com/arckadious/fizzbuzz/constant"

	"github.com/arckadious/fizzbuzz/model"
	"github.com/arckadious/fizzbuzz/repository"
	"github.com/arckadious/fizzbuzz/response"
)

// Fizz Class (manager child class)
type Fizz struct {
	*Manager //Fizz class has attributes and methods from manager parent class
	repoFizz *repository.Fizz
}

// NewFizz constructor Manager child Fizz
func NewFizz(m *Manager, repo *repository.Fizz) *Fizz {
	return &Fizz{
		Manager:  m,
		repoFizz: repo,
	}
}

// HandleFizz writes to 'w' a list of strings with numbers from 1 to limit, where: all multiples specified are replaced by text
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

// HandleStatistics write to 'w' the parameters corresponding to the most used request, as well as the number of hits for this request.
func (m *Fizz) HandleStatistics(w http.ResponseWriter) {

	res := m.GetApiResponse()

	msg, hits, noRows, err := m.repoFizz.GetMostRequestUsed()
	if err != nil {
		if noRows { // Database does not contain rows with checksum not empty
			res.StatusCode = http.StatusPartialContent
			res.WriteJSONResponse(w)
			return
		}
		res.SetErrorResponse(http.StatusInternalServerError, []response.ApiError{{Code: cst.ErrorInternalServerError, Message: err.Error()}}).WriteJSONResponse(w)
		return
	}

	var msgStruct model.Input
	err = json.Unmarshal([]byte(msg), &msgStruct)
	if err != nil {
		res.SetErrorResponse(http.StatusInternalServerError, []response.ApiError{{Code: cst.ErrorInternalServerError, Message: err.Error()}}).WriteJSONResponse(w)
		return
	}

	res.SetData(model.Output{Hits: hits, Request: msgStruct}).WriteJSONResponse(w)
}
