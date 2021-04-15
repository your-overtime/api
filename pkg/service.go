package pkg

import "errors"

type OvertimeService interface {
	CalcCurrentOverview(e Employee) (*Overview, error)
	CalcOverviewForThisYear(e Employee) (*Overview, error)
	StartActivity(desc string, employee Employee) (*Activity, error)
	StopRunningActivity(employee Employee) (*Activity, error)
	GetActivity(id uint, employee Employee) (*Activity, error)
	DelActivity(id uint, employee Employee) error
	AddHollyday(h Hollyday, employee Employee) (*Hollyday, error)
	GetHollyday(id uint, employee Employee) (*Hollyday, error)
	DelHollyday(id uint, employee Employee) error
}

var (
	ErrUserNotFound       = errors.New("User not found")
	ErrInvalidCredentials = errors.New("Login or password are wrong")
	ErrPermissionDenied   = errors.New("Permission denied")
)

type EmployeeService interface {
	FromToken(token string) (*Employee, error)
	Login(login string, password string) (*Employee, error)
	AddEmployee(employee Employee) (*Employee, error)
}
