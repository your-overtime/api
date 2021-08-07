package service

import (
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/your-overtime/api/pkg"
)

func (s *Service) GetWorkDays(start time.Time, end time.Time, employee pkg.Employee) ([]pkg.WorkDay, error) {
	wds, err := s.db.GetWorkDaysBetweenStartAndEnd(start, end, employee.ID)
	if err != nil {
		log.Debug(err)
		return nil, err
	}

	return wds, nil
}

func (s *Service) AddWorkDay(w pkg.WorkDay, employee pkg.Employee) (*pkg.WorkDay, error) {
	err := s.db.SaveWorkDay(&w)
	if err != nil {
		log.Debug(err)
		return nil, err
	}
	return &w, nil
}
