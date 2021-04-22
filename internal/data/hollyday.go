package data

import (
	"database/sql"
	"time"

	"git.goasum.de/jasper/overtime/pkg"
)

func (d *Db) SaveHollyday(a *pkg.Hollyday) error {
	tx := d.Conn.Save(a)
	return tx.Error
}

func (d *Db) GetHollyday(id uint) (*pkg.Hollyday, error) {
	a := pkg.Hollyday{}
	tx := d.Conn.First(&a, id)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return nil, tx.Error
}

func (d *Db) GetHollydaysBetweenStartAndEnd(start time.Time, end time.Time, employeeID uint) ([]pkg.Hollyday, error) {
	hollydays := []pkg.Hollyday{}
	tx := d.Conn.Where("user_id = ?", employeeID).Where("start between ? and ?", start, end).Or("end between ? and ?", start, end).Find(&hollydays)
	if tx.Error != nil && tx.Error != sql.ErrNoRows {
		return nil, tx.Error
	}

	return hollydays, nil
}
