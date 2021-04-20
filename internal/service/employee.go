package service

import (
	"fmt"

	"git.goasum.de/jasper/overtime/pkg"
	"golang.org/x/crypto/bcrypt"
)

func (s *Service) FromToken(token string) (*pkg.Employee, error) {
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

func (s *Service) Login(login string, password string) (*pkg.Employee, error) {
	e, err := s.db.GetEmployeeByLogin(login)
	if err != nil {
		return nil, err
	}

	if comparePasswords(e.Password, password) {
		return e, nil
	}

	return nil, pkg.ErrInvalidCredentials
}

func (s *Service) SaveEmployee(employee pkg.Employee) (*pkg.Employee, error) {
	err := s.db.SaveEmployee(&employee)
	if err != nil {
		return nil, err
	}
	return &employee, nil
}

func (s *Service) DeleteEmployee(employeeID string) error {
	tx := s.db.Conn.Model(pkg.Employee{}).Delete(employeeID)
	return tx.Error
}

func (s *Service) GetTokens(employee pkg.Employee) ([]pkg.Token, error) {
	ts, err := s.db.GetTokens(employee)
	if err != nil {
		return nil, err
	}
	return ts, nil
}

func (s *Service) SaveToken(token pkg.Token, employee pkg.Employee) (*pkg.Token, error) {
	token.UserID = employee.ID
	tx := s.db.Conn.Save(&token)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &token, nil
}

func (s *Service) DeleteToken(tokenID uint, employee pkg.Employee) error {
	var t pkg.Token
	tx := s.db.Conn.First(&t, tokenID)
	if tx.Error != nil {
		return tx.Error
	}
	if t.UserID == employee.ID {
		tx := s.db.Conn.Delete(&employee)
		return tx.Error
	}
	return pkg.ErrPermissionDenied
}
