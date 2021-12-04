package data

import (
	"strings"

	"github.com/your-overtime/api/pkg"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	log "github.com/sirupsen/logrus"
)

func (d *Db) SaveUser(user *UserDB) error {
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

func (d *Db) GetUser(id uint) (*UserDB, error) {
	e := UserDB{}
	tx := d.Conn.First(&e, id)
	if tx.Error != nil {
		log.Debug(tx.Error)
		return nil, tx.Error
	}

	return &e, nil
}

func (d *Db) GetTokens(e UserDB) ([]TokenDB, error) {
	var ts []TokenDB
	tx := d.Conn.Where("user_id = ?", e.ID).Find(&ts)
	if tx.Error != nil {
		log.Debug(tx.Error)
		return nil, tx.Error
	}

	return ts, nil
}

func (d *Db) GetTokenByToken(token string) (*TokenDB, error) {
	var t TokenDB
	tx := d.Conn.Where("token = ?", token).First(&t)
	if tx.Error != nil {
		log.Debug(tx.Error)
		return nil, tx.Error
	}

	return &t, nil
}

func (d *Db) SaveToken(token *TokenDB) error {
	tx := d.Conn.Save(token)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (d *Db) GetUserByToken(token string) (*UserDB, error) {
	t := TokenDB{}
	tx := d.Conn.Where("token = ?", token).First(&t)
	if tx.Error != nil {
		log.Debug(tx.Error)
		return nil, tx.Error
	}
	e := UserDB{}
	tx = d.Conn.First(&e, t.UserID)
	if tx.Error != nil {
		log.Debug(tx.Error)
		return nil, tx.Error
	}
	return &e, nil
}

func (d *Db) GetUserByLogin(login string) (*UserDB, error) {
	e := &UserDB{}
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
