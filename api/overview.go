package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// GetOverview godoc
// @Tags overview
// @Summary Retrieves overview of your overtime
// @Produce json
// @Success 200 {object} pkg.Overview
// @Router /overview [get]
// @Security BasicAuth
// @Security ApiKeyAuth
func (a *API) GetOverview(c *gin.Context) {
	e, err := a.getEmployeeFromRequest(c)
	if err != nil {
		log.Debug(err)
		c.JSON(http.StatusUnauthorized, err)
		return
	}
	overview, err := a.os.CalcOverview(*e, time.Now())
	if err != nil {
		log.Debug(err)
		c.JSON(http.StatusInternalServerError, err)
	} else {
		c.JSON(http.StatusOK, overview)
	}
}
