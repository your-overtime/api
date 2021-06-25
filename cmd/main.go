package main

import (
	"os"
	"strings"
	"time"

	"github.com/your-overtime/api/api"
	"github.com/your-overtime/api/internal/data"
	"github.com/your-overtime/api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	log "github.com/sirupsen/logrus"
)

func main() {
	godotenv.Load()
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
	log.SetReportCaller(true)
	if strings.ToLower(os.Getenv("DEBUG")) == "true" {
		log.SetLevel(log.DebugLevel)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	if len(os.Getenv("TZ")) > 0 {
		loc, err := time.LoadLocation(os.Getenv("TZ"))
		if err != nil {
			panic(err)
		}
		time.Local = loc
	}

	db, err := data.Init(
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"),
	)
	if err != nil {
		panic(err)
	}

	ovs := service.Init(db)
	api := api.Init(ovs, os.Getenv("ADMIN_TOKEN"))
	api.Start(os.Getenv("HOST"))
}
