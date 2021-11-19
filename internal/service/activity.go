package service

import (
	"time"

	"github.com/your-overtime/api/pkg"

	log "github.com/sirupsen/logrus"
)

func (s *Service) SumActivityBetweenStartAndEndInMinutes(start time.Time, end time.Time, employeeID uint) (int64, error) {
	activities, err := s.db.GetActivitiesBetweenStartAndEnd(start, end, employeeID)
	if err != nil {
		log.Debug(err)
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

func (s *Service) StartActivity(desc string, employee pkg.Employee) (*pkg.Activity, error) {
	ca, _ := s.db.GetRunningActivityByEmployeeID(employee.ID)
	now := time.Now()
	if ca != nil {
		ca.End = &now
		if err := s.db.SaveActivity(ca); err != nil {
			return nil, err
		}
	}
	a := s.startActivityHook(pkg.Activity{
		UserID:      employee.ID,
		Start:       &now,
		Description: desc,
	})

	err := s.db.SaveActivity(&a)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (s *Service) AddActivity(a pkg.Activity, employee pkg.Employee) (*pkg.Activity, error) {
	// handle activities without end as new started activities
	if a.End == nil {
		return s.StartActivity(a.Description, employee)
	}

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

func (s *Service) UpdateActivity(a pkg.Activity, employee pkg.Employee) (*pkg.Activity, error) {
	err := s.db.SaveActivity(&a)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (s *Service) DelActivity(id uint, employee pkg.Employee) error {
	a, err := s.GetActivity(id, employee)
	if err != nil {
		return err
	}
	tx := s.db.Conn.Delete(a)
	return tx.Error
}
