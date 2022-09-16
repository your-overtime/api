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
// @Param date query string false "Calculation date" format date-time
// @Security BasicAuth
// @Security ApiKeyAuth
func (a *API) GetOverview(c *gin.Context) {
	os, err := a.getOvertimeServiceForUserFromRequest(c)
	if err != nil {
		log.Debug(err)
		c.JSON(http.StatusUnauthorized, err)
		return
	}

	strDate := c.Query("date")
	var date time.Time
	if len(strDate) == 0 {
		date = time.Now()
	} else {
		date, err = time.Parse(time.RFC3339, strDate)
		if err != nil {
			log.Debug(err)
			c.JSON(http.StatusBadRequest, err)
			return
		}
	}

	overview, err := os.CalcOverview(date)
	if err != nil {
		log.Debug(err)
		c.JSON(http.StatusInternalServerError, err)
	} else {
		c.JSON(http.StatusOK, overview)
	}
}
