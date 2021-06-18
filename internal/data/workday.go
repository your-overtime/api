package data

import (
	"time"

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

func (d *Db) DeleteWorkDay(day time.Time, userID uint) error {
	tx := d.Conn.Delete(&pkg.WorkDay{}, d.Conn.Where("DATE(day) = DATE(?)", day).Where("user_id = ?", userID))
	return tx.Error
}
