package overtime

import (
	"fmt"
	"math"
	"time"

	"git.goasum.de/jasper/overtime/internal/data"
	"git.goasum.de/jasper/overtime/pkg"
)

type service struct {
	db *data.Db
}

func Init(db *data.Db) pkg.OvertimeService {
	return &service{
		db: db,
	}
}

func (s *service) SumActivityBetweenStartAndEndInMinutes(start time.Time, end time.Time, employeeID uint) (int64, error) {
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
		return 1
	case time.Wednesday:
		return 2
	case time.Thursday:
		return 3
	case time.Friday:
		return 4
	case time.Saturday:
		return 5
	case time.Sunday:
		return 6
	default:
		return 0
	}
}

func (s *service) SumHollydaysBetweenStartAndEndInMinutes(start time.Time, end time.Time, employee pkg.Employee) (int64, error) {
	hollydays, err := s.db.GetHollydaysBetweenStartAndEnd(start, end, employee.ID)
	if err != nil {
		return 0, err
	}
	freeTimeInMinutes := int64(0)
	for _, a := range hollydays {
		fmt.Println(a.Description, " ", a.Start, " ", a.End)
		s := a.Start
		c := 1
		for {
			if weekDayToInt(s.Weekday()) < 5 {
				freeTimeInMinutes += int64(employee.WeekWorkingTimeInMinutes) / 5
			}
			s = s.AddDate(0, 0, 1)
			c++
			if s.Unix() > a.End.Unix() {
				break
			}
		}
	}
	return freeTimeInMinutes, nil
}

func (s *service) calcOvertimeAndActivetime(start time.Time, now time.Time, e *pkg.Employee, wn int, wdNumber int) (int64, int64, error) {
	at, err := s.SumActivityBetweenStartAndEndInMinutes(start, now, e.ID)
	if err != nil {
		return 0, 0, err
	}

	ft, err := s.SumHollydaysBetweenStartAndEndInMinutes(start, now, *e)
	if err != nil {
		return 0, 0, err
	}

	diff := now.Sub(start)
	ds, _ := math.Modf(diff.Hours() / 24)

	ot := at + ft - int64(e.WeekWorkingTimeInMinutes/7)*int64(ds)
	if wdNumber < 5 {
		ot = at + ft - int64(e.WeekWorkingTimeInMinutes/7)*int64(ds) - int64(e.WeekWorkingTimeInMinutes/uint((5)))
	}

	return at, ot, nil
}

func (s *service) CalcOverview(e pkg.Employee) (*pkg.Overview, error) {
	now := time.Now()
	yyyy, mm, dd := now.Date()
	// TODO: sum working hours (maybe for the running year) and subtract e.WeekWorkingTime per week and hollydays
	wd := now.Weekday()
	_, wn := now.ISOWeek()
	wdNumber := weekDayToInt(wd)
	// This day
	dStart := time.Date(yyyy, mm, dd, 0, 0, 0, 0, now.Location())
	at, ot, err := s.calcOvertimeAndActivetime(dStart, now, &e, wn, wdNumber)
	if err != nil {
		return nil, err
	}
	// This week
	wStart := time.Date(yyyy, mm, dd-wdNumber, 0, 0, 0, 0, now.Location())
	wat, wot, err := s.calcOvertimeAndActivetime(wStart, now, &e, wn, wdNumber)
	if err != nil {
		return nil, err
	}
	// This month
	mStart := time.Date(yyyy, mm, 01, 0, 0, 0, 0, now.Location())
	mat, mot, err := s.calcOvertimeAndActivetime(mStart, now, &e, wn, wdNumber)
	if err != nil {
		return nil, err
	}
	// This year
	yStart := time.Date(yyyy, 01, 01, 0, 0, 0, 0, now.Location())
	yat, yot, err := s.calcOvertimeAndActivetime(yStart, now, &e, wn, wdNumber)
	if err != nil {
		return nil, err
	}
	o := &pkg.Overview{
		Date:                         now,
		WeekNumber:                   wn,
		Employee:                     &e,
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

func (s *service) StartActivity(desc string, employee pkg.Employee) (*pkg.Activity, error) {
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

func (s *service) AddActivity(a pkg.Activity, employee pkg.Employee) (*pkg.Activity, error) {
	err := s.db.SaveActivity(&a)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (s *service) StopRunningActivity(employee pkg.Employee) (*pkg.Activity, error) {
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

func (s *service) GetActivity(id uint, employee pkg.Employee) (*pkg.Activity, error) {
	a, err := s.db.GetActivity(id)
	if err != nil {
		return nil, err
	}

	if a.UserID != employee.ID {
		return nil, pkg.ErrPermissionDenied
	}

	return a, nil
}

func (s *service) GetActivities(start time.Time, end time.Time, employee pkg.Employee) ([]pkg.Activity, error) {
	a, err := s.db.GetActivitiesBetweenStartAndEnd(start, end, employee.ID)
	if err != nil {
		return nil, err
	}

	return a, nil
}
func (s *service) DelActivity(id uint, employee pkg.Employee) error {
	a, err := s.GetActivity(id, employee)
	if err != nil {
		return err
	}
	tx := s.db.Conn.Delete(a)
	return tx.Error
}

func (s *service) AddHollyday(h pkg.Hollyday, employee pkg.Employee) (*pkg.Hollyday, error) {
	err := s.db.SaveHollyday(&h)
	if err != nil {
		return nil, err
	}
	return &h, nil
}

func (s *service) GetHollyday(id uint, employee pkg.Employee) (*pkg.Hollyday, error) {
	h, err := s.db.GetHollyday(id)
	if err != nil {
		return nil, err
	}

	if h.UserID != employee.ID {
		return nil, pkg.ErrPermissionDenied
	}

	return h, nil
}

func (s *service) GetHollydays(start time.Time, end time.Time, employee pkg.Employee) ([]pkg.Hollyday, error) {
	h, err := s.db.GetHollydaysBetweenStartAndEnd(start, end, employee.ID)
	if err != nil {
		return nil, err
	}

	return h, nil
}

func (s *service) DelHollyday(id uint, employee pkg.Employee) error {
	h, err := s.GetHollyday(id, employee)
	if err != nil {
		return err
	}
	tx := s.db.Conn.Delete(h)
	return tx.Error
}
