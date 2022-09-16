package api

import (
	"encoding/base64"
	"errors"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/your-overtime/api/v2/pkg"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// API struct
type API struct {
	mos        pkg.MainOvertimeService
	Router     *gin.Engine
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
	readonly := true
	if len(token) > 0 {
		user, err = a.mos.FromToken(token)
		readonly = a.mos.IsReadonlyToken(token)
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
				readonly = false
			default:
				user, err = a.mos.FromToken(authHeaderSlice[1])
				readonly = a.mos.IsReadonlyToken(authHeaderSlice[1])
			}

		}
	}

	if user != nil && err == nil {
		if readonly {
			return a.mos.GetOrCreateReadonlyInstanceForUser(user), nil
		}
		return a.mos.GetOrCreateInstanceForUser(user), nil
	}

	return nil, pkg.ErrUserNotFound
}

func (a *API) CreateEndpoints() {
	api := a.Router.Group("/api")

	v2 := api.Group("/v2")

	v2.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// Overview
	v2.GET("/overview", a.GetOverview)

	// activity
	v2.PATCH("/activity/stop", a.StopActivity)
	v2.POST("/activity", a.CreateActivity)
	v2.PUT("/activity/:id", a.UpdateActivity)
	v2.GET("/activity/:id", a.GetActivity)
	v2.GET("/activity", a.GetActivities)
	v2.DELETE("/activity/:id", a.DeleteActivity)

	// holiday
	v2.POST("/holiday", a.CreateHoliday)
	v2.PUT("/holiday/:id", a.UpdateHoliday)
	v2.GET("/holiday/:id", a.GetHoliday)
	v2.GET("/holiday", a.GetHolidays)
	v2.DELETE("/holiday/:id", a.DeleteHoliday)

	// workday
	v2.GET("/workday", a.GetWorkDays)
	v2.POST("/workday", a.CreateWorkDay)

	// token
	v2.GET("/token", a.GetTokens)
	v2.POST("/token", a.CreateToken)
	v2.DELETE("/token/:id", a.DeleteToken)

	// account
	v2.GET("/account", a.GetAccount)
	v2.PATCH("/account", a.UpdateAccount)

	// webhook
	v2.POST("/webhook", a.CreateWebhook)
	v2.GET("/webhook", a.GetWebhooks)

	// user
	authorizedV2 := v2.Group("/", a.adminAuth())
	{
		authorizedV2.POST("/user", a.CreateUser)
	}

	// ical / caldav
	v2.GET("/activities.ics", a.ICalActivities)
	v2.GET("/holidays.ics", a.ICalHolidays)
}

// Init API server
func Init(mos pkg.MainOvertimeService, adminToken string) *API {
	return &API{
		Router:     gin.Default(),
		mos:        mos,
		adminToken: adminToken,
	}
}

func (a API) Start(host string) {
	a.CreateEndpoints()
	panic(a.Router.Run(host))
}
