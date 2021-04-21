package api

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"git.goasum.de/jasper/overtime/internal/service"
	"git.goasum.de/jasper/overtime/pkg"
	"github.com/gin-gonic/gin"
)

// API struct
type API struct {
	os         *service.Service
	router     *gin.Engine
	host       string
	adminToken string
}

func (a *API) adminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.FormValue("adminToken")
		log.Debug(token, " ", a.adminToken, " ", token == a.adminToken)
		if token != a.adminToken {
			c.AbortWithError(http.StatusUnauthorized, errors.New("Invalid token"))
		}
	}
}

func (a *API) getEmployeeFromRequest(c *gin.Context) (*pkg.Employee, error) {
	token := c.Request.FormValue("token")
	if len(token) > 0 {
		return a.os.FromToken(token)
	}
	authHeaderSlice := strings.Split(c.Request.Header.Get("Authorization"), " ")
	if len(authHeaderSlice) == 2 {
		switch strings.ToLower(authHeaderSlice[0]) {
		case "basic":
			payload, err := base64.StdEncoding.DecodeString(authHeaderSlice[1])
			if err != nil {
				log.Debug(err)
				return nil, pkg.ErrUserNotFound
			}
			basicAuth := strings.Split(string(payload), ":")
			if len(basicAuth) != 2 {
				log.Debug(pkg.ErrUserNotFound, " ", basicAuth, " ", authHeaderSlice)
				return nil, pkg.ErrUserNotFound
			}
			return a.os.Login(basicAuth[0], basicAuth[1])
		default:
			return a.os.FromToken(authHeaderSlice[1])
		}

	}

	return nil, pkg.ErrUserNotFound
}

func (a *API) createEndPoints() {
	api := a.router.Group("/api")

	v1 := api.Group("/v1")
	v1.GET("overview", func(c *gin.Context) {
		e, err := a.getEmployeeFromRequest(c)
		if err != nil {
			log.Debug(err)
			c.JSON(http.StatusBadRequest, err)
			return
		}
		overview, err := a.os.CalcOverview(*e)
		if err != nil {
			log.Debug(err)
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, overview)
		}
	})
	v1.POST("/activity/:desc", func(c *gin.Context) {
		e, err := a.getEmployeeFromRequest(c)
		if err != nil {
			log.Debug(err)
			c.JSON(http.StatusBadRequest, err)
			return
		}
		fmt.Println(e)
		desc := c.Param("desc")
		ac, err := a.os.StartActivity(desc, *e)
		if err == pkg.ErrActivityIsRunning {
			c.JSON(http.StatusConflict, err.Error())
		} else if err != nil {
			log.Debug(err)
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, ac)
		}
	})
	v1.DELETE("/activity", func(c *gin.Context) {
		e, err := a.getEmployeeFromRequest(c)
		if err != nil {
			log.Debug(err)
			c.JSON(http.StatusBadRequest, err)
			return
		}
		ac, err := a.os.StopRunningActivity(*e)
		if err != nil {
			log.Debug(err)
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, ac)
		}
	})
	v1.POST("/activity", func(c *gin.Context) {
		e, err := a.getEmployeeFromRequest(c)
		if err != nil {
			log.Debug(err)
			c.JSON(http.StatusBadRequest, err)
			return
		}
		var ia pkg.InputActivity
		err = c.Bind(&ia)
		if err != nil {
			log.Debug(err)
			c.JSON(http.StatusBadRequest, err)
			return
		}
		act := pkg.Activity{
			UserID:      e.ID,
			Start:       ia.Start,
			End:         ia.End,
			Description: ia.Description,
		}
		ac, err := a.os.AddActivity(act, *e)
		if err != nil {
			log.Debug(err)
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, ac)
		}
	})
	v1.GET("/activity/:id", func(c *gin.Context) {
		e, err := a.getEmployeeFromRequest(c)
		if err != nil {
			log.Debug(err)
			c.JSON(http.StatusBadRequest, err)
			return
		}
		id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
		h, err := a.os.GetActivity(uint(id), *e)
		if err != nil {
			log.Debug(err)
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, h)
		}
	})
	v1.GET("/activity", func(c *gin.Context) {
		e, err := a.getEmployeeFromRequest(c)
		if err != nil {
			log.Debug(err)
			c.JSON(http.StatusBadRequest, err)
			return
		}
		start, err := time.Parse(time.RFC3339, c.Query("start"))
		if err != nil {
			log.Debug(err)
			c.JSON(http.StatusBadRequest, err)
			return
		}
		end, err := time.Parse(time.RFC3339, c.Query("end"))
		if err != nil {
			log.Debug(err)
			c.JSON(http.StatusBadRequest, err)
			return
		}
		h, err := a.os.GetActivities(start, end, *e)
		if err != nil {
			log.Debug(err)
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, h)
		}
	})
	v1.DELETE("/activity/:id", func(c *gin.Context) {
		e, err := a.getEmployeeFromRequest(c)
		if err != nil {
			log.Debug(err)
			c.JSON(http.StatusBadRequest, err)
			return
		}
		id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
		err = a.os.DelActivity(uint(id), *e)
		if err != nil {
			log.Debug(err)
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, "")
		}
	})
	v1.POST("/hollyday", func(c *gin.Context) {
		e, err := a.getEmployeeFromRequest(c)
		if err != nil {
			log.Debug(err)
			c.JSON(http.StatusBadRequest, err)
			return
		}
		var ih pkg.InputHollyday
		err = c.Bind(ih)
		if err != nil {
			log.Debug(err)
			c.JSON(http.StatusBadRequest, err)
			return
		}
		ho := pkg.Hollyday{
			UserID:      e.ID,
			Start:       ih.Start,
			End:         ih.End,
			Description: ih.Description,
		}
		h, err := a.os.AddHollyday(ho, *e)
		if err != nil {
			log.Debug(err)
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, h)
		}
	})
	v1.GET("/hollyday/:id", func(c *gin.Context) {
		e, err := a.getEmployeeFromRequest(c)
		if err != nil {
			log.Debug(err)
			c.JSON(http.StatusBadRequest, err)
			return
		}
		id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
		h, err := a.os.GetHollyday(uint(id), *e)
		if err != nil {
			log.Debug(err)
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, h)
		}
	})
	v1.GET("/hollyday", func(c *gin.Context) {
		e, err := a.getEmployeeFromRequest(c)
		if err != nil {
			log.Debug(err)
			c.JSON(http.StatusBadRequest, err)
			return
		}
		start, err := time.Parse(time.RFC3339Nano, c.Query("start"))
		if err != nil {
			log.Debug(start, err)
			c.JSON(http.StatusBadRequest, err)
			return
		}
		end, err := time.Parse(time.RFC3339Nano, c.Query("end"))
		if err != nil {
			log.Debug(end, err)
			c.JSON(http.StatusBadRequest, err)
			return
		}
		h, err := a.os.GetHollydays(start, end, *e)
		if err != nil {
			log.Debug(err)
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, h)
		}
	})
	v1.DELETE("/hollyday/:id", func(c *gin.Context) {
		e, err := a.getEmployeeFromRequest(c)
		if err != nil {
			log.Debug(err)
			c.JSON(http.StatusBadRequest, err)
			return
		}
		id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
		err = a.os.DelHollyday(uint(id), *e)
		if err != nil {
			log.Debug(err)
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, "")
		}
	})
	v1.POST("/token", func(c *gin.Context) {
		e, err := a.getEmployeeFromRequest(c)
		if err != nil {
			log.Debug(err)
			c.JSON(http.StatusBadRequest, err)
			return
		}
		var it pkg.InputToken
		err = c.Bind(&it)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}
		t, err := a.os.CreateToken(it, *e)
		if err != nil {
			log.Debug(err)
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusCreated, t)
		}
	})
	v1.GET("/token", func(c *gin.Context) {
		e, err := a.getEmployeeFromRequest(c)
		if err != nil {
			log.Debug(err)
			c.JSON(http.StatusBadRequest, err)
			return
		}

		ts, err := a.os.GetTokens(*e)
		if err != nil {
			log.Debug(err)
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusCreated, ts)
		}
	})
	v1.DELETE("/token/:id", func(c *gin.Context) {
		e, err := a.getEmployeeFromRequest(c)
		if err != nil {
			log.Debug(err)
			c.JSON(http.StatusBadRequest, err)
			return
		}
		id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
		err = a.os.DeleteToken(uint(id), *e)
		if err != nil {
			log.Debug(err)
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, "token deleted")
		}
	})
	authorizedV1 := v1.Group("/", a.adminAuth())
	{
		authorizedV1.POST("/employee", func(c *gin.Context) {
			var ie pkg.InputEmployee
			err := c.Bind(&ie)
			if err != nil {
				log.Debug(err)
				c.JSON(http.StatusBadRequest, err)
				return
			}
			e, err := a.os.SaveEmployee(ie.ToEmployee())
			if err != nil {
				log.Debug(err)
				c.JSON(http.StatusInternalServerError, err)
			} else {
				c.JSON(http.StatusCreated, e)
			}
		})
	}
}

// Init API server
func Init(os *service.Service, adminToken string) *API {
	return &API{
		router:     gin.Default(),
		os:         os,
		adminToken: adminToken,
	}
}

// Start API server
func (a API) Start(host string) {
	a.createEndPoints()
	panic(a.router.Run(host))
}
