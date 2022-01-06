package data

import (
	"database/sql"
	"errors"
	"time"

	"github.com/your-overtime/api/pkg"
	"gorm.io/gorm"
)

func (d *Db) SaveActivity(a *ActivityDB) error {
	if len(a.Description) == 0 {
		return pkg.ErrEmptyDescriptionNotAllowed
	}

	tx := d.Conn.Save(a)
	return tx.Error
}

func (d *Db) GetRunningActivityByUserID(eID uint) (*ActivityDB, error) {
	a := ActivityDB{}
	tx := d.Conn.Where("user_id = ? and end is null", eID).First(&a)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, pkg.ErrNoActivityIsRunning
		}
		return nil, tx.Error
	}
	return &a, nil
}

func (d *Db) GetActivity(id uint, userID uint) (*ActivityDB, error) {
	a := ActivityDB{}
	tx := d.Conn.Where("id =? and user_id = ?", id, userID).First(&a)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &a, nil
}

func (d *Db) GetActivitiesBetweenStartAndEnd(start time.Time, end time.Time, userID uint) ([]ActivityDB, error) {
	activities := []ActivityDB{}
	tx := d.Conn.Where("user_id = ?", userID).
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

func (d *Db) MigrateActivityDuration() error {
	activities := []pkg.Activity{}
	if err := d.Conn.Find(&activities).Error; err != nil {
		return err
	}
	for _, a := range activities {
		if a.ActualDurationInMinutes == 0 && a.End != nil {
			a.ActualDurationInMinutes = uint(a.End.Sub(*a.Start).Minutes())
			if a.EventualDurationInMinutes == 0 {
				a.EventualDurationInMinutes = a.ActualDurationInMinutes
			}
			if err := d.Conn.Save(&a).Error; err != nil {
				return err
			}
		}
	}
	return nil
}
