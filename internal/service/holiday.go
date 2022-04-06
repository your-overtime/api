package service

import (
	"errors"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/your-overtime/api/internal/data"
	"github.com/your-overtime/api/pkg"
	"github.com/your-overtime/api/pkg/utils"
)

func (s *Service) CountUsedHolidaysBetweenStartAndEnd(start time.Time, end time.Time) (uint, error) {
	useHolidays := uint(0)
	hs, err := s.db.GetHolidaysBetweenStartAndEndByType(start, end, pkg.HolidayTypeFree, s.user.ID)
	if err != nil {
		return 0, err
	}

	for _, h := range hs {
		h := h
		cs := utils.DayStart(h.Start)
		e := utils.DayEnd(h.End)

		useHolidays += uint(float64(e.Unix()-cs.Unix()) / 60 / 60 / 24)
	}

	return useHolidays, nil
}

func (s *Service) SumHolidaysBetweenStartAndEndInMinutes(start time.Time, end time.Time) (int64, bool, error) {
	isLegal := false
	holidays, err := s.db.GetHolidaysBetweenStartAndEnd(start, end, s.user.ID)
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
				dayFreeTimeInMinutes = int64(s.user.WeekWorkingTimeInMinutes / 5)
			} else {
				dayFreeTimeInMinutes = int64(s.user.WeekWorkingTimeInMinutes / s.user.NumWorkingDays)
			}
			freeTimeInMinutes += dayFreeTimeInMinutes
			st = st.AddDate(0, 0, 1)
		}
	}
	return freeTimeInMinutes, isLegal, nil
}

func (s *Service) AddHoliday(h pkg.Holiday) (*pkg.Holiday, error) {
	if s.readonly {
		return nil, pkg.ErrReadOnlyAccess
	}
	err := s.db.SaveHoliday(&data.HolidayDB{Holiday: h})
	if err != nil {
		log.Debug(err)
		return nil, err
	}
	return &h, nil
}

func (s *Service) GetHoliday(id uint) (*pkg.Holiday, error) {
	h, err := s.db.GetHoliday(id)
	if err != nil {
		log.Debug(err)
		return nil, err
	}

	if h.UserID != s.user.ID {
		return nil, pkg.ErrPermissionDenied
	}

	return &h.Holiday, nil
}

func (s *Service) UpdateHoliday(activity pkg.Holiday) (*pkg.Holiday, error) {
	// only needed in client
	return nil, errors.New("not implemented")
}

func (s *Service) GetHolidays(start time.Time, end time.Time) ([]pkg.Holiday, error) {
	hDBs, err := s.db.GetHolidaysBetweenStartAndEnd(start, end, s.user.ID)
	if err != nil {
		log.Debug(err)
		return nil, err
	}

	return castHolidayDBToPkgArray(hDBs), nil
}

func (s *Service) GetHolidaysByType(start time.Time, end time.Time, hType pkg.HolidayType) ([]pkg.Holiday, error) {
	hDBs, err := s.db.GetHolidaysBetweenStartAndEndByType(start, end, hType, s.user.ID)
	if err != nil {
		log.Debug(err)
		return nil, err
	}

	return castHolidayDBToPkgArray(hDBs), nil
}

func (s *Service) DelHoliday(id uint) error {
	if s.readonly {
		return pkg.ErrReadOnlyAccess
	}
	h, err := s.GetHoliday(id)
	if err != nil {
		log.Debug(err)
		return err
	}
	tx := s.db.Conn.Delete(h)
	return tx.Error
}

func castHolidayDBToPkgArray(hDBs []data.HolidayDB) []pkg.Holiday {
	hs := make([]pkg.Holiday, len(hDBs))
	for i := range hDBs {
		hs[i] = hDBs[i].Holiday
	}

	return hs
}
