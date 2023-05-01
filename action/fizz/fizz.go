// This package extract and validate data
package fizz

import (
	"encoding/json"
	"net/http"

	"github.com/arckadious/fizzbuzz/manager"
	"github.com/arckadious/fizzbuzz/model"
	"github.com/arckadious/fizzbuzz/response"
)

type FizzAction struct {
	mng *manager.Fizz
}

// NewFizzAction constructor FizzAction
func NewFizzAction(mng *manager.Fizz) *FizzAction {
	return &FizzAction{mng}
}

// HandleFizz extract and validate data
func (ac *FizzAction) HandleFizz(w http.ResponseWriter, r *http.Request) {

	res := ac.mng.GetApiResponse() //generic response, unchanged, send 200 OK with json body response
	var input model.Input

	//extract input
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		res.SetBadRequestResponse([]response.ApiError{{Code: response.ErrorInvalidData, Message: "Data JSON input bad format."}}).WriteJSONResponse(w)
		return
	}

	//Validate input
	if err := ac.mng.GetValidator().Struct(input); err != nil {
		res.SetBadRequestResponse([]response.ApiError{{Code: response.ErrorInvalidField, Message: err.Error()}}).WriteJSONResponse(w)
		return
	}

	ac.mng.HandleFizz(w, input)

}

// HandleStatistics send back statistics about what the most used request has been
func (ac *FizzAction) HandleStatistics(w http.ResponseWriter, r *http.Request) {
	ac.mng.HandleStatistics(w)
}
