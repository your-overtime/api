package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/emersion/go-ical"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func (a *API) ICalActivities(c *gin.Context) {
	yot, err := a.getOvertimeServiceForUserFromRequest(c)
	if err != nil {
		c.String(http.StatusUnauthorized, "401 Unauthorized")
	}

	now := time.Now()
	end := time.Date(now.Year(), 12, 31, 23, 59, 59, 0, now.Location())
	start := time.Date(now.Year()-1, 1, 1, 0, 0, 0, 0, now.Location())
	activitites, err := yot.GetActivities(start, end)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
	cal := ical.NewCalendar()
	cal.Props.SetText(ical.PropProductID, "-//Your Overtime//Activities")
	cal.Props.SetText(ical.PropVersion, "2.0")
	for _, ac := range activitites {
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
