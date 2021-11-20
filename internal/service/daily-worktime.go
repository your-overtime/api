package service

import (
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/your-overtime/api/pkg"
)

func (s *Service) StaticCalculation(employee pkg.Employee, day time.Time) (uint, error) {
	if strings.Contains(employee.WorkingDays, day.Weekday().String()) {
		return uint(employee.WeekWorkingTimeInMinutes) / uint(len(employee.WorkingDaysAsArray())), nil
	}

	return 0, nil
}

func (s *Service) DynamicCalculation(employee pkg.Employee, day time.Time) (uint, error) {
	weekStart := time.Date(day.Year(), day.Month(), day.Day()-weekDayToInt(day.Weekday())+1, 0, 0, 0, 0, day.Location())
	dayWorkTimeInMinutes := uint(employee.WeekWorkingTimeInMinutes) / uint(employee.NumWorkingDays)

	wds, err := s.db.GetWorkDaysBetweenStartAndEnd(weekStart, day, employee.ID)
	if err != nil {
		log.Debugln(err)
		return 0, err
	}
	existingWDs := uint(0)
	for _, wd := range wds {
		if wd.ActiveTime > 0 || wd.IsHoliday || wd.Overtime > 0 {
			existingWDs += 1
		}
	}

	// Fix first week of the year
	if weekStart.Year() != day.Year() {
		existingWDs += 31 - uint(weekStart.Day())
	}

	dayActiveTimeInMinutes, err := s.SumActivityBetweenStartAndEndInMinutes(time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, day.Location()), day, employee.ID)
	if err != nil {
		log.Debugln(err)
		return 0, err
	}

	if existingWDs >= employee.NumWorkingDays ||
		(dayActiveTimeInMinutes == 0 && 7-weekDayToInt(day.Weekday())-int(employee.NumWorkingDays)+int(existingWDs) >= 0) {
		return 0, nil
	}

	return dayWorkTimeInMinutes, nil
}

func (s *Service) CalcDailyWorktime(employee pkg.Employee, day time.Time) (uint, error) {
	if len(employee.WorkingDaysAsArray()) > 0 {
		return s.StaticCalculation(employee, day)
	}

	return s.DynamicCalculation(employee, day)
}
