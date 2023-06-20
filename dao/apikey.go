package dao

type ApiKey struct {
	ID  int64  `gorm:"type:integer;primary_key;auto_increment"`
	Key string `gorm:"type:text;"`
}

func GetAllApiKeys() ([]ApiKey, error) {
	var keys []ApiKey
	err := DB.Find(&keys).
		Error
	Handle(err)
	return keys, err
}

func GetApiKeyByID(id int64) (ApiKey, error) {
	var key ApiKey
	err := DB.Where("id = ?", id).
		First(&key).
		Error
	Handle(err)
	return key, err
}

func InsertApiKey(key ApiKey) error {
	err := DB.Create(&key).
		Error
	Handle(err)
	return err
}

func DeleteApiKey(key ApiKey) error {
	err := DB.Where("id = ?", key.ID).
		Delete(&key).
		Error
	Handle(err)
	return err
}
