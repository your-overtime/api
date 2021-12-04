package service

import (
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/your-overtime/api/internal/data"
	"github.com/your-overtime/api/pkg"
)

func (s *Service) GetWorkDays(start time.Time, end time.Time, user pkg.User) ([]pkg.WorkDay, error) {
	wdDBs, err := s.db.GetWorkDaysBetweenStartAndEnd(start, end, user.ID)
	if err != nil {
		log.Debug(err)
		return nil, err
	}

	wds := make([]pkg.WorkDay, len(wdDBs))
	for i := range wdDBs {
		wds[i] = wdDBs[i].WorkDay
	}

	return wds, nil
}

func (s *Service) AddWorkDay(w pkg.WorkDay, user pkg.User) (*pkg.WorkDay, error) {
	err := s.db.SaveWorkDay(&data.WorkDayDB{WorkDay: w})
	if err != nil {
		log.Debug(err)
		return nil, err
	}
	return &w, nil
}
