package service

import (
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"github.com/your-overtime/api/pkg"
	"github.com/your-overtime/api/pkg/stringutils"
	"golang.org/x/crypto/bcrypt"
)

func createSHA256Hash(v string) string {
	return fmt.Sprintf("%x",
		sha256.Sum256([]byte(v)),
	)
}

func (s *Service) FromToken(token string) (*pkg.Employee, error) {
	hashedToken := createSHA256Hash(token)

	t, err := s.db.GetTokenByToken(token)
	if err == nil {
		t.Token = hashedToken

		err = s.db.SaveToken(t)
		if err != nil {
			log.Debug(err)
			return nil, err
		}
	}

	return s.db.GetEmployeeByToken(hashedToken)
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

func (s *Service) SaveEmployee(employee pkg.Employee, adminToken string) (*pkg.Employee, error) {
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
		case "WeekWorkingTimeInMinutes":
			employee.WeekWorkingTimeInMinutes = uint(fields[f].(float64))
		case "NumWorkingDays":
			employee.NumWorkingDays = uint(fields[f].(float64))
		default:
			return nil, pkg.ErrBadRequest
		}
	}
	dbE, err := s.SaveEmployee(employee, "")
	if err != nil {
		if err.(*mysql.MySQLError).Number == 1062 {
			return nil, pkg.ErrDuplicateValue
		}
		return nil, err
	}

	return dbE, nil
}

func (s *Service) DeleteEmployee(login string, adminToken string) error {
	tx := s.db.Conn.Model(pkg.Employee{}).Delete("login = ?", login)
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
	// TODO add database method to create token?
	token := pkg.Token{
		UserID: employee.ID,
		Name:   it.Name,
		Token:  stringutils.RandString(40),
	}

	tx := s.db.Conn.Create(&token)
	if tx.Error != nil {
		log.Debug(tx.Error)
		return nil, tx.Error
	}

	respToken := token
	token.Token = createSHA256Hash(token.Token)
	err := s.db.SaveToken(&token)
	if err != nil {
		log.Debug(err)
		return nil, err
	}

	return &respToken, nil
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

func (s *Service) GetAccount() (*pkg.Employee, error) {
	return nil, errors.New("not implemented")
}
