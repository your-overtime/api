package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/your-overtime/api/v2/pkg"
)

// CreateWebhook godoc
// @Tags webhook
// @Summary create a webhook
// @Produce json
// @Consume json
// @Param webhook body pkg.Webhook true "Webhook"
// @Success 201 {object} pkg.Webhook
// @Router /webhook [post]
// @Security BasicAuth
// @Security ApiKeyAuth
func (a *API) CreateWebhook(c *gin.Context) {
	os, err := a.getOvertimeServiceForUserFromRequest(c)
	if err != nil {
		log.Debug(err)
		c.JSON(http.StatusUnauthorized, err)
		return
	}
	var hook pkg.WebhookInput
	if err := c.Bind(&hook); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	saved, err := os.CreateWebhook(hook)
	if err != nil {
		log.Debug(err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusCreated, saved)
}

// GetWebhooks godoc
// @Tags webhook
// @Summary Receives users registered webhooks
// @Produce json
// @Consume json
// @Success 200 {object} []pkg.Webhook
// @Router /webhook [get]
// @Security BasicAuth
// @Security ApiKeyAuth
func (a *API) GetWebhooks(c *gin.Context) {
	os, err := a.getOvertimeServiceForUserFromRequest(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err)
		return
	}
	hooks, err := os.GetWebhooks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, hooks)
}
