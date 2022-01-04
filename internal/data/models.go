package data

import (
	"time"

	"github.com/your-overtime/api/pkg"
	"gorm.io/gorm"
)

type ModelExtensions struct {
	CreatedAt time.Time
	UpdatedAt time.Time      `gorm:"index" json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type UserDB struct {
	ModelExtensions `gorm:"embedded"`
	pkg.User        `gorm:"embedded"`
}

func (UserDB) TableName() string {
	return "employees"
}

type TokenDB struct {
	ModelExtensions `gorm:"embedded"`
	pkg.Token       `gorm:"embedded"`
}

func (TokenDB) TableName() string {
	return "tokens"
}

type ActivityDB struct {
	ModelExtensions `gorm:"embedded"`
	pkg.Activity    `gorm:"embedded"`
}

func (ActivityDB) TableName() string {
	return "activities"
}

type HolidayDB struct {
	ModelExtensions `gorm:"embedded"`
	pkg.Holiday     `gorm:"embedded"`
}

func (HolidayDB) TableName() string {
	return "holidays"
}

type WorkDayDB struct {
	ModelExtensions `gorm:"embedded"`
	pkg.WorkDay     `gorm:"embedded"`
}

func (WorkDayDB) TableName() string {
	return "work_days"
}

type WebhookDB struct {
	ModelExtensions `gorm:"embedded"`
	pkg.Webhook     `gorm:"embedded"`
}

func (WebhookDB) TableName() string {
	return "webhooks"
}
