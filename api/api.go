package api

import (
	"encoding/base64"
	"errors"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/your-overtime/api/pkg"

	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// API struct
type API struct {
	mos        pkg.MainOvertimeService
	router     *gin.Engine
	adminToken string
}

func (a *API) adminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.FormValue("adminToken")
		log.Debug(token, " ", a.adminToken, " ", token == a.adminToken)
		if token != a.adminToken {
			c.AbortWithError(http.StatusUnauthorized, errors.New("invalid token"))
		}
	}
}

func (a *API) getOvertimeServiceForUserFromRequest(c *gin.Context) (pkg.OvertimeService, error) {
	token := c.Request.FormValue("token")
	var (
		user *pkg.User
		err  error
	)
	if len(token) > 0 {
		user, err = a.mos.FromToken(token)
	} else {
		authHeaderSlice := strings.Split(c.Request.Header.Get("Authorization"), " ")
		if len(authHeaderSlice) == 2 {
			switch strings.ToLower(authHeaderSlice[0]) {
			case "basic":
				payload, payloadErr := base64.StdEncoding.DecodeString(authHeaderSlice[1])
				if payloadErr != nil {
					log.Debug(payloadErr)
					return nil, pkg.ErrUserNotFound
				}
				basicAuth := strings.Split(string(payload), ":")
				if len(basicAuth) != 2 {
					log.Debug(pkg.ErrUserNotFound, " ", basicAuth, " ", authHeaderSlice)
					return nil, pkg.ErrUserNotFound
				}
				user, err = a.mos.Login(basicAuth[0], basicAuth[1])
			default:
				user, err = a.mos.FromToken(authHeaderSlice[1])
			}

		}
	}

	if user != nil && err == nil {
		return a.mos.GetOrCreateInstanceForUser(user), nil
	}

	return nil, pkg.ErrUserNotFound
}

func (a *API) createEndPoints() {
	api := a.router.Group("/api")

	v1 := api.Group("/v1")

	v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Overview
	v1.GET("/overview", a.GetOverview)

	// activity
	v1.POST("/activity/:desc", a.StartActivity)
	v1.DELETE("/activity", a.StopActivity)
	v1.POST("/activity", a.CreateActivity)
	v1.PUT("/activity/:id", a.UpdateActivity)
	v1.GET("/activity/:id", a.GetActivity)
	v1.GET("/activity", a.GetActivities)
	v1.DELETE("/activity/:id", a.DeleteActivity)

	// holiday
	v1.POST("/holiday", a.CreateHoliday)
	v1.PUT("/holiday/:id", a.UpdateHoliday)
	v1.GET("/holiday/:id", a.GetHoliday)
	v1.GET("/holiday", a.GetHolidays)
	v1.DELETE("/holiday/:id", a.DeleteHoliday)

	// workday
	v1.GET("/workday", a.GetWorkDays)
	v1.POST("/workday", a.CreateWorkDay)

	// token
	v1.GET("/token", a.GetTokens)
	v1.POST("/token", a.CreateToken)
	v1.DELETE("/token/:id", a.DeleteToken)

	// account
	v1.GET("account", a.GetAccount)
	v1.PATCH("account", a.UpdateAccount)

	// webhook
	v1.POST("/webhook", a.CreateWebhook)
	v1.GET("/webhook", a.GetWebhooks)

	// user
	authorizedV1 := v1.Group("/", a.adminAuth())
	{
		authorizedV1.POST("/user", a.CreateUser)
	}
}

// Init API server
func Init(mos pkg.MainOvertimeService, adminToken string) *API {
	return &API{
		router:     gin.Default(),
		mos:        mos,
		adminToken: adminToken,
	}
}

func (a API) Start(host string) {
	a.createEndPoints()
	panic(a.router.Run(host))
}
