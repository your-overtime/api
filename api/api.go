package api

import (
	"encoding/base64"
	"errors"
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
}

func auth(accessTokens map[string]bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c)

		if _, exists := accessTokens[token]; !exists {
			c.AbortWithError(http.StatusUnauthorized, errors.New("Invalid token"))
		}
	}
}

func extractToken(c *gin.Context) string {
	token := ""
	authHeaderSlice := strings.Split(c.Request.Header.Get("Authorization"), " ")
	if len(authHeaderSlice) > 1 {
		token = authHeaderSlice[1]
	}

	if len(token) == 0 {
		token = c.Request.FormValue("token")
	}

	return strings.TrimSpace(token)
}

func tokenSliceToMap(accessTokens []string) map[string]bool {
	tokenMap := map[string]bool{}
	for _, token := range accessTokens {
		token := token
		tokenMap[token] = true
	}

	return tokenMap
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

		overview, err := a.os.CalcCurrentOverview()
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, traces)
		}
	})

	authorizedV1 := v1.Group("/", auth(s.accessTokens))
	{
		authorizedV1.POST("/traces", func(c *gin.Context) {
			traces := []pkg.Trace{}
			err := c.ShouldBindJSON(&traces)
			if err != nil {
				c.JSON(
					http.StatusBadRequest,
					map[string]string{"msg": "bad request"},
				)
			} else {
				err := s.traceService.SaveTraces(pkg.SaveTracesOptions{
					Traces: traces,
				})
				if err != nil {
					c.JSON(
						http.StatusInternalServerError,
						map[string]string{"msg": "something went wrong"},
					)
				} else {
					c.JSON(
						http.StatusCreated,
						map[string]string{"msg": "traces are saved"},
					)
				}
			}
		})
		authorizedV1.POST("/tracer", func(c *gin.Context) {
			tracer := pkg.Tracer{}
			err := c.ShouldBindJSON(&tracer)
			if err != nil {
				c.JSON(
					http.StatusBadRequest,
					map[string]string{"msg": "bad request"},
				)
			} else {
				tracer, err := s.traceService.SaveTracer(pkg.SaveTracerOptions{
					Tracer: tracer,
				})
				if err != nil {
					c.JSON(
						http.StatusInternalServerError,
						map[string]string{"msg": "something went wrong"},
					)
				} else {
					c.JSON(
						http.StatusCreated,
						tracer,
					)
				}
			}
		})
	}
}

// Init API server
func Init(host string, service pkg.TraceService, accessTokens []string) *Server {
	return &Server{
		router:       gin.Default(),
		host:         host,
		traceService: service,
		accessTokens: tokenSliceToMap(accessTokens),
	}
}

// Start API server
func (s Server) Start() {
	s.createEndPoints()
	s.router.Run(s.host)
}
