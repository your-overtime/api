package data

import (
	"fmt"

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

	conn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.Conn = conn

	// Migrate the schema
	conn.AutoMigrate(&pkg.Activity{})
	conn.AutoMigrate(&pkg.Employee{})
	conn.AutoMigrate(&pkg.Token{})
	conn.AutoMigrate(&pkg.Hollyday{})

	return &db, err
}
