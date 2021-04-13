package overtime

import (
	"errors"

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
	return nil, errors.New("not implemented yet")
}
func (s *service) StartActivity(activity pkg.Activity, employee pkg.Employee) error {
	return errors.New("not implemented yet")
}
func (s *service) StopRunningActivity(employee pkg.Employee) (*pkg.Activity, error) {
	return nil, errors.New("not implemented yet")
}

func (s *service) GetActivity(id string, employee pkg.Employee) (*pkg.Activity, error) {
	return nil, errors.New("not implemented yet")
}

func (s *service) DelActivity(id string, employee pkg.Employee) error {
	return errors.New("not implemented yet")
}

func (s *service) AddHollyday(h pkg.Hollyday, employee pkg.Employee) (*pkg.Hollyday, error) {
	return nil, errors.New("not implemented yet")
}

func (s *service) GetHollyday(id string, employee pkg.Employee) (*pkg.Hollyday, error) {
	return nil, errors.New("not implemented yet")
}

func (s *service) DelHollyday(id string, employee pkg.Employee) error {
	return errors.New("not implemented yet")
}
