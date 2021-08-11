package data

import (
	"crypto/sha256"
	"fmt"
	"strings"

	"github.com/your-overtime/api/pkg"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	log "github.com/sirupsen/logrus"
)

// TODO: remove in following releases
func (d *Db) MirgrateTokensToHashedTokens() error {
	tokens := []pkg.Token{}

	tx := d.Conn.Find(&tokens)
	if tx.Error != nil {
		return tx.Error
	}

	for _, t := range tokens {
		if len(t.Token) == 40 {
			t.Token = fmt.Sprintf("%x", sha256.Sum256([]byte(t.Token)))
			if err := d.Conn.Save(&t).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

func (d *Db) SaveEmployee(user *pkg.Employee) error {
	if !strings.HasPrefix(user.Password, "$2a$") || len(user.Password) < 60 {
		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Debug(err)
			return err
		}
		user.Password = string(hash)
	}
	var tx *gorm.DB
	if user.ID == 0 {
		tx = d.Conn.Create(&user)
	} else {
		tx = d.Conn.Save(&user)
	}

	if tx.Error != nil {
		log.Debug(tx.Error)
		return tx.Error
	}

	return nil
}

func (d *Db) GetEmployee(id uint) (*pkg.Employee, error) {
	e := pkg.Employee{}
	tx := d.Conn.First(&e, id)
	if tx.Error != nil {
		log.Debug(tx.Error)
		return nil, tx.Error
	}

	return &e, nil
}

func (d *Db) GetTokens(e pkg.Employee) ([]pkg.Token, error) {
	var ts []pkg.Token
	tx := d.Conn.Where("user_id = ?", e.ID).Find(&ts)
	if tx.Error != nil {
		log.Debug(tx.Error)
		return nil, tx.Error
	}

	return ts, nil
}

func (d *Db) GetTokenByToken(token string) (*pkg.Token, error) {
	var t pkg.Token
	tx := d.Conn.Where("token = ?", token).First(&t)
	if tx.Error != nil {
		log.Debug(tx.Error)
		return nil, tx.Error
	}

	return &t, nil
}

func (d *Db) SaveToken(token *pkg.Token) error {
	tx := d.Conn.Save(token)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (d *Db) GetEmployeeByToken(token string) (*pkg.Employee, error) {
	t := pkg.Token{}
	tx := d.Conn.Where("token = ?", token).First(&t)
	if tx.Error != nil {
		log.Debug(tx.Error)
		return nil, tx.Error
	}
	e := pkg.Employee{}
	tx = d.Conn.First(&e, t.UserID)
	if tx.Error != nil {
		log.Debug(tx.Error)
		return nil, tx.Error
	}
	return &e, nil
}

func (d *Db) GetEmployeeByLogin(login string) (*pkg.Employee, error) {
	e := &pkg.Employee{}
	tx := d.Conn.Where("login = ?", login).First(e)
	if tx.Error != nil {
		log.Debug(tx.Error)
		return nil, tx.Error
	}

	if e == nil {
		return nil, pkg.ErrUserNotFound
	}

	return e, nil
}
