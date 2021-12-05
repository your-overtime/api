package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/your-overtime/api/pkg"
)

// CreateUser godoc
// @Tags user
// @Summary creates a user
// @Produce json
// @Consume json
// @Param bottles body pkg.InputUser true "Input user"
// @Success 200 {object} pkg.User
// @Router /user [post]
// @Security AdminAuth
func (a *API) CreateUser(c *gin.Context) {
	var ie pkg.InputUser
	err := c.Bind(&ie)
	if err != nil {
		log.Debug(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}
	os := a.mos.GetOrCreateInstanceForUser(&pkg.User{})
	e, err := os.SaveUser(ie.ToUser(), "")
	if err != nil {
		log.Debug(err)
		c.JSON(http.StatusInternalServerError, err)
	} else {
		c.JSON(http.StatusCreated, e)
	}
}
