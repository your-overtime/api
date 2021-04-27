package data

import (
	"database/sql"
	"time"

	"git.goasum.de/jasper/overtime/pkg"
)

func (d *Db) SaveActivity(a *pkg.Activity) error {
	tx := d.Conn.Save(a)
	return tx.Error
}

func (d *Db) GetRunningActivityByEmployeeID(eID uint) (*pkg.Activity, error) {
	a := pkg.Activity{}
	tx := d.Conn.Where("user_id = ? and end is null", eID).First(&a)
	if tx.Error != nil {
		if tx.Error == sql.ErrNoRows {
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
	return nil, tx.Error
}

func (d *Db) GetActivitiesBetweenStartAndEnd(start time.Time, end time.Time, employeeID uint) ([]pkg.Activity, error) {
	activities := []pkg.Activity{}
	tx := d.Conn.Where("user_id = ?", employeeID).
		Where("start between ? and ?", start, end).
		Or(
			d.Conn.Where("end is not null AND end between ? and ?", start, end).
				Where("start not between ? and ?", start, end),
		).Find(&activities)
	if tx.Error != nil && tx.Error != sql.ErrNoRows {
		return nil, tx.Error
	}

	return activities, nil
}
