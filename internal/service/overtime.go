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
	var activityTimeInMinutes int64
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

func (s *Service) SumHollydaysBetweenStartAndEndInMinutes(start time.Time, end time.Time, employee pkg.Employee) (int64, error) {
	hollydays, err := s.db.GetHollydaysBetweenStartAndEnd(start, end, employee.ID)
	if err != nil {
		return 0, err
	}
	freeTimeInMinutes := int64(0)
	for _, a := range hollydays {
		s := a.Start
		if a.Start.Unix() < start.Unix() {
			s = start
		}
		for {
			if weekDayToInt(s.Weekday()) < 6 {
				freeTimeInMinutes += int64(employee.WeekWorkingTimeInMinutes / 5)
			}
			s = s.AddDate(0, 0, 1)
			if s.Unix() > a.End.Unix() {
				break
			}
		}
	}
	return freeTimeInMinutes, nil
}

func (s *Service) calcOvertimeAndActivetime(start time.Time, end time.Time, e *pkg.Employee) (int64, int64, error) {
	workingDays := strings.Split(e.WorkingDays, ",")

	var (
		overtimeInMinutes   int64
		activeTimeInMinutes int64
	)
	st := start
	for {
		be := time.Date(st.Year(), st.Month(), st.Day(), 0, 0, 0, 0, st.Location())
		en := time.Date(st.Year(), st.Month(), st.Day(), 23, 59, 59, 0, st.Location())
		wd, err := s.db.GetWorkDay(be, e.ID)
		if err != nil {
			log.Debug(err)
		}
		if wd != nil {
			activeTimeInMinutes = wd.ActiveTime
			overtimeInMinutes = wd.Overtime
			continue
		}

		var dayWorkTimeInMinutes int64
		if st.Weekday() < 6 && len(e.WorkingDays) == 0 {
			dayWorkTimeInMinutes = int64(e.WeekWorkingTimeInMinutes / 5)
		} else if strings.Contains(e.WorkingDays, st.Weekday().String()) {
			dayWorkTimeInMinutes = int64(e.WeekWorkingTimeInMinutes / uint(len(workingDays)))
		} else {
			continue
		}

		at, err := s.SumActivityBetweenStartAndEndInMinutes(be, en, e.ID)
		if err != nil {
			return 0, 0, err
		}

		ft, err := s.SumHollydaysBetweenStartAndEndInMinutes(be, en, *e)
		if err != nil {
			return 0, 0, err
		}
		dayOvertimeInMinutes := at + ft - dayWorkTimeInMinutes
		tx := s.db.Conn.Save(pkg.WorkDay{
			Day:        be,
			Overtime:   dayOvertimeInMinutes,
			ActiveTime: at,
		})
		if tx.Error != nil {
			return 0, 0, tx.Error
		}
		st := st.AddDate(0, 0, 1)
		if st.Unix() > end.Unix() {
			break
		}
		overtimeInMinutes += dayOvertimeInMinutes
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

	if a.UserID != employee.ID {
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
