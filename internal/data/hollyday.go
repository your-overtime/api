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
		d.Conn.Where("? between start and end", start, end).Or(
			d.Conn.Where("end between ? and ?", start, end).Where("not start between ? and ?", start, end),
		)).Find(&hollydays)
	if tx.Error != nil && tx.Error != sql.ErrNoRows {
		return nil, tx.Error
	}

	return hollydays, nil
}
