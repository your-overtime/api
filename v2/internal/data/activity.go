package data

import (
	"database/sql"
	"errors"
	"time"

	"github.com/your-overtime/api/v2/pkg"
	"github.com/your-overtime/api/v2/pkg/utils"
	"gorm.io/gorm"
)

func (d *Db) SaveActivity(a *ActivityDB) error {
	if len(a.Description) == 0 {
		return pkg.ErrEmptyDescriptionNotAllowed
	}
	if a.End != nil && a.End.Before(*a.Start) {
		return pkg.ErrStartMustBeBeforeEnd
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
	if err := HandleErr(tx.Error); err != nil {
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
			actualDuration := utils.DurationInMinutes(a.End.Sub(*a.Start))
			a.ActualDurationInMinutes = actualDuration
			a.EventualDurationInMinutes = a.ActualDurationInMinutes
			if err := d.Conn.Save(&a).Error; err != nil {
				return err
			}
		}
	}
	return nil
}
