package data

import (
	"database/sql"
	"errors"
	"time"

	"github.com/your-overtime/api/pkg"
	"gorm.io/gorm"
)

func (d *Db) SaveActivity(a *pkg.Activity) error {
	now := time.Now()
	// TODO don't handle domain logic in data layer
	if !(a.Start.Year() == now.Year() && a.Start.Month() == now.Month() && a.Start.Day() == now.Day()) {
		err := d.DeleteWorkDay(time.Date(a.Start.Day(), a.Start.Month(), a.Start.Day(), 0, 0, 0, 0, a.Start.Location()), a.UserID)
		if err != nil {
			return err
		}
	}
	tx := d.Conn.Save(a)
	return tx.Error
}

func (d *Db) GetRunningActivityByEmployeeID(eID uint) (*pkg.Activity, error) {
	a := pkg.Activity{}
	tx := d.Conn.Where("user_id = ? and end is null", eID).First(&a)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, pkg.ErrNoActivityIsRunning
		}
		return nil, tx.Error
	}
	return &a, nil
}

func (d *Db) GetActivity(id uint) (*pkg.Activity, error) {
	a := pkg.Activity{}
	tx := d.Conn.First(&a, id)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &a, nil
}

func (d *Db) GetActivitiesBetweenStartAndEnd(start time.Time, end time.Time, employeeID uint) ([]pkg.Activity, error) {
	activities := []pkg.Activity{}
	tx := d.Conn.Where("user_id = ?", employeeID).
		Where("end IS NULL AND start >= ? AND start <= ?", start, end).Or(
		d.Conn.Where("end IS NOT NULL").
			Where("start >= ? AND start <= ?", start, end).
			Where("end >= ? AND end <= ?", start, end),
	).Find(&activities)

	if tx.Error != nil && tx.Error != sql.ErrNoRows {
		return nil, tx.Error
	}

	return activities, nil
}
