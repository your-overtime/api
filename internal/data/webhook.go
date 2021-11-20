package data

import "github.com/your-overtime/api/pkg"

type Webhook struct {
	pkg.Model
	pkg.Webhook
	UserID uint
}

func (d *Db) SaveWebhook(webhook Webhook) (*pkg.Webhook, error) {
	err := d.Conn.Save(&webhook).Error

	if err != nil {
		return nil, err
	}
	return &webhook.Webhook, nil
}

func (d *Db) GetWebhooksByUserID(userID uint) ([]pkg.Webhook, error) {
	var hooks []pkg.Webhook
	err := d.Conn.Where("user_id = ?", userID).Find(&hooks).Error
	if err != nil {
		return nil, err
	}
	return hooks, nil
}
