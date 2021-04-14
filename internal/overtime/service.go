package overtime

import (
	"errors"
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

func (s *service) CalcCurrentOverview(e pkg.Employee) (*pkg.Overview, error) {
	now := time.Now()
	_, week := now.ISOWeek()
	o := &pkg.Overview{
		Date:       time.Now(),
		WeekNumber: week,
		Employee:   &e,
	}
	cra, err := s.db.GetRunningActivityByEmployeeID(e.ID)
	if err == nil {
		o.ActiveActivity = cra
	}

	// TODO: sum working hours (maybe for the running year) and subtract e.WeekWorkingTime per week and hollydays

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
	return nil, errors.New("not implemented yet")
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
