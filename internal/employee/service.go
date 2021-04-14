package employee

import (
	"fmt"

	"git.goasum.de/jasper/overtime/internal/data"
	"git.goasum.de/jasper/overtime/pkg"
	"golang.org/x/crypto/bcrypt"
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
	fmt.Println(token)
	return s.db.GetEmployeeByToken(token)
}

func comparePasswords(hashedPw string, plainPw string) bool {
	bytePlain := []byte(plainPw)
	byteHash := []byte(hashedPw)
	err := bcrypt.CompareHashAndPassword(byteHash, bytePlain)
	if err != nil {
		return false
	}

	return true
}

func (s *service) Login(login string, password string) (*pkg.Employee, error) {
	e, err := s.db.GetEmployeeByLogin(login)
	if err != nil {
		return nil, err
	}

	if comparePasswords(e.Password, password) {
		return e, nil
	}

	return nil, pkg.ErrInvalidCredentials
}

func (s *service) AddEmployee(employee pkg.Employee) (*pkg.Employee, error) {
	err := s.db.SaveEmployee(&employee)
	if err != nil {
		return nil, err
	}
	return &employee, nil
}
