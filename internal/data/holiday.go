package data

import (
	"database/sql"
	"time"

	"github.com/your-overtime/api/pkg"
)

func (d *Db) SaveHoliday(a *pkg.Holiday) error {
	tx := d.Conn.Save(a)
	return tx.Error
}

func (d *Db) CountHolidaysBetweenStartAndEnd(start time.Time, end time.Time, employeeID uint) (uint, error) {
	nh := int64(0)
	tx := d.Conn.Model(pkg.Holiday{}).
		Where("user_id = ?", employeeID).
		Where("type = ?", pkg.HolidayTypeFree).Count(&nh)

	if tx.Error != nil {
		return 0, tx.Error
	}

	return uint(nh), nil
}

func (d *Db) GetHoliday(id uint) (*pkg.Holiday, error) {
	h := pkg.Holiday{}
	tx := d.Conn.First(&h, id)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &h, nil
}

func (d *Db) GetHolidaysBetweenStartAndEnd(start time.Time, end time.Time, employeeID uint) ([]pkg.Holiday, error) {
	holidays := []pkg.Holiday{}
	tx := d.Conn.Where("user_id = ?", employeeID).Where("? <= end AND ? >= start", start, end).Find(&holidays)
	if tx.Error != nil && tx.Error != sql.ErrNoRows {
		return nil, tx.Error
	}

	return holidays, nil
}

func (d *Db) GetHolidaysBetweenStartAndEndByType(start time.Time, end time.Time, hType pkg.HolidayType, employeeID uint) ([]pkg.Holiday, error) {
	holidays := []pkg.Holiday{}
	tx := d.Conn.Where("user_id = ?", employeeID).Where("? <= `end` AND ? >= `start` AND `type` = ?", start, end, hType).Find(&holidays)
	if tx.Error != nil && tx.Error != sql.ErrNoRows {
		return nil, tx.Error
	}

	return holidays, nil
}
