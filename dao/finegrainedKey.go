package dao

// FineGrainedKey Fine-grained permission control key.
type FineGrainedKey struct {
	ID          int64       `gorm:"type:integer;primary_key;auto_increment"`
	Hash        string      `gorm:"type:text"`                 // SHA256 hash value
	ParentID    int64       `gorm:"type:integer"`              // id in api_keys table
	Desc        string      `gorm:"type:text"`                 // Description of the key, e.g. sk-xx...xxxx
	Type        string      `gorm:"type:text"`                 // whitelist or blacklist
	List        *Serialized `gorm:"type:text;serializer:json"` // model name list saved in JSON array form
	Expire      int64       `gorm:"type:integer"`              // Unix timestamp of expiration time
	RemainCalls int64       `gorm:"type:integer"`              // Remaining number of calls
}

func GetAllFineGrainedKeys() ([]FineGrainedKey, error) {
	var keys []FineGrainedKey
	err := DB.Find(&keys).
		Error
	Handle(err)
	return keys, err
}

func GetFineGrainedKeyById(id int64) (FineGrainedKey, error) {
	key := FineGrainedKey{}
	err := DB.Where("id = ?", id).
		First(&key).
		Error
	Handle(err)
	return key, err
}

func GetFineGrainedKeyByHash(hash string) (FineGrainedKey, error) {
	key := FineGrainedKey{}
	err := DB.Where("hash = ?", hash).
		First(&key).
		Error
	Handle(err)
	return key, err
}

func GetFineGrainedKeysByParentID(parentID int64) ([]FineGrainedKey, error) {
	var keys []FineGrainedKey
	err := DB.Where("parent_id = ?", parentID).
		Find(&keys).
		Error
	Handle(err)
	return keys, err
}

func InsertFineGrainedKey(key FineGrainedKey) error {
	err := DB.Create(&key).
		Error
	Handle(err)
	return err
}

func UpdateFineGrainedKey(key FineGrainedKey) error {
	err := DB.Save(&key).
		Error
	Handle(err)
	return err
}

func DeleteFineGrainedKey(key FineGrainedKey) error {
	err := DB.Delete(&key).
		Error
	Handle(err)
	return err
}
