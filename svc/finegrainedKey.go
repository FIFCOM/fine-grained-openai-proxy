package svc

import (
	"crypto/rand"
	"encoding/json"
	"fine-grained-openai-proxy/dao"
)

type FineGrainedKey struct {
	ID          int64  `json:"id"`
	Hash        string `json:"hash"`
	ParentID    int64  `json:"parent_id"`
	Desc        string `json:"desc"`
	Type        string `json:"type"`
	List        string `json:"list"`
	Expire      int64  `json:"expire"`
	RemainCalls int64  `json:"remain_calls"`
}

type FineGrainedKeySvc struct {
}

// All Get all fine-grained keys
func (s *FineGrainedKeySvc) All() ([]*FineGrainedKey, error) {
	fgs, err := dao.GetAllFineGrainedKeys()
	if err != nil {
		return nil, err
	}
	var ret []*FineGrainedKey
	for _, fg := range fgs {
		ret = append(ret, s.conv(&fg))
	}
	return ret, nil
}

// ByID Get fine-grained key by ID
func (s *FineGrainedKeySvc) ByID(id int64) (*FineGrainedKey, error) {
	fg, err := dao.GetFineGrainedKeyById(id)
	if err != nil {
		return nil, err
	}
	return s.conv(&fg), nil
}

// ByHash Get fine-grained key by SHA-256 hash
func (s *FineGrainedKeySvc) ByHash(h string) (*FineGrainedKey, error) {
	fg, err := dao.GetFineGrainedKeyByHash(h)
	if err != nil {
		return nil, err
	}
	return s.conv(&fg), nil
}

// ByParentID Get fine-grained key by parent ID
func (s *FineGrainedKeySvc) ByParentID(pid int64) ([]*FineGrainedKey, error) {
	fgs, err := dao.GetFineGrainedKeysByParentID(pid)
	if err != nil {
		return nil, err
	}
	var ret []*FineGrainedKey
	for _, fg := range fgs {
		ret = append(ret, s.conv(&fg))
	}
	return ret, nil
}

// Insert Insert new fine-grained key
func (s *FineGrainedKeySvc) Insert(fg *FineGrainedKey) (string, error) {
	newKey := "sk-" + s.randomString(48)
	desc := newKey[0:8] + "..." + newKey[len(newKey)-4:]
	// convert json to array
	var list []int64
	_ = json.Unmarshal([]byte(fg.List), &list)

	dfg := dao.FineGrainedKey{
		Hash:        Hash(newKey),
		ParentID:    fg.ParentID,
		Desc:        desc,
		Type:        fg.Type,
		List:        &dao.Serialized{List: list},
		Expire:      fg.Expire,
		RemainCalls: fg.RemainCalls,
	}
	if err := dao.InsertFineGrainedKey(dfg); err != nil {
		return "", err
	}
	return newKey, nil
}

// Update update fine-grained key
func (s *FineGrainedKeySvc) Update(fg *FineGrainedKey) error {
	originFg, err := dao.GetFineGrainedKeyById(fg.ID)
	if err != nil {
		return err
	}
	list := &dao.Serialized{}
	_ = json.Unmarshal([]byte(fg.List), &list.List)
	dfg := dao.FineGrainedKey{
		ID:          fg.ID,
		Hash:        originFg.Hash,
		ParentID:    fg.ParentID,
		Desc:        originFg.Desc,
		Type:        fg.Type,
		List:        list,
		Expire:      fg.Expire,
		RemainCalls: fg.RemainCalls,
	}
	if err := dao.UpdateFineGrainedKey(dfg); err != nil {
		return err
	}
	return nil
}

// Delete delete fine-grained key
func (s *FineGrainedKeySvc) Delete(fg *FineGrainedKey) error {
	list := &dao.Serialized{}
	_ = json.Unmarshal([]byte(fg.List), &list.List)
	dfg := dao.FineGrainedKey{
		ID: fg.ID,
	}
	if err := dao.DeleteFineGrainedKey(dfg); err != nil {
		return err
	}
	return nil
}

// conv Convert dao.FineGrainedKey to svc.FineGrainedKey
func (s *FineGrainedKeySvc) conv(fg *dao.FineGrainedKey) *FineGrainedKey {
	var list []byte
	list, _ = json.Marshal(fg.List.List)
	return &FineGrainedKey{
		ID:          fg.ID,
		Hash:        fg.Hash,
		ParentID:    fg.ParentID,
		Desc:        fg.Desc,
		Type:        fg.Type,
		List:        string(list),
		Expire:      fg.Expire,
		RemainCalls: fg.RemainCalls,
	}
}

// randomString Generate a random string of the specified length
func (s *FineGrainedKeySvc) randomString(length int) string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	bytes := make([]byte, length)
	_, _ = rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = chars[b%byte(len(chars))]
	}
	return string(bytes)
}
