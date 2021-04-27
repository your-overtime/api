package data

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"git.goasum.de/jasper/overtime/pkg"
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

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user,
		pw,
		host,
		name,
	)

	fmt.Println(dsn)

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
	conn.AutoMigrate(&pkg.Hollyday{})
	conn.AutoMigrate(&pkg.WorkDay{})

	return &db, err
}
