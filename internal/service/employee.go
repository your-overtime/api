package service

import (
	"fmt"

	utils "git.goasum.de/jasper/go-utils/pkg/string_utils"
	"git.goasum.de/jasper/overtime/pkg"
	log "github.com/sirupsen/logrus"
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
		log.Debug(err)
		return false
	}

	return true
}

func (s *Service) Login(login string, password string) (*pkg.Employee, error) {
	e, err := s.db.GetEmployeeByLogin(login)
	if err != nil {
		log.Debug(err)
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
		log.Debug(err)
		return nil, err
	}
	return &employee, nil
}

func (s *Service) UpdateAccount(fields map[string]interface{}, employee pkg.Employee) (*pkg.Employee, error) {
	for f := range fields {
		switch f {
		case "Name":
			employee.Name = fields[f].(string)
		case "Surname":
			employee.Surname = fields[f].(string)
		case "Password":
			employee.Password = fields[f].(string)
		case "Login":
			employee.Login = fields[f].(string)
		case "WorkingDays":
			employee.WorkingDays = fields[f].(string)
		case "WeekWorkingTimeInMinutes":
			employee.WeekWorkingTimeInMinutes = fields[f].(uint)
		default:
			return nil, pkg.ErrBadRequest
		}
	}
	dbE, err := s.SaveEmployee(employee)
	if err != nil {
		return nil, err
	}

	return dbE, nil
}

func (s *Service) DeleteEmployee(employeeID string) error {
	tx := s.db.Conn.Model(pkg.Employee{}).Delete(employeeID)
	return tx.Error
}

func (s *Service) GetTokens(employee pkg.Employee) ([]pkg.Token, error) {
	ts, err := s.db.GetTokens(employee)
	if err != nil {
		log.Debug(err)
		return nil, err
	}
	return ts, nil
}

func (s *Service) CreateToken(it pkg.InputToken, employee pkg.Employee) (*pkg.Token, error) {
	token := pkg.Token{
		UserID: employee.ID,
		Name:   it.Name,
		Token:  utils.RandString(40),
	}
	tx := s.db.Conn.Create(&token)
	if tx.Error != nil {
		log.Debug(tx.Error)
		return nil, tx.Error
	}
	return &token, nil
}

func (s *Service) DeleteToken(tokenID uint, employee pkg.Employee) error {
	var t pkg.Token
	tx := s.db.Conn.First(&t, tokenID)
	if tx.Error != nil {
		log.Debug(tx.Error)
		return tx.Error
	}
	if t.UserID == employee.ID {
		tx := s.db.Conn.Delete(&employee)
		log.Debug(tx.Error)
		return tx.Error
	}
	return pkg.ErrPermissionDenied
}
