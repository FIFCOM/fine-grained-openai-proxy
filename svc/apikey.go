package svc

import (
	"fine-grained-openai-proxy/dao"
)

type ApiKey struct {
	ID  int64  `json:"id"`
	Key string `json:"key"`
}

type ApiKeySvc struct {
}

// All Get all OpenAI API keys
func (s *ApiKeySvc) All() ([]*ApiKey, error) {
	apikeys, err := dao.GetAllApiKeys()
	if err != nil {
		return nil, err
	}
	var ret []*ApiKey
	for _, apikey := range apikeys {
		ret = append(ret, s.conv(&apikey))
	}
	return ret, nil
}

// ByID Get OpenAI API key by ID
func (s *ApiKeySvc) ByID(id int64) (*ApiKey, error) {
	apikey, err := dao.GetApiKeyByID(id)
	if err != nil {
		return nil, err
	}
	return s.conv(&apikey), nil
}

// Insert Insert a new OpenAI API key
func (s *ApiKeySvc) Insert(key string) error {
	apikey := dao.ApiKey{Key: key}
	if err := dao.InsertApiKey(apikey); err != nil {
		return err
	}
	return nil
}

// Delete Delete an OpenAI API key
func (s *ApiKeySvc) Delete(key *ApiKey) error {
	if err := dao.DeleteApiKey(dao.ApiKey(*key)); err != nil {
		return err
	}
	return nil
}

// conv Convert dao.ApiKey to svc.ApiKey
func (s *ApiKeySvc) conv(key *dao.ApiKey) *ApiKey {
	return &ApiKey{
		ID:  key.ID,
		Key: key.Key,
	}
}
