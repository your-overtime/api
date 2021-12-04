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
	e, err := a.getUserFromRequest(c)
	if err != nil {
		log.Debug(err)
		c.JSON(http.StatusUnauthorized, err)
		return
	} else {
		c.JSON(http.StatusOK, e)
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
	e, err := a.getUserFromRequest(c)
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
	e, err = a.os.UpdateAccount(payload, *e)
	if err != nil {
		log.Debug(err)
		if errors.Is(err, pkg.ErrDuplicateValue) {
			c.JSON(http.StatusBadRequest, err)
		} else {
			c.JSON(http.StatusInternalServerError, err)
		}
	} else {
		c.JSON(http.StatusOK, e)
	}
}
