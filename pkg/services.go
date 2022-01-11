package pkg

import (
	"errors"
	"time"
)

type ActivityService interface {
	StartActivity(desc string) (*Activity, error)
	AddActivity(activity Activity) (*Activity, error)
	UpdateActivity(activity Activity) (*Activity, error)
	StopRunningActivity() (*Activity, error)
	GetActivity(id uint) (*Activity, error)
	GetActivities(start time.Time, end time.Time) ([]Activity, error)
	DelActivity(id uint) error
}

type HolidayService interface {
	AddHoliday(h Holiday) (*Holiday, error)
	UpdateHoliday(h Holiday) (*Holiday, error)
	GetHoliday(id uint) (*Holiday, error)
	GetHolidays(start time.Time, end time.Time) ([]Holiday, error)
	GetHolidaysByType(start time.Time, end time.Time, hType HolidayType) ([]Holiday, error)
	DelHoliday(id uint) error
}

type WorkDayService interface {
	GetWorkDays(start time.Time, end time.Time) ([]WorkDay, error)
	AddWorkDay(w WorkDay) (*WorkDay, error)
}

type UserService interface {
	SaveUser(user User, adminToken string) (*User, error)
	DeleteUser(login string, adminToken string) error

	UpdateAccount(fields map[string]interface{}, user User) (*User, error)
	GetAccount() (*User, error)

	CreateToken(token InputToken) (*Token, error)
	DeleteToken(tokenID uint) error
	GetTokens() ([]Token, error)
}

type WebhookService interface {
	CreateWebhook(webhook WebhookInput) (*Webhook, error)
	GetWebhooks() ([]Webhook, error)
}

type OvertimeService interface {
	CalcOverview(day time.Time) (*Overview, error)

	ActivityService
	HolidayService
	WorkDayService
	UserService
	WebhookService
}

type MainOvertimeService interface {
	GetOrCreateInstanceForUser(user *User) OvertimeService

	FromToken(token string) (*User, error)
	Login(login string, password string) (*User, error)
}

var (
	ErrUserNotFound               = errors.New("User not found")
	ErrInvalidCredentials         = errors.New("Login or password are wrong")
	ErrBadRequest                 = errors.New("Bad request")
	ErrPermissionDenied           = errors.New("Permission denied")
	ErrActivityIsRunning          = errors.New("A activity is currently running")
	ErrNoActivityIsRunning        = errors.New("No activity is currently running")
	ErrDuplicateValue             = errors.New("Duplicate value")
	ErrEmptyDescriptionNotAllowed = errors.New("empty description is not allowed")
)
