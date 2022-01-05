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
} // @Name InputUser

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
} // @Name User

type InputToken struct {
	Name string
} // @Name InputToken

type Token struct {
	ID uint `gorm:"primaryKey"`
	InputToken
	UserID uint
	Token  string
} // @Name Token

func (e *User) WorkingDaysAsArray() []string {
	if len(e.WorkingDays) > 0 {
		return strings.Split(e.WorkingDays, ",")
	}
	return nil
}

type Activity struct {
	ID uint `gorm:"primaryKey"`
	InputActivity
	ActualDuration   time.Duration
	EventualDuration time.Duration
	UserID           uint
} // @Name Activity

type InputActivity struct {
	Start       *time.Time `format:"date-time" extensions:"x-nullable"`
	End         *time.Time `format:"date-time" extensions:"x-nullable"`
	Description string
} // @Name InputActivity

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
} // @Name Holiday

type InputHoliday struct {
	Start       time.Time `format:"date-time"`
	End         time.Time `format:"date-time"`
	Description string
	Type        HolidayType
} // @Name InputHoliday

type WorkDay struct {
	ID uint `gorm:"primaryKey"`
	InputWorkDay
	IsHoliday bool
} // @Name WorkDay

type InputWorkDay struct {
	Day        time.Time `gorm:"UNIQUE_INDEX:compositeindex;index;not null" format:"date-time"`
	Overtime   int64
	ActiveTime int64
	UserID     uint `gorm:"UNIQUE_INDEX:compositeindex;index;not null"`
} // @Name InputWorkDay

type Overview struct {
	Date                         time.Time `format:"date-time"`
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
	ActiveActivity               *Activity `extensions:"x-nullable"`
} // @Name Overtime

type WebhookInput struct {
	HeaderKey   string
	HeaderValue string
	TargetURL   string
	ReadOnly    bool
} // @Name WebhookInput

type Webhook struct {
	ID uint `gorm:"primaryKey"`
	WebhookInput
	UserID uint
} // @Name Webhook

type WebhookBody struct {
	Event   WebhookEvent
	Payload interface{}
}

type WebhookEvent string

const (
	StartActivityEvent WebhookEvent = "start_activity"
	EndActivityEvent   WebhookEvent = "end_activity"
)
