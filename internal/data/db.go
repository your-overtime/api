package data

import (
	"fmt"
	"net/url"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/your-overtime/api/pkg"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Db struct
type Db struct {
	Conn *gorm.DB
}

// Init function return Db
func Init(user string, pw string, host string, name string) (*Db, error) {
	db := Db{}

	tz := os.Getenv("TZ")
	if len(tz) == 0 {
		tz = time.Now().Location().String()
	}

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=%s",
		url.QueryEscape(user),
		url.QueryEscape(pw),
		host,
		url.QueryEscape(name),
		url.QueryEscape(tz),
	)

	conn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Debug(err)
		return nil, err
	}
	db.Conn = conn

	// Migrate the schema
	conn.AutoMigrate(&pkg.Activity{})
	conn.AutoMigrate(&pkg.Employee{})
	conn.AutoMigrate(&pkg.Token{})
	conn.AutoMigrate(&pkg.Holiday{})
	conn.AutoMigrate(&pkg.WorkDay{})
	conn.AutoMigrate(&pkg.WebhookInput{})
	db.MirgrateTokensToHashedTokens()
	return &db, err
}
