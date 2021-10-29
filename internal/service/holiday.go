package service

import (
	"errors"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/your-overtime/api/pkg"
)

func (s *Service) CountHolidaysBetweenStartAndEnd(start time.Time, end time.Time, e pkg.Employee) (uint, error) {
	return s.db.CountHolidaysBetweenStartAndEnd(start, end, e.ID)
}

func (s *Service) SumHolidaysBetweenStartAndEndInMinutes(start time.Time, end time.Time, e pkg.Employee) (int64, bool, error) {
	isLegal := false
	holidays, err := s.db.GetHolidaysBetweenStartAndEnd(start, end, e.ID)
	if err != nil {
		log.Debug(err)
		return 0, isLegal, err
	}
	freeTimeInMinutes := int64(0)
	for _, a := range holidays {
		st := a.Start
		if a.Start.Unix() < start.Unix() {
			st = start
		}

		for {
			if st.Unix() > end.Unix() || st.Unix() > a.End.Unix() {
				break
			}
			dayFreeTimeInMinutes := int64(0)
			if a.Type == pkg.HolidayTypeLegalHoliday {
				// Fix legal holidays
				isLegal = true
				dayFreeTimeInMinutes = int64(e.WeekWorkingTimeInMinutes / 5)
			} else {
				dayFreeTimeInMinutes = int64(e.WeekWorkingTimeInMinutes / e.NumWorkingDays)
			}
			freeTimeInMinutes += dayFreeTimeInMinutes
			st = st.AddDate(0, 0, 1)
		}
	}
	return freeTimeInMinutes, isLegal, nil
}

func (s *Service) AddHoliday(h pkg.Holiday, employee pkg.Employee) (*pkg.Holiday, error) {
	err := s.db.SaveHoliday(&h)
	if err != nil {
		log.Debug(err)
		return nil, err
	}
	return &h, nil
}

func (s *Service) GetHoliday(id uint, employee pkg.Employee) (*pkg.Holiday, error) {
	h, err := s.db.GetHoliday(id)
	if err != nil {
		log.Debug(err)
		return nil, err
	}

	if h.UserID != employee.ID {
		return nil, pkg.ErrPermissionDenied
	}

	return h, nil
}

func (s *Service) UpdateHoliday(activity pkg.Holiday, employee pkg.Employee) (*pkg.Holiday, error) {
	// only needed in client
	return nil, errors.New("not implemented")
}

func (s *Service) GetHolidays(start time.Time, end time.Time, employee pkg.Employee) ([]pkg.Holiday, error) {
	h, err := s.db.GetHolidaysBetweenStartAndEnd(start, end, employee.ID)
	if err != nil {
		log.Debug(err)
		return nil, err
	}

	return h, nil
}

func (s *Service) GetHolidaysByType(start time.Time, end time.Time, hType pkg.HolidayType, employee pkg.Employee) ([]pkg.Holiday, error) {
	h, err := s.db.GetHolidaysBetweenStartAndEndByType(start, end, hType, employee.ID)
	if err != nil {
		log.Debug(err)
		return nil, err
	}

	return h, nil
}

func (s *Service) DelHoliday(id uint, employee pkg.Employee) error {
	h, err := s.GetHoliday(id, employee)
	if err != nil {
		log.Debug(err)
		return err
	}
	tx := s.db.Conn.Delete(h)
	return tx.Error
}
