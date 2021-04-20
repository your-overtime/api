package data

import (
	"strings"

	"git.goasum.de/jasper/overtime/pkg"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func (d *Db) SaveEmployee(user *pkg.Employee) error {
	if !strings.HasPrefix(user.Password, "$2a$") || len(user.Password) < 60 {
		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
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
		return tx.Error
	}

	return nil
}

func (d *Db) GetEmployee(id int) (*pkg.Employee, error) {
	e := pkg.Employee{}
	tx := d.Conn.First(&e, id)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &e, nil
}

func (d *Db) GetTokens(e pkg.Employee) ([]pkg.Token, error) {
	var ts []pkg.Token
	tx := d.Conn.Where("user_id = ?", e.ID).Find(&ts)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return ts, nil
}

func (d *Db) GetEmployeeByToken(token string) (*pkg.Employee, error) {
	t := pkg.Token{}
	tx := d.Conn.Where("token = ?", token).First(&t)
	if tx.Error != nil {
		return nil, tx.Error
	}
	e := pkg.Employee{}
	tx = d.Conn.First(&e, t.ID)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &e, nil
}

func (d *Db) GetEmployeeByLogin(login string) (*pkg.Employee, error) {
	e := pkg.Employee{}
	tx := d.Conn.Where("login = ?", login).First(&e)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &e, nil
}
