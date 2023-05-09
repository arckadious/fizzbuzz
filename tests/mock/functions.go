package mock

import "github.com/stretchr/testify/mock"

// FizzRepositoryMock mock class from FizzRepository class
type FizzRepositoryMock struct {
	mock.Mock
}

// GetMostRequestUsed is a mock function of GetMostRequestUsed method from FizzRepository class
func (m *FizzRepositoryMock) GetMostRequestUsed() (msg string, hits int, noRows bool, err error) {
	args := m.Called()
	return args.String(0), args.Int(1), args.Bool(2), args.Error(3)
}
