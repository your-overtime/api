package service

import (
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"github.com/your-overtime/api/internal/data"
	"github.com/your-overtime/api/pkg"
	"github.com/your-overtime/api/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

func createSHA256Hash(v string) string {
	return fmt.Sprintf("%x",
		sha256.Sum256([]byte(v)),
	)
}

func (s *Service) FromToken(token string) (*pkg.User, error) {
	hashedToken := createSHA256Hash(token)

	uDB, err := s.db.GetUserByToken(hashedToken)
	if err != nil {
		return nil, err
	}

	return &uDB.User, nil
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

func (s *Service) Login(login string, password string) (*pkg.User, error) {
	e, err := s.db.GetUserByLogin(login)
	if err != nil {
		log.Debug(err)
		return nil, err
	}

	if comparePasswords(e.Password, password) {
		return &e.User, nil
	}

	return nil, pkg.ErrInvalidCredentials
}

func (s *Service) SaveUser(user pkg.User, adminToken string) (*pkg.User, error) {
	var (
		u   *data.UserDB
		err error
	)
	if user.ID != 0 {
		u, err = s.db.GetUser(user.ID)
		u.User = user
	} else {
		u = &data.UserDB{
			User: user,
		}
	}

	err = s.db.SaveUser(u)
	if err != nil {
		log.Debug(err)
		return nil, err
	}
	return &u.User, nil
}

func (s *Service) UpdateAccount(fields map[string]interface{}, user pkg.User) (*pkg.User, error) {
	for f := range fields {
		switch f {
		case "Name":
			user.Name = fields[f].(string)
		case "Surname":
			user.Surname = fields[f].(string)
		case "Password":
			user.Password = fields[f].(string)
		case "Login":
			user.Login = fields[f].(string)
		case "WeekWorkingTimeInMinutes":
			user.WeekWorkingTimeInMinutes = utils.SafeGetUInt(fields[f])
		case "NumWorkingDays":
			user.NumWorkingDays = utils.SafeGetUInt(fields[f])
		case "NumHolidays":
			user.NumHolidays = utils.SafeGetUInt(fields[f])
		default:
			return nil, pkg.ErrBadRequest
		}
	}
	dbE, err := s.SaveUser(user, "")
	if err != nil {
		if err.(*mysql.MySQLError).Number == 1062 {
			return nil, pkg.ErrDuplicateValue
		}
		return nil, err
	}

	return dbE, nil
}

func (s *Service) DeleteUser(login string, adminToken string) error {
	tx := s.db.Conn.Model(data.UserDB{}).Delete("login = ?", login)
	return tx.Error
}

func (s *Service) GetTokens(user pkg.User) ([]pkg.Token, error) {
	uDB, err := s.db.GetUser(user.ID)
	if err != nil {
		log.Debug(err)
		return nil, err
	}
	tsDB, err := s.db.GetTokens(*uDB)
	if err != nil {
		log.Debug(err)
		return nil, err
	}
	ts := make([]pkg.Token, len(tsDB))
	for i := 0; i < len(tsDB); i++ {
		ts[i] = uDB.Tokens[i]
	}
	return ts, nil
}

func (s *Service) CreateToken(it pkg.InputToken, user pkg.User) (*pkg.Token, error) {
	// TODO add database method to create token?
	token := pkg.Token{
		UserID: user.ID,
		InputToken: pkg.InputToken{
			Name: it.Name,
		},
		Token: utils.RandString(40),
	}

	tx := s.db.Conn.Create(&token)
	if tx.Error != nil {
		log.Debug(tx.Error)
		return nil, tx.Error
	}

	respToken := token
	token.Token = createSHA256Hash(token.Token)
	err := s.db.SaveToken(&data.TokenDB{
		Token: token,
	})
	if err != nil {
		log.Debug(err)
		return nil, err
	}

	return &respToken, nil
}

func (s *Service) DeleteToken(tokenID uint, user pkg.User) error {
	var t pkg.Token
	tx := s.db.Conn.First(&t, tokenID)
	if tx.Error != nil {
		log.Debug(tx.Error)
		return tx.Error
	}
	if t.UserID == user.ID {
		tx := s.db.Conn.Delete(&user)
		log.Debug(tx.Error)
		return tx.Error
	}
	return pkg.ErrPermissionDenied
}

func (s *Service) GetAccount() (*pkg.User, error) {
	return nil, errors.New("not implemented")
}
