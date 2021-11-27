package data

import (
	"database/sql"
	"errors"
	"time"

	"github.com/your-overtime/api/pkg"
	"gorm.io/gorm"
)

func (d *Db) SaveActivity(a *pkg.Activity) error {
	if len(a.Description) == 0 {
		return pkg.ErrEmptyDescriptionNotAllowed
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
		Where(
			d.Conn.Where("end IS NULL AND start >= ? AND start <= ?", start, end).Or(
				d.Conn.Where("end IS NOT NULL").
					Where("start >= ? AND start <= ?", start, end).
					Where("end >= ? AND end <= ?", start, end),
			)).Find(&activities)

	if tx.Error != nil && tx.Error != sql.ErrNoRows {
		return nil, tx.Error
	}

	return activities, nil
}
