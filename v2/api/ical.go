package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/emersion/go-ical"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// ICalActivities godoc
// @Tags activities.ics
// @Summary Get a activities by start and end
// @Produce text/calendar
// @Param start query string true "Start date" format date-time default 01.01 of the current year
// @Param end query string true "End date" format date-time default now
// @Success 200 file activities as ical
// @Router /activities.ics [get]
// @Security BasicAuth
// @Security ApiKeyAuth
func (a *API) ICalActivities(c *gin.Context) {
	yot, err := a.getOvertimeServiceForUserFromRequest(c)
	if err != nil {
		c.String(http.StatusUnauthorized, "401 Unauthorized")
	}

	var (
		start time.Time
		end   time.Time
	)

	now := time.Now()
	if len(c.Query("start")) > 0 {
		start, err = time.Parse(time.RFC3339, c.Query("start"))
		if err != nil {
			log.Debug(err)
			c.JSON(http.StatusBadRequest, err)
			return
		}
	} else {
		start = time.Date(now.Year()-1, 1, 1, 0, 0, 0, 0, now.Location())
	}

	if len(c.Query("end")) > 0 {
		end, err = time.Parse(time.RFC3339, c.Query("end"))
		if err != nil {
			log.Debug(err)
			c.JSON(http.StatusBadRequest, err)
			return
		}
	} else {
		end = time.Date(now.Year(), 12, 31, 23, 59, 59, 0, now.Location())
	}

	activities, err := yot.GetActivities(start, end)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
	cal := ical.NewCalendar()
	cal.Props.SetText(ical.PropProductID, "-//Your Overtime//Activities")
	cal.Props.SetText(ical.PropVersion, "2.0")
	for _, ac := range activities {
		event := ical.NewEvent()
		end := ac.End
		if end == nil {
			end = &now
		}
		event.Props.SetText(ical.PropUID, fmt.Sprintf("%d", ac.ID))
		event.Props.SetDateTime(ical.PropDateTimeStart, *ac.Start)
		event.Props.SetDateTime(ical.PropDateTimeEnd, *end)
		event.Props.SetDateTime(ical.PropDateTimeStamp, now)
		event.Props.SetText(ical.PropSummary, ac.Description)
		cal.Children = append(cal.Children, event.Component)
	}
	err = ical.NewEncoder(c.Writer).Encode(cal)
	if err != nil {
		log.Debugln(err)
	}
}

// ICalHolidays godoc
// @Tags holidays.ics
// @Summary Get a holidays by start and end
// @Produce text/calendar
// @Param start query string true "Start date" format date-time deafault 01.01 of the current year
// @Param end query string true "End date" format date-time default now
// @Success 200 file holidays as ical
// @Router /holidays.ics [get]
// @Security BasicAuth
// @Security ApiKeyAuth
func (a *API) ICalHolidays(c *gin.Context) {
	yot, err := a.getOvertimeServiceForUserFromRequest(c)
	if err != nil {
		c.String(http.StatusUnauthorized, "401 Unauthorized")
	}

	var (
		start time.Time
		end   time.Time
	)

	now := time.Now()
	if len(c.Query("start")) > 0 {
		start, err = time.Parse(time.RFC3339, c.Query("start"))
		if err != nil {
			log.Debug(err)
			c.JSON(http.StatusBadRequest, err)
			return
		}
	} else {
		start = time.Date(now.Year()-1, 1, 1, 0, 0, 0, 0, now.Location())
	}

	if len(c.Query("end")) > 0 {
		end, err = time.Parse(time.RFC3339, c.Query("end"))
		if err != nil {
			log.Debug(err)
			c.JSON(http.StatusBadRequest, err)
			return
		}
	} else {
		end = time.Date(now.Year(), 12, 31, 23, 59, 59, 0, now.Location())
	}

	holidays, err := yot.GetHolidays(start, end)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
	cal := ical.NewCalendar()
	cal.Props.SetText(ical.PropProductID, "-//Your Overtime//Holidays")
	cal.Props.SetText(ical.PropVersion, "2.0")
	for _, h := range holidays {
		event := ical.NewEvent()
		event.Props.SetText(ical.PropUID, fmt.Sprintf("%d", h.ID))
		event.Props.SetDateTime(ical.PropDateTimeStart, h.Start)
		event.Props.SetDateTime(ical.PropDateTimeEnd, h.End)
		event.Props.SetDateTime(ical.PropDateTimeStamp, now)
		event.Props.SetText(ical.PropSummary, fmt.Sprintf("%s (%s)", h.Description, h.Type))
		cal.Children = append(cal.Children, event.Component)
	}
	err = ical.NewEncoder(c.Writer).Encode(cal)
	if err != nil {
		log.Debugln(err)
	}
}
