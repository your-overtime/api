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

		hs, mf := math.Modf(diff.Hours())

		activityTimeInMinutes += int64(hs*60 + mf*60)
		fmt.Println(a.ID, " ", a.Start, " ", a.End, " ", mf*60)
	}
	return activityTimeInMinutes, nil
}

func (s *service) SumHollydaysBetweenStartAndEndInMinutes(start time.Time, end time.Time, employee pkg.Employee) (int64, error) {
	hollydays, err := s.db.GetHollydayBetweenStartAndEnd(start, end, employee.ID)
	if err != nil {
		return 0, err
	}
	var freeTimeInMinutes int64
	for _, a := range hollydays {
		var (
			hs   float64
			mf   float64
			mins int64
		)
		if a.Start.Unix() >= start.Unix() {
			diff := a.End.Sub(a.Start)
			hs, mf = math.Modf(diff.Hours())
			mins = int64(hs*60 + mf*60)
		} else {
			diff := a.End.Sub(start)
			hs, mf = math.Modf(diff.Hours())
			mins = int64(hs*60 + mf*60)
			mins = int64(mins - int64(employee.WeekWorkingTimeInMinutes)/5)
		}

		freeTimeInMinutes += mins
		fmt.Println("Holly ", a.ID, " ", hs*60+mf*60, " ", mins, " ", int64(employee.WeekWorkingTimeInMinutes)/5)
	}
	return freeTimeInMinutes, nil
}

func (s *service) CalcOverviewForThisYear(e pkg.Employee) (*pkg.Overview, error) {
	now := time.Now()
	yyyy, _, _ := now.Date()
	wd := now.Weekday()
	_, wn := now.ISOWeek()
	var wdNumber int
	switch wd {
	case time.Tuesday:
		wdNumber = 1
	case time.Wednesday:
		wdNumber = 2
	case time.Thursday:
		wdNumber = 3
	case time.Friday:
		wdNumber = 4
	case time.Saturday:
		wdNumber = 5
	case time.Sunday:
		wdNumber = 6
	default:
		wdNumber = 0
	}
	wStart := time.Date(yyyy, 01, 01, 0, 0, 0, 0, now.Location())
	at, err := s.SumActivityBetweenStartAndEndInMinutes(wStart, now, e.ID)
	if err != nil {
		return nil, err
	}

	ft, err := s.SumHollydaysBetweenStartAndEndInMinutes(wStart, now, e)
	if err != nil {
		return nil, err
	}

	ot := at + ft - int64(e.WeekWorkingTimeInMinutes*uint(wn-1)) - int64(e.WeekWorkingTimeInMinutes/uint((7-wdNumber)))

	o := &pkg.Overview{
		Date:               time.Now(),
		WeekNumber:         wn,
		Employee:           &e,
		ActiveTimeThisWeek: at,
		OvertimeInMinutes:  ot,
	}
	cra, err := s.db.GetRunningActivityByEmployeeID(e.ID)
	if err == nil {
		o.ActiveActivity = cra
	}

	return o, nil
}

func (s *service) CalcCurrentOverview(e pkg.Employee) (*pkg.Overview, error) {
	now := time.Now()
	yyyy, mm, dd := now.Date()
	// TODO: sum working hours (maybe for the running year) and subtract e.WeekWorkingTime per week and hollydays
	wd := now.Weekday()
	_, wn := now.ISOWeek()
	var wdNumber int
	switch wd {
	case time.Tuesday:
		wdNumber = 1
	case time.Wednesday:
		wdNumber = 2
	case time.Thursday:
		wdNumber = 3
	case time.Friday:
		wdNumber = 4
	case time.Saturday:
		wdNumber = 5
	case time.Sunday:
		wdNumber = 6
	default:
		wdNumber = 0
	}
	wStart := time.Date(yyyy, mm, dd-wdNumber, 0, 0, 0, 0, now.Location())
	at, err := s.SumActivityBetweenStartAndEndInMinutes(wStart, now, e.ID)
	if err != nil {
		return nil, err
	}
	ft, err := s.SumHollydaysBetweenStartAndEndInMinutes(wStart, now, e)
	if err != nil {
		return nil, err
	}

	ot := at + ft - int64(e.WeekWorkingTimeInMinutes/uint((5-wdNumber)))

	o := &pkg.Overview{
		Date:               time.Now(),
		WeekNumber:         wn,
		Employee:           &e,
		ActiveTimeThisWeek: at,
		OvertimeInMinutes:  ot,
	}
	cra, err := s.db.GetRunningActivityByEmployeeID(e.ID)
	if err == nil {
		o.ActiveActivity = cra
	}

	return o, nil
}

func (s *service) StartActivity(desc string, employee pkg.Employee) (*pkg.Activity, error) {
	now := time.Now()
	a := pkg.Activity{
		UserID:      employee.ID,
		Start:       &now,
		Description: desc,
	}
	err := s.db.SaveActivity(&a)
	return &a, err
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

func (s *service) DelHollyday(id uint, employee pkg.Employee) error {
	h, err := s.GetHollyday(id, employee)
	if err != nil {
		return err
	}
	tx := s.db.Conn.Delete(h)
	return tx.Error
}
