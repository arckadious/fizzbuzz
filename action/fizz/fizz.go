// This package extract and validate data
package fizz

import (
	"encoding/json"
	"net/http"

	cst "github.com/arckadious/fizzbuzz/constant"

	"github.com/arckadious/fizzbuzz/manager"
	"github.com/arckadious/fizzbuzz/model"
	"github.com/arckadious/fizzbuzz/response"
)

// FizzAction class
type FizzAction struct {
	mng *manager.Fizz
}

// New constructor FizzAction
func New(mng *manager.Fizz) *FizzAction {
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
		res.SetErrorResponse(http.StatusBadRequest, []response.ApiError{{Code: cst.ErrorInvalidData, Message: "Data JSON input bad format."}}).WriteJSONResponse(w)
		return
	}

	//Validate input
	if err := ac.mng.GetValidator().Struct(input); err != nil {
		res.SetErrorResponse(http.StatusBadRequest, []response.ApiError{{Code: cst.ErrorInvalidField, Message: err.Error()}}).WriteJSONResponse(w)
		return
	}

	ac.mng.HandleFizz(w, input)

}

// HandleStatistics send back statistics about what the most used request has been
func (ac *FizzAction) HandleStatistics(w http.ResponseWriter, r *http.Request) {
	ac.mng.HandleStatistics(w)
}
