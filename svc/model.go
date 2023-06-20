package svc

import (
	"fine-grained-openai-proxy/dao"
)

type Model struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type ModelSvc struct {
}

// All Get all available models
func (s *ModelSvc) All() ([]*Model, error) {
	models, err := dao.GetAllModels()
	if err != nil {
		return nil, err
	}
	var ret []*Model
	for _, model := range models {
		ret = append(ret, s.conv(&model))
	}
	return ret, nil
}

// ByID Get Model by ID
func (s *ModelSvc) ByID(id int64) (*Model, error) {
	model, err := dao.GetModelByID(id)
	if err != nil {
		return nil, err
	}
	return s.conv(&model), nil
}

// ByName Get Model by Name
func (s *ModelSvc) ByName(name string) (*Model, error) {
	model, err := dao.GetModelByName(name)
	if err != nil {
		return nil, err
	}
	return s.conv(&model), nil
}

// Insert Insert new Model
func (s *ModelSvc) Insert(name string) error {
	model := dao.Model{Name: name}
	if err := dao.InsertModel(model); err != nil {
		return err
	}
	return nil
}

// Truncate Truncate Model table
func (s *ModelSvc) Truncate() error {
	if err := dao.TruncateModels(); err != nil {
		return err
	}
	return nil
}

// conv Convert dao.Model to svc.Model
func (s *ModelSvc) conv(model *dao.Model) *Model {
	return &Model{
		ID:   model.ID,
		Name: model.Name,
	}
}
