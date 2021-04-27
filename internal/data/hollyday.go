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
	h := pkg.Hollyday{}
	tx := d.Conn.First(&h, id)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &h, nil
}

func (d *Db) GetHollydaysBetweenStartAndEnd(start time.Time, end time.Time, employeeID uint) ([]pkg.Hollyday, error) {
	hollydays := []pkg.Hollyday{}
	tx := d.Conn.Where("user_id = ?", employeeID).Where(
		d.Conn.Where("DATE(?) between DATE(start) and DATE(end)", start).Or(
			d.Conn.Where("DATE(end) between DATE(?) and DATE(?)", start, end).Where("DATE(start) not between DATE(?) and DATE(?)", start, end),
		)).Find(&hollydays)
	if tx.Error != nil && tx.Error != sql.ErrNoRows {
		return nil, tx.Error
	}

	return hollydays, nil
}
