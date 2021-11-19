package pkg

import (
	"errors"
	"time"
)

type ActivityService interface {
	StartActivity(desc string, employee Employee) (*Activity, error)
	AddActivity(activity Activity, employee Employee) (*Activity, error)
	UpdateActivity(activity Activity, employee Employee) (*Activity, error)
	StopRunningActivity(employee Employee) (*Activity, error)
	GetActivity(id uint, employee Employee) (*Activity, error)
	GetActivities(start time.Time, end time.Time, employee Employee) ([]Activity, error)
	DelActivity(id uint, employee Employee) error
}

type HolidayService interface {
	AddHoliday(h Holiday, employee Employee) (*Holiday, error)
	UpdateHoliday(h Holiday, employee Employee) (*Holiday, error)
	GetHoliday(id uint, employee Employee) (*Holiday, error)
	GetHolidays(start time.Time, end time.Time, employee Employee) ([]Holiday, error)
	GetHolidaysByType(start time.Time, end time.Time, hType HolidayType, employee Employee) ([]Holiday, error)
	DelHoliday(id uint, employee Employee) error
}

type WorkDayService interface {
	GetWorkDays(start time.Time, end time.Time, employee Employee) ([]WorkDay, error)
	AddWorkDay(w WorkDay, employee Employee) (*WorkDay, error)
}

type EmployeeService interface {
	SaveEmployee(employee Employee, adminToken string) (*Employee, error)
	DeleteEmployee(login string, adminToken string) error

	UpdateAccount(fields map[string]interface{}, employee Employee) (*Employee, error)
	GetAccount() (*Employee, error)

	CreateToken(token InputToken, employee Employee) (*Token, error)
	DeleteToken(tokenID uint, employee Employee) error
	GetTokens(employee Employee) ([]Token, error)
}

type WebhookService interface {
	CreateWebhook(webhook Webhook) (*Webhook, error)
	GetWebhooks(employee Employee) ([]Webhook, error)
}
type OvertimeService interface {
	CalcOverview(e Employee, day time.Time) (*Overview, error)

	ActivityService
	HolidayService
	WorkDayService
	EmployeeService
	WebhookService
}

var (
	ErrUserNotFound        = errors.New("User not found")
	ErrInvalidCredentials  = errors.New("Login or password are wrong")
	ErrBadRequest          = errors.New("Bad request")
	ErrPermissionDenied    = errors.New("Permission denied")
	ErrActivityIsRunning   = errors.New("A activity is currently running")
	ErrNoActivityIsRunning = errors.New("No activity is currently running")
	ErrDuplicateValue      = errors.New("Duplicate value")
)
