package pkg

import (
	"errors"
	"time"
)

type ActivityService interface {
	StartActivity(desc string, user User) (*Activity, error)
	AddActivity(activity Activity, user User) (*Activity, error)
	UpdateActivity(activity Activity, user User) (*Activity, error)
	StopRunningActivity(user User) (*Activity, error)
	GetActivity(id uint, user User) (*Activity, error)
	GetActivities(start time.Time, end time.Time, user User) ([]Activity, error)
	DelActivity(id uint, user User) error
}

type HolidayService interface {
	AddHoliday(h Holiday, user User) (*Holiday, error)
	UpdateHoliday(h Holiday, user User) (*Holiday, error)
	GetHoliday(id uint, user User) (*Holiday, error)
	GetHolidays(start time.Time, end time.Time, user User) ([]Holiday, error)
	GetHolidaysByType(start time.Time, end time.Time, hType HolidayType, user User) ([]Holiday, error)
	DelHoliday(id uint, user User) error
}

type WorkDayService interface {
	GetWorkDays(start time.Time, end time.Time, user User) ([]WorkDay, error)
	AddWorkDay(w WorkDay, user User) (*WorkDay, error)
}

type UserService interface {
	SaveUser(user User, adminToken string) (*User, error)
	DeleteUser(login string, adminToken string) error

	UpdateAccount(fields map[string]interface{}, user User) (*User, error)
	GetAccount() (*User, error)

	CreateToken(token InputToken, user User) (*Token, error)
	DeleteToken(tokenID uint, user User) error
	GetTokens(user User) ([]Token, error)
}

type WebhookService interface {
	CreateWebhook(webhook WebhookInput, user User) (*Webhook, error)
	GetWebhooks(user User) ([]Webhook, error)
}
type OvertimeService interface {
	CalcOverview(e User, day time.Time) (*Overview, error)

	ActivityService
	HolidayService
	WorkDayService
	UserService
	WebhookService
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
