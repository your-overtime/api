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
	tx := d.Conn.Where("user_id = ?", employeeID).Where("? < end AND ? > start", start, end).Find(&holidays)
	if tx.Error != nil && tx.Error != sql.ErrNoRows {
		return nil, tx.Error
	}

	return holidays, nil
}
