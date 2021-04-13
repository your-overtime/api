package data

import "git.goasum.de/jasper/overtime/pkg"

func (d *Db) SaveEmployee(user *pkg.Employee) error {

	tx := d.Conn.Save(&user)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}
