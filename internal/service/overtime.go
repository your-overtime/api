package service

import (
	"strings"
	"time"

	"git.goasum.de/jasper/overtime/internal/data"
	"git.goasum.de/jasper/overtime/pkg"
	log "github.com/sirupsen/logrus"
)

type Service struct {
	db *data.Db
}

func Init(db *data.Db) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) SumActivityBetweenStartAndEndInMinutes(start time.Time, end time.Time, employeeID uint) (int64, error) {
	activities, err := s.db.GetActivitiesBetweenStartAndEnd(start, end, employeeID)
	if err != nil {
		return 0, err
	}
	activityTimeInMinutes := int64(0)
	for _, a := range activities {
		var diff time.Duration
		if a.End == nil {
			diff = end.Sub(*a.Start)
		} else {
			diff = a.End.Sub(*a.Start)
		}

		activityTimeInMinutes += int64(diff.Minutes())
	}
	return activityTimeInMinutes, nil
}

func weekDayToInt(wd time.Weekday) int {
	switch wd {
	case time.Tuesday:
		return 2
	case time.Wednesday:
		return 3
	case time.Thursday:
		return 4
	case time.Friday:
		return 5
	case time.Saturday:
		return 6
	case time.Sunday:
		return 7
	default:
		return 1
	}
}

func (s *Service) SumHollydaysBetweenStartAndEndInMinutes(start time.Time, end time.Time, e pkg.Employee) (int64, error) {
	workingDays := strings.Split(e.WorkingDays, ",")
	hollydays, err := s.db.GetHollydaysBetweenStartAndEnd(start, end, e.ID)
	if err != nil {
		return 0, err
	}
	freeTimeInMinutes := int64(0)
	for _, a := range hollydays {
		st := a.Start
		if a.Start.Unix() < start.Unix() {
			st = start
		}
		for {
			if st.Unix() > end.Unix() || st.Unix() > a.End.Unix() {
				break
			}
			dayFreeTimeInMinutes := int64(0)
			if a.LegalHollyday {
				// Fix legal hollydays
				dayFreeTimeInMinutes = int64(e.WeekWorkingTimeInMinutes / 5)
			} else if weekDayToInt(st.Weekday()) < 6 && len(e.WorkingDays) == 0 {
				dayFreeTimeInMinutes = int64(e.WeekWorkingTimeInMinutes / 5)
			} else if strings.Contains(e.WorkingDays, st.Weekday().String()) {
				dayFreeTimeInMinutes = int64(e.WeekWorkingTimeInMinutes / uint(len(workingDays)))
			} else {
				st = st.AddDate(0, 0, 1)
				if st.Unix() > end.Unix() || st.Unix() > a.End.Unix() {
					break
				}
				continue
			}
			freeTimeInMinutes += dayFreeTimeInMinutes
			st = st.AddDate(0, 0, 1)
		}
	}
	return freeTimeInMinutes, nil
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
		en := time.Date(st.Year(), st.Month(), st.Day(), 23, 59, 59, 0, st.Location())
		if end.Unix() < en.Unix() {
			en = end
		}
		isNowDay := (be.Year() == now.Year() && be.Month() == now.Month() && be.Day() == now.Day())
		if !isNowDay {
			wd, err := s.db.GetWorkDay(be, e.ID)
			if err != nil {
				log.Debug(err)
			}
			if wd != nil && err == nil {
				activeTimeInMinutes += wd.ActiveTime
				overtimeInMinutes += wd.Overtime
				st = st.AddDate(0, 0, 1)
				continue
			}
		}

		dayWorkTimeInMinutes, err := s.CalcDailyWorktime(*e)
		if err != nil {
			return 0, 0, err
		}
		// if weekDayToInt(st.Weekday()) < 6 && len(e.WorkingDays) == 0 {
		//  	dayWorkTimeInMinutes = int64(e.WeekWorkingTimeInMinutes / 5)
		// } else if strings.Contains(e.WorkingDays, st.Weekday().String()) {
		//  	dayWorkTimeInMinutes = int64(e.WeekWorkingTimeInMinutes / uint(len(workingDays)))
		// }

		at, err := s.SumActivityBetweenStartAndEndInMinutes(be, en, e.ID)
		if err != nil {
			return 0, 0, err
		}

		ft, err := s.SumHollydaysBetweenStartAndEndInMinutes(be, en, *e)
		if err != nil {
			return 0, 0, err
		}
		dayOvertimeInMinutes := at + ft - int64(dayWorkTimeInMinutes)
		if !isNowDay {
			tx := s.db.Conn.Create(&pkg.WorkDay{
				Day:        be,
				Overtime:   dayOvertimeInMinutes,
				ActiveTime: at,
				UserID:     e.ID,
			})
			if tx.Error != nil {
				return 0, 0, tx.Error
			}
		}
		overtimeInMinutes += dayOvertimeInMinutes
		activeTimeInMinutes += at
		st = st.AddDate(0, 0, 1)
	}

	return activeTimeInMinutes, overtimeInMinutes, nil
}

func (s *Service) CalcOverview(e pkg.Employee) (*pkg.Overview, error) {
	now := time.Now()
	yyyy, mm, dd := now.Date()
	wd := now.Weekday()
	wdNumber := weekDayToInt(wd)
	// This day
	dStart := time.Date(yyyy, mm, dd, 0, 0, 0, 0, now.Location())
	at, ot, err := s.calcOvertimeAndActivetime(dStart, now, &e)
	if err != nil {
		return nil, err
	}
	// This week
	wStart := time.Date(yyyy, mm, dd-wdNumber+1, 0, 0, 0, 0, now.Location())
	wat, wot, err := s.calcOvertimeAndActivetime(wStart, now, &e)
	if err != nil {
		return nil, err
	}
	// This month
	mStart := time.Date(yyyy, mm, 01, 0, 0, 0, 0, now.Location())
	mat, mot, err := s.calcOvertimeAndActivetime(mStart, now, &e)
	if err != nil {
		return nil, err
	}
	// This year
	yStart := time.Date(yyyy, 01, 01, 0, 0, 0, 0, now.Location())
	yat, yot, err := s.calcOvertimeAndActivetime(yStart, now, &e)
	if err != nil {
		return nil, err
	}
	_, wn := now.ISOWeek()
	o := &pkg.Overview{
		Date:                         now,
		WeekNumber:                   wn,
		ActiveTimeThisDayInMinutes:   at,
		ActiveTimeThisWeekInMinutes:  wat,
		ActiveTimeThisMonthInMinutes: mat,
		ActiveTimeThisYearInMinutes:  yat,
		OvertimeThisDayInMinutes:     ot,
		OvertimeThisWeekInMinutes:    wot,
		OvertimeThisMonthInMinutes:   mot,
		OvertimeThisYearInMinutes:    yot,
	}
	cra, err := s.db.GetRunningActivityByEmployeeID(e.ID)
	if err == nil {
		o.ActiveActivity = cra
	}

	return o, nil
}

func (s *Service) CalcDailyWorktime(employee pkg.Employee) (uint, error) {
	now := time.Now()
	acs, err := s.GetActivities(
		time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()),
		now,
		employee,
	)
	if err != nil {
		return 0, err
	}

	weekStart := time.Date(now.Year(), now.Month(), now.Day()-int(now.Weekday()), 0, 0, 0, 0, now.Location())

	activeTimeInMinutes, err := s.SumActivityBetweenStartAndEndInMinutes(weekStart, now, employee.ID)
	if err != nil {
		return 0, err
	}

	wd := uint(weekDayToInt(now.Weekday()))
	dayWorkTimeInMinutes := uint(employee.WeekWorkingTimeInMinutes) / uint(employee.NumWorkingDays)

	if len(acs) == 0 && (7-wd) > employee.NumWorkingDays {
		return 0, nil
	}

	if activeTimeInMinutes >= int64(employee.WeekWorkingTimeInMinutes) {
		return 0, nil
	}

	return dayWorkTimeInMinutes, nil
}

func (s *Service) StartActivity(desc string, employee pkg.Employee) (*pkg.Activity, error) {
	ca, _ := s.db.GetRunningActivityByEmployeeID(employee.ID)
	if ca != nil {
		return nil, pkg.ErrActivityIsRunning
	}
	now := time.Now()
	a := pkg.Activity{
		UserID:      employee.ID,
		Start:       &now,
		Description: desc,
	}
	err := s.db.SaveActivity(&a)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (s *Service) AddActivity(a pkg.Activity, employee pkg.Employee) (*pkg.Activity, error) {
	err := s.db.SaveActivity(&a)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (s *Service) StopRunningActivity(employee pkg.Employee) (*pkg.Activity, error) {
	a, err := s.db.GetRunningActivityByEmployeeID(employee.ID)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	a.End = &now
	err = s.db.SaveActivity(a)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (s *Service) GetActivity(id uint, employee pkg.Employee) (*pkg.Activity, error) {
	a, err := s.db.GetActivity(id)
	if err != nil {
		return nil, err
	}

	if a != nil && a.UserID != employee.ID {
		return nil, pkg.ErrPermissionDenied
	}

	return a, nil
}

func (s *Service) GetActivities(start time.Time, end time.Time, employee pkg.Employee) ([]pkg.Activity, error) {
	a, err := s.db.GetActivitiesBetweenStartAndEnd(start, end, employee.ID)
	if err != nil {
		return nil, err
	}

	return a, nil
}
func (s *Service) DelActivity(id uint, employee pkg.Employee) error {
	a, err := s.GetActivity(id, employee)
	if err != nil {
		return err
	}
	tx := s.db.Conn.Delete(a)
	return tx.Error
}

func (s *Service) AddHollyday(h pkg.Hollyday, employee pkg.Employee) (*pkg.Hollyday, error) {
	err := s.db.SaveHollyday(&h)
	if err != nil {
		return nil, err
	}
	return &h, nil
}

func (s *Service) GetHollyday(id uint, employee pkg.Employee) (*pkg.Hollyday, error) {
	h, err := s.db.GetHollyday(id)
	if err != nil {
		return nil, err
	}

	if h.UserID != employee.ID {
		return nil, pkg.ErrPermissionDenied
	}

	return h, nil
}

func (s *Service) GetHollydays(start time.Time, end time.Time, employee pkg.Employee) ([]pkg.Hollyday, error) {
	h, err := s.db.GetHollydaysBetweenStartAndEnd(start, end, employee.ID)
	if err != nil {
		return nil, err
	}

	return h, nil
}

func (s *Service) DelHollyday(id uint, employee pkg.Employee) error {
	h, err := s.GetHollyday(id, employee)
	if err != nil {
		return err
	}
	tx := s.db.Conn.Delete(h)
	return tx.Error
}
