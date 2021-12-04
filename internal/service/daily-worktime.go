package service

import (
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/your-overtime/api/pkg"
	"github.com/your-overtime/api/pkg/utils"
)

// StaticCalculation returns the daily working time if the transfer day is in the WeekWorkingdays otherwise 0
func (s *Service) StaticCalculation(user pkg.User, day time.Time) (uint, error) {
	if strings.Contains(user.WorkingDays, day.Weekday().String()) {
		return uint(user.WeekWorkingTimeInMinutes) / uint(len(user.WorkingDaysAsArray())), nil
	}

	return 0, nil
}

// DynamicCalculation returns the daily working time if the number of working days in the current
// week is < then NumWorkingDays and an activity exists for the day passing.
// The method returns the daily working time if there are no activities when the number of days in the week is smaller
// than the number of working days.
func (s *Service) DynamicCalculation(user pkg.User, day time.Time) (uint, error) {
	weekStart := time.Date(day.Year(), day.Month(), day.Day()-weekDayToInt(day.Weekday())+1, 0, 0, 0, 0, day.Location())
	dayWorkTimeInMinutes := uint(user.WeekWorkingTimeInMinutes) / uint(user.NumWorkingDays)

	wds, err := s.db.GetWorkDaysBetweenStartAndEnd(weekStart, day, user.ID)
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
	if weekStart.Year() < day.Year() {
		existingWDs += 31 - uint(weekStart.Day())
	}

	dayActiveTimeInMinutes, err := s.SumActivityBetweenStartAndEndInMinutes(utils.DayStart(day), day, user.ID)
	if err != nil {
		log.Debugln(err)
		return 0, err
	}

	if existingWDs >= user.NumWorkingDays ||
		(dayActiveTimeInMinutes == 0 && 7-weekDayToInt(day.Weekday())-int(user.NumWorkingDays)+int(existingWDs) >= 0) {
		return 0, nil
	}

	return dayWorkTimeInMinutes, nil
}

// CalcDailyWorktime returns the daily working time and selects the calculation method used, depending on whether
// fixed working days are stored or not. If not the dynamic calculation method is used
func (s *Service) CalcDailyWorktime(user pkg.User, day time.Time) (uint, error) {
	if len(user.WorkingDaysAsArray()) > 0 {
		return s.StaticCalculation(user, day)
	}

	return s.DynamicCalculation(user, day)
}
