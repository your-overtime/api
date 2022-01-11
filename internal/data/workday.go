package data

import (
	"time"

	"database/sql"
)

func (d *Db) GetWorkDay(day time.Time, userID uint) (*WorkDayDB, error) {
	w := WorkDayDB{}
	tx := d.Conn.Where("DATE(day) = DATE(?)", day).Where("user_id = ?", userID).First(&w)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return &w, nil
}

func (d *Db) SaveWorkDay(w *WorkDayDB) error {
	tx := d.Conn.Create(w)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (d *Db) DeleteWorkDay(day time.Time, userID uint) error {
	tx := d.Conn.Delete(&WorkDayDB{}, d.Conn.Where("DATE(day) = DATE(?)", day).Where("user_id = ?", userID))
	return tx.Error
}

func (d *Db) GetWorkDaysBetweenStartAndEnd(start time.Time, end time.Time, userID uint) ([]WorkDayDB, error) {
	ws := []WorkDayDB{}
	tx := d.Conn.Where("user_id = ?", userID).
		Where("DATE(day) BETWEEN DATE(?) AND DATE(?)", start, end).Find(&ws)
	if tx.Error != nil && tx.Error != sql.ErrNoRows {
		return nil, tx.Error
	}

	return ws, nil
}
