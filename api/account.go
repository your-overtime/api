package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/your-overtime/api/pkg"
)

// GetWorkdays godoc
// @Tags account
// @Summary Retrieves account information
// @Produce json
// @Success 200 {object} pkg.User
// @Router /account [get]
// @Security BasicAuth
// @Security ApiKeyAuth
func (a *API) GetAccount(c *gin.Context) {
	os, err := a.getOvertimeServiceForUserFromRequest(c)
	if err != nil {
		log.Debug(err)
		c.JSON(http.StatusUnauthorized, err)
		return
	} else {
		user, err := os.GetAccount()
		if err != nil {
			log.Debug(err)
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, user)
	}
}

// UpdateAccount godoc
// @Tags account
// @Summary updates a account
// @Produce json
// @Consume json
// @Param account body map[string]interface{} true "input account fields"
// @Success 200 {object} pkg.User
// @Router /account [patch]
// @Security BasicAuth
// @Security ApiKeyAuth
func (a *API) UpdateAccount(c *gin.Context) {
	os, err := a.getOvertimeServiceForUserFromRequest(c)
	if err != nil {
		log.Debug(err)
		c.JSON(http.StatusUnauthorized, err)
		return
	}
	var payload map[string]interface{}
	err = c.Bind(&payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	user, _ := os.GetAccount()
	u, err := os.UpdateAccount(payload, *user)
	if err != nil {
		log.Debug(err)
		if errors.Is(err, pkg.ErrDuplicateValue) {
			c.JSON(http.StatusBadRequest, err)
		} else {
			c.JSON(http.StatusInternalServerError, err)
		}
	} else {
		c.JSON(http.StatusOK, u)
	}
}
