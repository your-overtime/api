package data

import (
	"time"

	"git.goasum.de/jasper/overtime/pkg"
)

func (d *Db) GetWorkDay(day time.Time, userID uint) (*pkg.WorkDay, error) {
	w := pkg.WorkDay{}
	tx := d.Conn.Where("day = ?", day).Where("user_id = ?", userID).First(&w)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return &w, nil
}
