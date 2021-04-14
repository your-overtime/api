package data

import "git.goasum.de/jasper/overtime/pkg"

func (d *Db) SaveActivity(a *pkg.Activity) error {
	tx := d.Conn.Save(a)
	return tx.Error
}

func (d *Db) GetRunningActivityByEmployeeID(eID uint) (*pkg.Activity, error) {
	a := pkg.Activity{}
	tx := d.Conn.Where("user_id = ? and end is null", eID).First(&a)
	return &a, tx.Error
}

func (d *Db) GetActivity(id uint) (*pkg.Activity, error) {
	a := pkg.Activity{}
	tx := d.Conn.First(&a, id)
	return &a, tx.Error
}
