package data

func (d *Db) SaveWebhook(webhook WebhookDB) (*WebhookDB, error) {
	err := d.Conn.Save(&webhook).Error

	if err != nil {
		return nil, err
	}
	return &webhook, nil
}

func (d *Db) GetWebhooksByUserID(userID uint) ([]WebhookDB, error) {
	var hooks []WebhookDB
	err := d.Conn.Where("user_id = ?", userID).Find(&hooks).Error
	if err != nil {
		return nil, err
	}
	return hooks, nil
}
