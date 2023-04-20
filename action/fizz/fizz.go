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

func NewFizzAction(mng *manager.Fizz) *FizzAction {
	return &FizzAction{mng}
}

func (ac *FizzAction) HandleFizz(w http.ResponseWriter, r *http.Request) {

	res := ac.mng.GetApiResponse() //generic response, unchanged, send 200 OK with json body response
	var input model.Input

	//extract input
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		res.SetBadRequestResponse([]response.ApiError{{Code: response.ErrorInvalidData, Message: "Data JSON input bad format."}}).WriteJsonResponse(w)
		return
	}

	//Validate input
	if err := ac.mng.GetValidator().Struct(input); err != nil {
		res.SetBadRequestResponse([]response.ApiError{{Code: response.ErrorInvalidField, Message: err.Error()}}).WriteJsonResponse(w)
		return
	}

	ac.mng.HandleFizz(w, input)

}

func (ac *FizzAction) HandleStatistics(w http.ResponseWriter, r *http.Request) {
	ac.mng.HandleStatistics(w)
}
