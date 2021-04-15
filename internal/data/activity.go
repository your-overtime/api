package data

import (
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
	return &a, tx.Error
}

func (d *Db) GetActivity(id uint) (*pkg.Activity, error) {
	a := pkg.Activity{}
	tx := d.Conn.First(&a, id)
	return &a, tx.Error
}

func (d *Db) GetActivitiesBetweenStartAndEnd(start time.Time, end time.Time, employeeID uint) ([]pkg.Activity, error) {
	activities := []pkg.Activity{}
	tx := d.Conn.Where("user_id = ?", employeeID).Where("start between ? and ?", start, end).Find(&activities)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return activities, nil
}
