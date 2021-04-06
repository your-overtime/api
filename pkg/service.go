package pkg

import "errors"

type OvertimeService interface {
	CalcCurrentOverview(e Employee) (*Overview, error)
	StartActivity(activity Activity, employee Employee) error
	StopRunningActivity(employee Employee) (Activity, error)
}

var (
	ErrUserNotFound       = errors.New("User not found")
	ErrInvalidCredentials = errors.New("Login or password are wrong")
)

type EmployeeService interface {
	FromToken(token string) (*Employee, error)
	Login(login string, password string) (*Employee, error)
}
