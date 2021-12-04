package service

import (
	"time"

	"github.com/your-overtime/api/internal/data"
	"github.com/your-overtime/api/pkg"
	"github.com/your-overtime/api/pkg/utils"

	log "github.com/sirupsen/logrus"
)

func (s *Service) SumActivityBetweenStartAndEndInMinutes(start time.Time, end time.Time, userID uint) (int64, error) {
	activities, err := s.db.GetActivitiesBetweenStartAndEnd(start, end, userID)
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

func (s *Service) StartActivity(desc string, user pkg.User) (*pkg.Activity, error) {
	ca, _ := s.db.GetRunningActivityByUserID(user.ID)
	now := time.Now()
	if ca != nil {
		if _, err := s.StopRunningActivity(user); err != nil {
			return nil, err
		}
	}
	orig := data.ActivityDB{
		Activity: pkg.Activity{
			UserID: user.ID,
			InputActivity: pkg.InputActivity{
				Start:       &now,
				Description: desc,
			},
		},
	}

	err := s.db.SaveActivity(&orig)
	if err != nil {
		return nil, err
	}
	hooked, modified := s.startActivityHook(&orig.Activity)
	log.Debug(hooked, modified)
	if modified {
		hooked.ID = orig.ID // ensure id is not changed
		orig.Activity = *hooked
		s.db.SaveActivity(&orig)
	}
	return hooked, nil
}

func (s *Service) AddActivity(a pkg.Activity, user pkg.User) (*pkg.Activity, error) {
	// handle activities without end as new started activities
	if a.End == nil {
		return s.StartActivity(a.Description, user)
	}

	err := s.db.SaveActivity(&data.ActivityDB{
		Activity: a,
	})
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (s *Service) StopRunningActivity(user pkg.User) (*pkg.Activity, error) {
	a, err := s.db.GetRunningActivityByUserID(user.ID)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	a.End = &now
	ac := s.endActivityHook(&a.Activity)
	a.Activity = *ac
	err = s.db.SaveActivity(a)
	if err != nil {
		return nil, err
	}
	return &a.Activity, nil
}

func (s *Service) GetActivity(id uint, user pkg.User) (*pkg.Activity, error) {
	a, err := s.db.GetActivity(id)
	if err != nil {
		return nil, err
	}

	if a != nil && a.UserID != user.ID {
		return nil, pkg.ErrPermissionDenied
	}

	return &a.Activity, nil
}

func (s *Service) GetActivities(start time.Time, end time.Time, user pkg.User) ([]pkg.Activity, error) {
	asDB, err := s.db.GetActivitiesBetweenStartAndEnd(start, end, user.ID)
	if err != nil {
		return nil, err
	}

	as := make([]pkg.Activity, len(asDB))
	for i := 0; i < len(as); i++ {
		as[i] = asDB[i].Activity
	}

	return as, nil
}

func (s *Service) UpdateActivity(a pkg.Activity, user pkg.User) (*pkg.Activity, error) {
	aDB, err := s.db.GetActivity(a.ID)
	if err != nil {
		return nil, err
	}

	aDB.Activity = a

	err = s.db.SaveActivity(aDB)
	if err != nil {
		return nil, err
	}
	// delete WorkingDay of the passing day to force recalculation
	now := time.Now()
	if !(a.Start.Year() == now.Year() && a.Start.Month() == now.Month() && a.Start.Day() == now.Day()) {
		err := s.db.DeleteWorkDay(utils.DayStart(*a.Start), a.UserID)
		if err != nil {
			return nil, err
		}
	}
	return &a, nil
}

func (s *Service) DelActivity(id uint, user pkg.User) error {
	a, err := s.GetActivity(id, user)
	if err != nil {
		return err
	}
	tx := s.db.Conn.Delete(a)
	return tx.Error
}
