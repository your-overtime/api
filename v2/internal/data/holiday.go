package data

import (
	"database/sql"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/your-overtime/api/v2/pkg"
)

func (d *Db) SaveHoliday(a *HolidayDB) error {
	tx := d.Conn.Save(a)
	return tx.Error
}

func (d *Db) GetHoliday(id uint) (*HolidayDB, error) {
	h := HolidayDB{}
	tx := d.Conn.First(&h, id)
	if err := HandleErr(tx.Error); err != nil {
		log.Debug(err)
		return nil, err
	}
	return &h, nil
}

func (d *Db) GetHolidaysBetweenStartAndEnd(start time.Time, end time.Time, userID uint) ([]HolidayDB, error) {
	holidays := []HolidayDB{}
	tx := d.Conn.Where("user_id = ?", userID).Where("? <= end AND ? >= start", start, end).Find(&holidays)
	if tx.Error != nil && tx.Error != sql.ErrNoRows {
		return nil, tx.Error
	}

	return holidays, nil
}

func (d *Db) GetHolidaysBetweenStartAndEndByType(start time.Time, end time.Time, hType pkg.HolidayType, userID uint) ([]HolidayDB, error) {
	holidays := []HolidayDB{}
	tx := d.Conn.Where("user_id = ?", userID).Where("? <= `end` AND ? >= `start` AND `type` = ?", start, end, hType).Find(&holidays)
	if tx.Error != nil && tx.Error != sql.ErrNoRows {
		return nil, tx.Error
	}

	return holidays, nil
}
