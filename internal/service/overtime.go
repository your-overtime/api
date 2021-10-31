package service

import (
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/your-overtime/api/internal/data"
	"github.com/your-overtime/api/pkg"
)

type Service struct {
	db *data.Db
}

func Init(db *data.Db) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) calcOvertimeAndActivetime(start time.Time, end time.Time, e *pkg.Employee) (int64, int64, error) {
	overtimeInMinutes := int64(0)
	activeTimeInMinutes := int64(0)
	now := time.Now()

	st := start
	for {
		if st.Unix() > end.Unix() {
			break
		}
		be := time.Date(st.Year(), st.Month(), st.Day(), 0, 0, 0, 0, st.Location())
		en := time.Date(st.Year(), st.Month(), st.Day(), 23, 59, 59, 59, st.Location())
		if end.Unix() < en.Unix() {
			en = end
		}
		isNowDay := (be.Year() == now.Year() && be.Month() == now.Month() && be.Day() == now.Day())
		if !isNowDay {
			wd, err := s.db.GetWorkDay(be, e.ID)
			if err != nil {
				log.Info(err)
			}
			if wd != nil && err == nil {
				activeTimeInMinutes += wd.ActiveTime
				overtimeInMinutes += wd.Overtime
				st = st.AddDate(0, 0, 1)
				continue
			}
		}

		dayWorkTimeInMinutes, err := s.CalcDailyWorktime(*e, en)
		if err != nil {
			log.Debug(err)
			return 0, 0, err
		}

		at, err := s.SumActivityBetweenStartAndEndInMinutes(be, en, e.ID)
		if err != nil {
			log.Debug(err)
			return 0, 0, err
		}

		ft, isLegal, err := s.SumHolidaysBetweenStartAndEndInMinutes(be, en, *e)
		if err != nil {
			log.Debug(err)
			return 0, 0, err
		}
		dayOvertimeInMinutes := at + ft - int64(dayWorkTimeInMinutes)
		if ft > 0 && dayWorkTimeInMinutes == 0 && !isLegal {
			dayOvertimeInMinutes = at
		}
		if !isNowDay {
			err = s.db.SaveWorkDay(&pkg.WorkDay{
				Day:        be,
				Overtime:   dayOvertimeInMinutes,
				ActiveTime: at,
				UserID:     e.ID,
				IsHoliday:  ft > 0,
			})
			if err != nil {
				return 0, 0, err
			}
		}
		overtimeInMinutes += dayOvertimeInMinutes
		activeTimeInMinutes += at
		st = st.AddDate(0, 0, 1)
	}

	return activeTimeInMinutes, overtimeInMinutes, nil
}

func (s *Service) CalcOverview(e pkg.Employee, day time.Time) (*pkg.Overview, error) {
	yyyy, mm, dd := day.Date()
	wd := day.Weekday()
	wdNumber := weekDayToInt(wd)
	// This year
	yStart := time.Date(yyyy, 01, 01, 0, 0, 0, 0, day.Location())
	yat, yot, err := s.calcOvertimeAndActivetime(yStart, day, &e)
	if err != nil {
		return nil, err
	}

	holidays, err := s.CountUsedHolidaysBetweenStartAndEnd(yStart, day, e)
	if err != nil {
		return nil, err
	}

	// This month
	mStart := time.Date(yyyy, mm, 01, 0, 0, 0, 0, day.Location())
	mat, mot, err := s.calcOvertimeAndActivetime(mStart, day, &e)
	if err != nil {
		return nil, err
	}
	// This week
	wStart := time.Date(yyyy, mm, dd-wdNumber+1, 0, 0, 0, 0, day.Location())
	wat, wot, err := s.calcOvertimeAndActivetime(wStart, day, &e)
	if err != nil {
		return nil, err
	}
	// This day
	dStart := time.Date(yyyy, mm, dd, 0, 0, 0, 0, day.Location())
	at, ot, err := s.calcOvertimeAndActivetime(dStart, day, &e)
	if err != nil {
		return nil, err
	}

	_, wn := day.ISOWeek()
	o := &pkg.Overview{
		Date:                         day,
		WeekNumber:                   wn,
		ActiveTimeThisDayInMinutes:   at,
		ActiveTimeThisWeekInMinutes:  wat,
		ActiveTimeThisMonthInMinutes: mat,
		ActiveTimeThisYearInMinutes:  yat,
		OvertimeThisDayInMinutes:     ot,
		OvertimeThisWeekInMinutes:    wot,
		OvertimeThisMonthInMinutes:   mot,
		OvertimeThisYearInMinutes:    yot,
		UsedHolidays:                 int(holidays),
		HolidaysStillAvailable:       int(e.NumHolidays - holidays),
	}
	cra, err := s.db.GetRunningActivityByEmployeeID(e.ID)
	if err == nil {
		o.ActiveActivity = cra
	}

	return o, nil
}

func (s *Service) CalcDailyWorktime(employee pkg.Employee, day time.Time) (uint, error) {
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
