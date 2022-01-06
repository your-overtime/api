package service

import (
	"time"

	"github.com/your-overtime/api/internal/data"
	"github.com/your-overtime/api/pkg"
	"github.com/your-overtime/api/pkg/utils"

	log "github.com/sirupsen/logrus"
)

func (s *Service) SumActivityBetweenStartAndEndInMinutes(start time.Time, end time.Time) (int64, error) {
	activities, err := s.db.GetActivitiesBetweenStartAndEnd(start, end, s.user.ID)
	if err != nil {
		log.Debug(err)
		return 0, err
	}
	activityTimeInMinutes := int64(0)
	for _, a := range activities {
		// var diff time.Duration
		if a.End == nil {

			diff := end.Sub(*a.Start)
			activityTimeInMinutes += int64(diff.Minutes())
		} else {
			activityTimeInMinutes += int64(a.EventualDurationInMinutes)
		}

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

func (s *Service) StartActivity(desc string) (*pkg.Activity, error) {
	ca, _ := s.db.GetRunningActivityByUserID(s.user.ID)
	now := time.Now()
	if ca != nil {
		if _, err := s.StopRunningActivity(); err != nil {
			return nil, err
		}
	}
	orig := data.ActivityDB{
		Activity: pkg.Activity{
			UserID: s.user.ID,
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
	s.startActivityHook(&orig.Activity)

	return &orig.Activity, nil
}

func (s *Service) AddActivity(a pkg.Activity) (*pkg.Activity, error) {
	// handle activities without end as new started activities
	if a.End == nil {
		return s.StartActivity(a.Description)
	} else {
		s.calculateDuration(&a)
	}

	err := s.db.SaveActivity(&data.ActivityDB{
		Activity: a,
	})
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (s *Service) StopRunningActivity() (*pkg.Activity, error) {
	a, err := s.db.GetRunningActivityByUserID(s.user.ID)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	a.End = &now
	s.calculateDuration(&a.Activity)
	s.endActivityHook(&a.Activity)
	err = s.db.SaveActivity(a)
	if err != nil {
		return nil, err
	}
	return &a.Activity, nil
}

func (s *Service) GetActivity(id uint) (*pkg.Activity, error) {
	a, err := s.db.GetActivity(id, s.user.ID)
	if err != nil {
		return nil, err
	}

	if a != nil && a.UserID != s.user.ID {
		return nil, pkg.ErrPermissionDenied
	}

	return &a.Activity, nil
}

func (s *Service) GetActivities(start time.Time, end time.Time) ([]pkg.Activity, error) {
	asDB, err := s.db.GetActivitiesBetweenStartAndEnd(start, end, s.user.ID)
	if err != nil {
		return nil, err
	}

	as := make([]pkg.Activity, len(asDB))
	for i := 0; i < len(as); i++ {
		as[i] = asDB[i].Activity
	}

	return as, nil
}

func (s *Service) UpdateActivity(a pkg.Activity) (*pkg.Activity, error) {
	aDB, err := s.db.GetActivity(a.ID, s.user.ID)
	if err != nil {
		return nil, err
	}

	aDB.InputActivity = a.InputActivity
	if a.End != nil {
		s.calculateDuration(&aDB.Activity)
		s.endActivityHook(&aDB.Activity)
	}
	if err := s.db.SaveActivity(aDB); err != nil {
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
	return &aDB.Activity, nil
}

func (s *Service) DelActivity(id uint) error {
	a, err := s.db.GetActivity(id, s.user.ID)
	if err != nil {
		return err
	}
	return s.db.Conn.Delete(a).Error
}

func (s *Service) calculateDuration(a *pkg.Activity) {
	if a.End != nil {
		actualDuration := utils.DurationInMinutes(a.End.Sub(*a.Start))
		a.ActualDurationInMinutes = actualDuration
		if a.EventualDurationInMinutes == 0 {
			a.EventualDurationInMinutes = actualDuration
		}
	}
}
