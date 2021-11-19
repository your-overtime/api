package pkg

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Model struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type InputEmployee struct {
	Name                     string
	Surname                  string
	Login                    string
	Password                 string
	WeekWorkingTimeInMinutes uint
	NumWorkingDays           uint
	NumHolidays              uint
}

func (u *InputEmployee) ToEmployee() Employee {
	return Employee{
		User: &User{
			Name:     u.Name,
			Surname:  u.Surname,
			Login:    u.Login,
			Password: u.Password,
		},
		WeekWorkingTimeInMinutes: u.WeekWorkingTimeInMinutes,
		NumWorkingDays:           u.NumWorkingDays,
		NumHolidays:              u.NumHolidays,
	}
}

type User struct {
	Model
	Name     string
	Surname  string
	Login    string `gorm:"unique"`
	Password string `json:"-"`
	Tokens   []Token
}

type InputToken struct {
	Name string
}

type Token struct {
	Model
	Name   string
	UserID uint
	Token  string
}

type Employee struct {
	*User                    `gorm:"embedded"`
	WeekWorkingTimeInMinutes uint
	NumWorkingDays           uint
	NumHolidays              uint
}

type Activity struct {
	Model
	Start       *time.Time
	End         *time.Time
	Description string
	UserID      uint
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
	Model
	Start       time.Time
	End         time.Time
	Description string
	Type        HolidayType
	UserID      uint
}

type InputHoliday struct {
	Start       time.Time
	End         time.Time
	Description string
	Type        HolidayType
}

type WorkDay struct {
	Model
	Day        time.Time `gorm:"UNIQUE_INDEX:compositeindex;index;not null"`
	Overtime   int64
	ActiveTime int64
	IsHoliday  bool
	UserID     uint `gorm:"UNIQUE_INDEX:compositeindex;index;not null"`
}

type InputWorkDay struct {
	Day        time.Time
	Overtime   int64
	ActiveTime int64
	UserID     uint
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

type Webhook struct {
	Model
	HeaderKey   string
	HeaderValue string
	TargetURL   string
	UserID      uint
	ReadOnly    bool
}
