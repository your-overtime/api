package pkg

import (
	"errors"
	"strings"
	"time"
)

type InputUser struct {
	Name                     string
	Surname                  string
	Login                    string
	Password                 string
	WorkingDays              string
	WeekWorkingTimeInMinutes uint
	NumWorkingDays           uint
	NumHolidays              uint
}

func (u *InputUser) ToUser() User {
	return User{
		Name:                     u.Name,
		Surname:                  u.Surname,
		Login:                    u.Login,
		Password:                 u.Password,
		WorkingDays:              u.WorkingDays,
		WeekWorkingTimeInMinutes: u.WeekWorkingTimeInMinutes,
		NumWorkingDays:           u.NumWorkingDays,
		NumHolidays:              u.NumHolidays,
	}
}

type User struct {
	ID                       uint `gorm:"primaryKey"`
	Tokens                   []Token
	Password                 string `json:"-"`
	Name                     string
	Surname                  string
	Login                    string `gorm:"unique"`
	WorkingDays              string
	WeekWorkingTimeInMinutes uint
	NumWorkingDays           uint
	NumHolidays              uint
}

type InputToken struct {
	Name string
}

type Token struct {
	ID uint `gorm:"primaryKey"`
	InputToken
	UserID uint
	Token  string
}

func (e *User) WorkingDaysAsArray() []string {
	if len(e.WorkingDays) > 0 {
		return strings.Split(e.WorkingDays, ",")
	}
	return nil
}

type Activity struct {
	ID uint `gorm:"primaryKey"`
	InputActivity
	UserID uint
}

type InputActivity struct {
	Start       *time.Time
	End         *time.Time
	Description string
}

type HolidayType string

const (
	HolidayTypeFree            HolidayType = "free"
	HolidayTypeSick            HolidayType = "sick"
	HolidayTypeLegalHoliday    HolidayType = "legal_holiday"
	HolidayTypeLegalUnpaidFree HolidayType = "unpaid_free"
)

var ErrInvalidType = errors.New("invalid type")

func StrToHolidayType(str string) (HolidayType, error) {
	switch strings.ToLower(str) {
	case "free":
		return HolidayTypeFree, nil
	case "sick":
		return HolidayTypeSick, nil
	case "legal_holiday":
		return HolidayTypeLegalHoliday, nil
	case "unpaid_free":
		return HolidayTypeLegalUnpaidFree, nil
	}

	return HolidayTypeFree, ErrInvalidType
}

type Holiday struct {
	ID uint `gorm:"primaryKey"`
	InputHoliday
	UserID uint
}

type InputHoliday struct {
	Start       time.Time
	End         time.Time
	Description string
	Type        HolidayType
}

type WorkDay struct {
	ID uint `gorm:"primaryKey"`
	InputWorkDay
	IsHoliday bool
}

type InputWorkDay struct {
	Day        time.Time `gorm:"UNIQUE_INDEX:compositeindex;index;not null"`
	Overtime   int64
	ActiveTime int64
	UserID     uint `gorm:"UNIQUE_INDEX:compositeindex;index;not null"`
}

type Overview struct {
	Date                         time.Time
	WeekNumber                   int
	UsedHolidays                 int
	HolidaysStillAvailable       int
	OvertimeThisDayInMinutes     int64
	ActiveTimeThisDayInMinutes   int64
	OvertimeThisWeekInMinutes    int64
	ActiveTimeThisWeekInMinutes  int64
	OvertimeThisMonthInMinutes   int64
	ActiveTimeThisMonthInMinutes int64
	OvertimeThisYearInMinutes    int64
	ActiveTimeThisYearInMinutes  int64
	ActiveActivity               *Activity
}

type WebhookInput struct {
	HeaderKey   string
	HeaderValue string
	TargetURL   string
	ReadOnly    bool
}

type Webhook struct {
	ID uint `gorm:"primaryKey"`
	WebhookInput
	UserID uint
}

type WebhookBody struct {
	Event   WebhookEvent
	Payload interface{}
}

type WebhookEvent string

const (
	StartActivityEvent WebhookEvent = "start_activity"
	EndActivityEvent   WebhookEvent = "end_activity"
)
