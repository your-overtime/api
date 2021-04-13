package api

import (
	"encoding/base64"
	"net/http"
	"strings"

	"git.goasum.de/jasper/overtime/pkg"
	"github.com/gin-gonic/gin"
)

// API struct
type API struct {
	os     pkg.OvertimeService
	es     pkg.EmployeeService
	router *gin.Engine
	host   string
}

func (a *API) getEmployeeFromRequest(c *gin.Context) (*pkg.Employee, error) {
	token := c.Request.FormValue("token")
	if len(token) > 0 {
		return a.es.FromToken(token)
	}
	authHeaderSlice := strings.Split(c.Request.Header.Get("Authorization"), " ")
	if len(authHeaderSlice) == 3 {
		switch strings.ToLower(authHeaderSlice[1]) {
		case "basic":
			payload := []byte{}
			_, err := base64.StdEncoding.Decode(payload, []byte(authHeaderSlice[2]))
			if err != nil {
				return nil, pkg.ErrUserNotFound
			}
			basicAuth := strings.Split(string(payload), ":")
			if len(basicAuth) != 2 {
				return nil, pkg.ErrUserNotFound
			}
			return a.es.Login(basicAuth[0], basicAuth[1])
		default:
			return a.es.FromToken(authHeaderSlice[2])
		}

	}

	return nil, pkg.ErrUserNotFound
}

func (a *API) createEndPoints() {
	api := a.router.Group("/api")

	v1 := api.Group("/v1")
	v1.GET("/current", func(c *gin.Context) {
		e, err := a.getEmployeeFromRequest(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
		}
		overview, err := a.os.CalcCurrentOverview(*e)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, overview)
		}
	})
	v1.POST("/activity", func(c *gin.Context) {
		e, err := a.getEmployeeFromRequest(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
		}
		var ac pkg.Activity
		err = c.Bind(a)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
		}
		err = a.os.StartActivity(ac, *e)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, ac)
		}
	})
	v1.DELETE("/activity", func(c *gin.Context) {
		e, err := a.getEmployeeFromRequest(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
		}
		ac, err := a.os.StopRunningActivity(*e)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, ac)
		}
	})
	v1.GET("/activity/:id", func(c *gin.Context) {
		e, err := a.getEmployeeFromRequest(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
		}
		id := c.Param("id")
		h, err := a.os.GetActivity(id, *e)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, h)
		}
	})
	v1.DELETE("/activity/:id", func(c *gin.Context) {
		e, err := a.getEmployeeFromRequest(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
		}
		id := c.Param("id")
		err = a.os.DelActivity(id, *e)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, "")
		}
	})
	v1.POST("/hollyday", func(c *gin.Context) {
		e, err := a.getEmployeeFromRequest(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
		}
		var ho pkg.Hollyday
		err = c.Bind(a)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
		}
		h, err := a.os.AddHollyday(ho, *e)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, h)
		}
	})
	v1.GET("/hollyday/:id", func(c *gin.Context) {
		e, err := a.getEmployeeFromRequest(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
		}
		id := c.Param("id")
		h, err := a.os.GetHollyday(id, *e)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, h)
		}
	})
	v1.DELETE("/hollyday/:id", func(c *gin.Context) {
		e, err := a.getEmployeeFromRequest(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
		}
		id := c.Param("id")
		err = a.os.DelHollyday(id, *e)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, "")
		}
	})
}

// Init API server
func Init(os pkg.OvertimeService, es pkg.EmployeeService) *API {
	return &API{
		router: gin.Default(),
		os:     os,
		es:     es,
	}
}

// Start API server
func (a API) Start(host string) {
	a.createEndPoints()
	a.router.Run(host)
}
