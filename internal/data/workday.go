package data

import (
	"time"

	"database/sql"
	"github.com/your-overtime/api/pkg"
)

func (d *Db) GetWorkDay(day time.Time, userID uint) (*pkg.WorkDay, error) {
	w := pkg.WorkDay{}
	tx := d.Conn.Where("DATE(day) = DATE(?)", day).Where("user_id = ?", userID).First(&w)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return &w, nil
}

func (d *Db) SaveWorkDay(w *pkg.WorkDay) error {
	tx := d.Conn.Create(w)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (d *Db) DeleteWorkDay(day time.Time, userID uint) error {
	tx := d.Conn.Delete(&pkg.WorkDay{}, d.Conn.Where("DATE(day) = DATE(?)", day).Where("user_id = ?", userID))
	return tx.Error
}

func (d *Db) GetWorkDaysBetweenStartAndEnd(start time.Time, end time.Time, employeeID uint) ([]pkg.WorkDay, error) {
	ws := []pkg.WorkDay{}
	tx := d.Conn.Where("user_id = ?", employeeID).
		Where("day >= ? and day <= ?", start, end).Find(&ws)
	if tx.Error != nil && tx.Error != sql.ErrNoRows {
		return nil, tx.Error
	}

	return ws, nil
}
