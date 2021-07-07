package data

import (
	"database/sql"
	"time"

	"github.com/your-overtime/api/pkg"
)

func (d *Db) SaveHoliday(a *pkg.Holiday) error {
	tx := d.Conn.Save(a)
	return tx.Error
}

func (d *Db) GetHoliday(id uint) (*pkg.Holiday, error) {
	h := pkg.Holiday{}
	tx := d.Conn.First(&h, id)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &h, nil
}

func (d *Db) GetHolidaysBetweenStartAndEnd(start time.Time, end time.Time, employeeID uint) ([]pkg.Holiday, error) {
	holidays := []pkg.Holiday{}
	tx := d.Conn.Where("user_id = ?", employeeID).
		Where(
			d.Conn.Where("start >= ? AND start <= ?", start, end).
				Or("end >= ? AND end <= ?", start, end).Or(
				d.Conn.Where("? <= end AND ? <= end", start, end).Where("? >= start AND ? >= start", start, end)),
		).Find(&holidays)
	if tx.Error != nil && tx.Error != sql.ErrNoRows {
		return nil, tx.Error
	}

	return holidays, nil
}
