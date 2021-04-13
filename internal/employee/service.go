package employee

import (
	"errors"

	"git.goasum.de/jasper/overtime/internal/data"
	"git.goasum.de/jasper/overtime/pkg"
)

type service struct {
	db *data.Db
}

func Init(db *data.Db) pkg.EmployeeService {
	return &service{
		db: db,
	}
}

func (s *service) FromToken(token string) (*pkg.Employee, error) {
	return nil, errors.New("not implemented yet")
}

func (s *service) Login(login string, password string) (*pkg.Employee, error) {
	return nil, errors.New("not implemented yet")
}

func (s *service) AddEmployee(employee pkg.Employee) (*pkg.Employee, error) {
	return nil, errors.New("not implemented yet")
}
