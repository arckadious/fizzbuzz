package manager

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/arckadious/fizzbuzz/model"
	"github.com/arckadious/fizzbuzz/repository"
)

// Manager fizz
type Fizz struct {
	*Manager //Fizz class has attributes and methods from manager parent class
	repo     *repository.RepositoryFizz
}

// NewFizz Constructor
func NewFizz(m *Manager, repo *repository.RepositoryFizz) *Fizz {
	return &Fizz{
		Manager: m,
		repo:    repo,
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

	//TO DO
	res.SetData("NOT IMPLEMENTED YET").WriteJsonResponse(w)
}
