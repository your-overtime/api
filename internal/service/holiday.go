package service

import (
	"errors"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/your-overtime/api/internal/data"
	"github.com/your-overtime/api/pkg"
	"github.com/your-overtime/api/pkg/utils"
)

func (s *Service) CountUsedHolidaysBetweenStartAndEnd(start time.Time, end time.Time, e pkg.User) (uint, error) {
	useHolidays := uint(0)
	hs, err := s.db.GetHolidaysBetweenStartAndEndByType(start, end, pkg.HolidayTypeFree, e.ID)
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

func (s *Service) SumHolidaysBetweenStartAndEndInMinutes(start time.Time, end time.Time, e pkg.User) (int64, bool, error) {
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

func (s *Service) AddHoliday(h pkg.Holiday, user pkg.User) (*pkg.Holiday, error) {
	err := s.db.SaveHoliday(&data.HolidayDB{Holiday: h})
	if err != nil {
		log.Debug(err)
		return nil, err
	}
	return &h, nil
}

func (s *Service) GetHoliday(id uint, user pkg.User) (*pkg.Holiday, error) {
	h, err := s.db.GetHoliday(id)
	if err != nil {
		log.Debug(err)
		return nil, err
	}

	if h.UserID != user.ID {
		return nil, pkg.ErrPermissionDenied
	}

	return &h.Holiday, nil
}

func (s *Service) UpdateHoliday(activity pkg.Holiday, user pkg.User) (*pkg.Holiday, error) {
	// only needed in client
	return nil, errors.New("not implemented")
}

func (s *Service) GetHolidays(start time.Time, end time.Time, user pkg.User) ([]pkg.Holiday, error) {
	hDBs, err := s.db.GetHolidaysBetweenStartAndEnd(start, end, user.ID)
	if err != nil {
		log.Debug(err)
		return nil, err
	}

	return castHolidayDBToPkgArray(hDBs), nil
}

func (s *Service) GetHolidaysByType(start time.Time, end time.Time, hType pkg.HolidayType, user pkg.User) ([]pkg.Holiday, error) {
	hDBs, err := s.db.GetHolidaysBetweenStartAndEndByType(start, end, hType, user.ID)
	if err != nil {
		log.Debug(err)
		return nil, err
	}

	return castHolidayDBToPkgArray(hDBs), nil
}

func (s *Service) DelHoliday(id uint, user pkg.User) error {
	h, err := s.GetHoliday(id, user)
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
