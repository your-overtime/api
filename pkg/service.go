package pkg

import (
	"errors"
	"time"
)

type OvertimeService interface {
	CalcOverview(e Employee) (*Overview, error)
	StartActivity(desc string, employee Employee) (*Activity, error)
	AddActivity(activity Activity, employee Employee) (*Activity, error)
	UpdateActivity(activity Activity, employee Employee) (*Activity, error)
	StopRunningActivity(employee Employee) (*Activity, error)
	GetActivity(id uint, employee Employee) (*Activity, error)
	GetActivities(start time.Time, end time.Time, employee Employee) ([]Activity, error)
	DelActivity(id uint, employee Employee) error
	AddHollyday(h Hollyday, employee Employee) (*Hollyday, error)
	UpdateHollyday(h Hollyday, employee Employee) (*Hollyday, error)
	GetHollyday(id uint, employee Employee) (*Hollyday, error)
	GetHollydays(start time.Time, end time.Time, employee Employee) ([]Hollyday, error)
	DelHollyday(id uint, employee Employee) error

	SaveEmployee(employee Employee, adminToken string) (*Employee, error)
	DeleteEmployee(login string, adminToken string) error
	CreateToken(token InputToken, employee Employee) (*Token, error)
	DeleteToken(tokenID uint, employee Employee) error
	GetTokens(employee Employee) ([]Token, error)
}

var (
	ErrUserNotFound       = errors.New("User not found")
	ErrInvalidCredentials = errors.New("Login or password are wrong")
	ErrPermissionDenied   = errors.New("Permission denied")
	ErrActivityIsRunning  = errors.New("A activity is currently running")
)
