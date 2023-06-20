package dao

/*
Model YOUR available OpenAI API model list.

List:

	resp = curl https://api.openai.com/v1/models \
	  -H "Authorization: Bearer $OPENAI_API_KEY"
	resp['data'][i]['id']

Sort By:

	resp['data'][i]['created'] ascending.
*/
type Model struct {
	ID   int64  `gorm:"type:integer;primary_key;auto_increment"`
	Name string `gorm:"type:text"`
}

func GetAllModels() ([]Model, error) {
	var models []Model
	err := DB.Find(&models).
		Error
	Handle(err)
	return models, err
}

func GetModelByID(id int64) (Model, error) {
	var model Model
	err := DB.Where("id = ?", id).
		First(&model).
		Error
	Handle(err)
	return model, err
}

func GetModelByName(name string) (Model, error) {
	var model Model
	err := DB.Where("name = ?", name).
		First(&model).
		Error
	Handle(err)
	return model, err
}

func InsertModel(model Model) error {
	err := DB.Create(&model).
		Error
	Handle(err)
	return err
}

func TruncateModels() error {
	err := DB.Exec("DELETE FROM models").
		Error
	Handle(err)
	return err
}
