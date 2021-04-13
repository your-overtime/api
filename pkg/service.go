package pkg

import "errors"

type OvertimeService interface {
	CalcCurrentOverview(e Employee) (*Overview, error)
	StartActivity(activity Activity, employee Employee) error
	StopRunningActivity(employee Employee) (*Activity, error)
	GetActivity(id string, employee Employee) (*Activity, error)
	DelActivity(id string, employee Employee) error
	AddHollyday(h Hollyday, employee Employee) (*Hollyday, error)
	GetHollyday(id string, employee Employee) (*Hollyday, error)
	DelHollyday(id string, employee Employee) error
}

var (
	ErrUserNotFound       = errors.New("User not found")
	ErrInvalidCredentials = errors.New("Login or password are wrong")
)

type EmployeeService interface {
	FromToken(token string) (*Employee, error)
	Login(login string, password string) (*Employee, error)
	AddEmployee(employee Employee) (*Employee, error)
}
