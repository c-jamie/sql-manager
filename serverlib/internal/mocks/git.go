package mocks

import "fmt"

type MockRepo struct {
}

func (mi *MockRepo) GetFile(file string) (string, error) {
	return fmt.Sprintf("select * from %s", file), nil
}
